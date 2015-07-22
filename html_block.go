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

import (
	"regexp"
	"strings"

	"github.com/opennota/byteutil"
)

var (
	htmlBlocks = []string{
		"address",
		"article",
		"aside",
		"base",
		"basefont",
		"blockquote",
		"body",
		"caption",
		"center",
		"col",
		"colgroup",
		"dd",
		"details",
		"dialog",
		"dir",
		"div",
		"dl",
		"dt",
		"fieldset",
		"figcaption",
		"figure",
		"footer",
		"form",
		"frame",
		"frameset",
		"h1",
		"head",
		"header",
		"hr",
		"html",
		"legend",
		"li",
		"link",
		"main",
		"menu",
		"menuitem",
		"meta",
		"nav",
		"noframes",
		"ol",
		"optgroup",
		"option",
		"p",
		"param",
		"pre",
		"section",
		"source",
		"summary",
		"table",
		"tbody",
		"td",
		"tfoot",
		"th",
		"thead",
		"title",
		"tr",
		"track",
		"ul",
	}

	htmlBlocksSet = make(map[string]bool)

	htmlSecond [256]bool

	rStartCond1 = regexp.MustCompile(`(?i)^(pre|script|style)([\n\t >]|$)`)
	rEndCond1   = regexp.MustCompile(`(?i)</(pre|script|style)>`)
	rStartCond6 *regexp.Regexp
	rStartCond7 = regexp.MustCompile(`(?i)^(/[a-z][a-z0-9-]*|[a-z][a-z0-9-]*(\s+[a-z_:][a-z0-9_.:-]*\s*=\s*("[^"]*"|'[^']*'|[ "'=<>\x60]))*\s*/?)>\s*$`)
)

func init() {
	for _, tag := range htmlBlocks {
		htmlBlocksSet[tag] = true
	}
	for _, b := range "!/?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		htmlSecond[b] = true
	}
	rStartCond6 = regexp.MustCompile(`(?i)^/?(` + strings.Join(htmlBlocks, "|") + `)(\s|$|>|/>)`)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func matchTagName(s string) string {
	if len(s) < 2 {
		return ""
	}

	i := 0
	if s[0] == '/' {
		i++
	}
	start := i
	max := min(15+i, len(s))
	for i < max && byteutil.IsLetter(s[i]) {
		i++
	}
	if i >= len(s) {
		return ""
	}

	switch s[i] {
	case ' ', '\n', '/', '>':
		return byteutil.ToLower(s[start:i])
	default:
		return ""
	}
}

func ruleHTMLBlock(s *StateBlock, startLine, endLine int, silent bool) (_ bool) {
	if !s.Md.HTML {
		return
	}

	shift := s.TShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.BMarks[startLine] + shift
	max := s.EMarks[startLine]

	if pos+1 >= max {
		return
	}

	src := s.Src

	if src[pos] != '<' {
		return
	}

	pos++
	b := src[pos]
	if !htmlSecond[b] {
		return
	}

	nextLine := startLine + 1

	var endCond func(string) bool

	if pos+2 < max && byteutil.IsLetter(b) && rStartCond1.MatchString(src[pos:]) {
		endCond = func(s string) bool {
			return rEndCond1.MatchString(s)
		}
	} else if strings.HasPrefix(src[pos:], "!--") {
		endCond = func(s string) bool {
			return strings.Contains(s, "-->")
		}
	} else if b == '?' {
		endCond = func(s string) bool {
			return strings.Contains(s, "?>")
		}
	} else if b == '!' && pos+1 < max && byteutil.IsUppercaseLetter(src[pos+1]) {
		endCond = func(s string) bool {
			return strings.Contains(s, ">")
		}
	} else if strings.HasPrefix(src[pos:], "![CDATA[") {
		endCond = func(s string) bool {
			return strings.Contains(s, "]]>")
		}
	} else if pos+2 < max && (byteutil.IsLetter(b) || b == '/' && byteutil.IsLetter(src[pos+1])) {
		terminator := true
		if rStartCond6.MatchString(src[pos:max]) {
		} else if rStartCond7.MatchString(src[pos:max]) {
			terminator = false
		} else {
			return
		}
		if silent {
			return terminator
		}
		endCond = func(s string) bool {
			return s == ""
		}
	} else {
		return
	}

	if silent {
		return true
	}

	if !endCond(src[pos:max]) {
		for nextLine < endLine {
			shift := s.TShift[nextLine]
			if shift < s.BlkIndent {
				break
			}
			pos := s.BMarks[nextLine] + shift
			max := s.EMarks[nextLine]
			if pos > max {
				pos = max
			}
			if endCond(src[pos:max]) {
				if pos != max {
					nextLine++
				}
				break
			}
			nextLine++
		}
	}

	s.Line = nextLine
	s.PushToken(&HTMLBlock{
		Content: s.Lines(startLine, nextLine, s.BlkIndent, true),
		Map:     [2]int{startLine, nextLine},
	})

	return true
}
