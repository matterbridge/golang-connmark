// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package markdown

func ruleBackticks(s *StateInline, silent bool) (_ bool) {
	pos := s.Pos
	src := s.Src

	if src[pos] != '`' {
		return
	}

	start := pos
	pos++
	max := s.PosMax

	for pos < max && src[pos] == '`' {
		pos++
	}

	marker := src[start:pos]

	end := pos

	for {
		for start = end; start < max && src[start] != '`'; start++ {
			// do nothing
		}
		if start >= max {
			break
		}
		end = start + 1

		for end < max && src[end] == '`' {
			end++
		}

		if end-start == len(marker) {
			if !silent {
				s.PushToken(&CodeInline{
					Content: normalizeInlineCode(src[pos:start]),
				})
			}
			s.Pos = end
			return true
		}
	}

	if !silent {
		s.Pending.WriteString(marker)
	}

	s.Pos += len(marker)

	return true
}
