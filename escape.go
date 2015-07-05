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

var escaped = make([]bool, 256)

func init() {
	for _, b := range "\\!\"#$%&'()*+,./:;<=>?@[]^_`{|}~-" {
		escaped[b] = true
	}
}

func ruleEscape(s *StateInline, silent bool) (_ bool) {
	pos := s.Pos
	src := s.Src

	if src[pos] != '\\' {
		return
	}

	pos++
	max := s.PosMax

	if pos < max {
		b := src[pos]

		if b < 0x7f && escaped[b] {
			if !silent {
				s.Pending.WriteByte(b)
			}
			s.Pos += 2
			return true
		}

		if b == '\n' {
			if !silent {
				s.PushToken(&Hardbreak{})
			}

			pos++

			for pos < max && src[pos] == ' ' {
				pos++
			}

			s.Pos = pos
			return true
		}
	}

	if !silent {
		s.Pending.WriteByte('\\')
	}

	s.Pos++

	return true
}
