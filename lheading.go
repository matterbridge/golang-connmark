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

var under [256]bool

func init() {
	under['-'], under['='] = true, true
}

func ruleLHeading(s *StateBlock, startLine, endLine int, silent bool) (_ bool) {
	nextLine := startLine + 1

	if nextLine >= endLine {
		return
	}

	shift := s.TShift[nextLine]
	if shift < s.BlkIndent {
		return
	}

	if shift-s.BlkIndent > 3 {
		return
	}

	pos := s.BMarks[nextLine] + shift
	max := s.EMarks[nextLine]

	if pos >= max {
		return
	}

	src := s.Src
	marker := src[pos]

	if !under[marker] {
		return
	}

	pos = s.SkipBytes(pos, marker)

	pos = s.SkipSpaces(pos)

	if pos < max {
		return
	}

	pos = s.BMarks[startLine] + s.TShift[startLine]

	s.Line = nextLine + 1

	hLevel := 1
	if marker == '-' {
		hLevel++
	}

	s.PushOpeningToken(&HeadingOpen{
		HLevel: hLevel,
		Map:    [2]int{startLine, s.Line},
	})
	s.PushToken(&Inline{
		Content: strings.TrimSpace(src[pos:s.EMarks[startLine]]),
		Map:     [2]int{startLine, s.Line - 1},
	})
	s.PushClosingToken(&HeadingClose{HLevel: hLevel})

	return true
}
