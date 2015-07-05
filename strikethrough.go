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

func ruleStrikeThrough(s *StateInline, silent bool) (_ bool) {
	start := s.Pos
	max := s.PosMax
	src := s.Src

	if src[start] != '~' {
		return
	}

	if silent {
		return
	}

	canOpen, canClose, delims := scanDelims(s, start)
	startCount := delims
	if !canOpen {
		s.Pos += startCount
		s.Pending.WriteString(src[start:s.Pos])
		return true
	}

	stack := startCount / 2
	if stack <= 0 {
		return
	}
	s.Pos = start + startCount

	var found bool
	for s.Pos < max {
		if src[s.Pos] == '~' {
			canOpen, canClose, delims = scanDelims(s, s.Pos)
			count := delims
			tagCount := count / 2
			if canClose {
				if tagCount >= stack {
					s.Pos += count - 2
					found = true
					break
				}
				stack -= tagCount
				s.Pos += count
				continue
			}

			if canOpen {
				stack += tagCount
			}
			s.Pos += count
			continue
		}

		s.Md.Inline.SkipToken(s)
	}

	if !found {
		s.Pos = start
		return
	}

	s.PosMax = s.Pos
	s.Pos = start + 2

	s.PushOpeningToken(&StrikethroughOpen{})

	s.Md.Inline.Tokenize(s)

	s.PushClosingToken(&StrikethroughClose{})

	s.Pos = s.PosMax + 2
	s.PosMax = max

	return true
}
