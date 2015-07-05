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

import "strings"

func ruleHeading(s *StateBlock, startLine, _ int, silent bool) (_ bool) {
	shift := s.TShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.BMarks[startLine] + shift
	max := s.EMarks[startLine]
	src := s.Src

	if pos >= max || src[pos] != '#' {
		return
	}

	pos++

	level := 1
	for pos < max && src[pos] == '#' && level <= 6 {
		level++
		pos++
	}

	if level > 6 || (pos < max && src[pos] != ' ') {
		return
	}

	if silent {
		return true
	}

	max = s.SkipBytesBack(max, ' ', pos)
	tmp := s.SkipBytesBack(max, '#', pos)
	if tmp > pos && src[tmp-1] == ' ' {
		max = tmp
	}

	s.Line = startLine + 1

	s.PushOpeningToken(&HeadingOpen{
		HLevel: level,
		Map:    [2]int{startLine, s.Line},
	})

	if pos < max {
		s.PushToken(&Inline{
			Content: strings.TrimSpace(src[pos:max]),
			Map:     [2]int{startLine, s.Line},
		})
	}
	s.PushClosingToken(&HeadingClose{HLevel: level})

	return true
}
