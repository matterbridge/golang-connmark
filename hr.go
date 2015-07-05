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

var hr [256]bool

func init() {
	hr['*'], hr['-'], hr['_'] = true, true, true
}

func ruleHR(s *StateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.TShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.BMarks[startLine] + shift
	src := s.Src

	marker := src[pos]

	if !hr[marker] {
		return
	}

	pos++
	max := s.EMarks[startLine]

	count := 1
	for pos < max {
		c := src[pos]
		pos++
		if c != marker && c != ' ' {
			return
		}
		if c == marker {
			count++
		}
	}

	if count < 3 {
		return
	}

	if silent {
		return true
	}

	s.Line = startLine + 1
	s.PushToken(&Hr{
		Map: [2]int{startLine, s.Line},
	})

	return true
}
