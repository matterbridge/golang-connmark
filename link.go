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

func ruleLink(s *StateInline, silent bool) (_ bool) {
	pos := s.Pos
	src := s.Src

	if src[pos] != '[' {
		return
	}

	labelStart := pos + 1
	labelEnd := parseLinkLabel(s, pos, true)
	if labelEnd < 0 {
		return
	}

	var href, title, label string
	oldPos := pos
	pos = labelEnd + 1
	max := s.PosMax
	if pos < max && src[pos] == '(' {
		pos = skipws(src, pos+1, max)
		if pos >= max {
			return
		}

		url, endpos, ok := parseLinkDestination(src, pos, s.PosMax)
		if ok {
			url = normalizeLink(url)
			if validateLink(url) {
				href = url
				pos = endpos
			}
		}

		start := pos
		pos = skipws(src, pos, max)
		if pos >= max {
			return
		}

		title, _, endpos, ok = parseLinkTitle(src, pos, s.PosMax)
		if pos < max && start != pos && ok {
			pos = skipws(src, endpos, max)
		}

		if pos >= max || src[pos] != ')' {
			s.Pos = oldPos
			return
		}

		pos++

	} else {
		if s.Env.References == nil {
			return
		}

		pos = skipws(src, pos, max)

		if pos < max && src[pos] == '[' {
			start := pos + 1
			pos = parseLinkLabel(s, pos, false)
			if pos >= 0 {
				label = src[start:pos]
				pos++
			} else {
				pos = labelEnd + 1
			}
		} else {
			pos = labelEnd + 1
		}

		if label == "" {
			label = src[labelStart:labelEnd]
		}

		ref, ok := s.Env.References[normalizeReference(label)]
		if !ok {
			s.Pos = oldPos
			return
		}

		href = ref["href"]
		title = ref["title"]
	}

	if !silent {
		s.Pos = labelStart
		s.PosMax = labelEnd

		s.PushOpeningToken(&LinkOpen{
			Href:  href,
			Title: title,
		})

		s.Md.Inline.Tokenize(s)

		s.PushClosingToken(&LinkClose{})
	}

	s.Pos = pos
	s.PosMax = max

	return true
}
