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

var paragraphTerminatedBy []BlockRule

func ruleParagraph(s *StateBlock, startLine, _ int, _ bool) bool {
	nextLine := startLine + 1
	endLine := s.LineMax

outer:
	for ; nextLine < endLine && !s.IsLineEmpty(nextLine); nextLine++ {
		shift := s.TShift[nextLine]
		if shift < 0 || shift-s.BlkIndent > 3 {
			continue
		}

		for _, r := range paragraphTerminatedBy {
			if r(s, nextLine, endLine, true) {
				break outer
			}
		}
	}

	content := strings.TrimSpace(s.Lines(startLine, nextLine, s.BlkIndent, false))

	s.Line = nextLine

	s.PushOpeningToken(&ParagraphOpen{
		Map: [2]int{startLine, s.Line},
	})
	s.PushToken(&Inline{
		Content: content,
		Map:     [2]int{startLine, s.Line},
	})
	s.PushClosingToken(&ParagraphClose{})

	return true
}
