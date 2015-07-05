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

type ParserBlock struct {
}

type BlockRule func(*StateBlock, int, int, bool) bool

var blockRules []BlockRule

func (b ParserBlock) Parse(src []byte, md *Markdown, env *Environment) []Token {
	str, bMarks, eMarks, tShift := normalizeAndIndex(src)
	bMarks = append(bMarks, len(str))
	eMarks = append(eMarks, len(str))
	tShift = append(tShift, 0)
	var s StateBlock
	s.BMarks = bMarks
	s.EMarks = eMarks
	s.TShift = tShift
	s.LineMax = len(bMarks) - 1
	s.Src = str
	s.Md = md
	s.Env = env

	b.Tokenize(&s, s.Line, s.LineMax)

	return s.Tokens
}

func (ParserBlock) Tokenize(s *StateBlock, startLine, endLine int) {
	line := startLine
	hasEmptyLines := false
	maxNesting := s.Md.MaxNesting

	for line < endLine {
		line = s.SkipEmptyLines(line)
		s.Line = line
		if line >= endLine {
			break
		}

		if s.TShift[line] < s.BlkIndent {
			break
		}

		if s.Level >= maxNesting {
			s.Line = endLine
			break
		}

		for _, r := range blockRules {
			if r(s, line, endLine, false) {
				break
			}
		}

		s.Tight = !hasEmptyLines

		if s.IsLineEmpty(s.Line - 1) {
			hasEmptyLines = true
		}

		line = s.Line

		if line < endLine && s.IsLineEmpty(line) {
			hasEmptyLines = true
			line++

			if line < endLine && s.ParentType == ptList && s.IsLineEmpty(line) {
				break
			}
			s.Line = line
		}
	}
}
