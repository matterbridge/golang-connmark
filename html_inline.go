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

func ruleHTMLInline(s *StateInline, silent bool) (_ bool) {
	if !s.Md.HTML {
		return
	}

	pos := s.Pos
	src := s.Src
	if pos+2 >= s.PosMax || src[pos] != '<' {
		return
	}

	if !htmlSecond[src[pos+1]] {
		return
	}

	match := matchHTML(src[pos:])
	if match == "" {
		return
	}

	if !silent {
		s.PushToken(&HTMLInline{Content: match})
	}

	s.Pos += len(match)

	return true
}
