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

func ruleNewline(s *StateInline, silent bool) (_ bool) {
	pos := s.Pos
	src := s.Src

	if src[pos] != '\n' {
		return
	}

	pending := s.Pending.Bytes()
	n := len(pending) - 1

	if !silent {
		if n >= 0 && pending[n] == ' ' {
			if n >= 1 && pending[n-1] == ' ' {
				n -= 2
				for n >= 0 && pending[n] == ' ' {
					n--
				}
				s.Pending.Truncate(n + 1)
				s.PushToken(&Hardbreak{})
			} else {
				s.Pending.Truncate(n)
				s.PushToken(&Softbreak{})
			}
		} else {
			s.PushToken(&Softbreak{})
		}
	}

	pos++
	max := s.PosMax

	for pos < max && src[pos] == ' ' {
		pos++
	}

	s.Pos = pos

	return true
}
