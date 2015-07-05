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

var fence [256]bool

func init() {
	fence['~'], fence['`'] = true, true
}

func ruleFence(s *StateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.TShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.BMarks[startLine] + shift
	max := s.EMarks[startLine]
	src := s.Src

	if pos+3 > max {
		return
	}

	marker := src[pos]

	if !fence[marker] {
		return
	}

	mem := pos
	pos = s.SkipBytes(pos, marker)
	len := pos - mem
	if len < 3 {
		return
	}

	params := strings.TrimSpace(src[pos:max])

	if strings.IndexByte(params, '`') >= 0 {
		return
	}

	if silent {
		return true
	}

	nextLine := startLine
	haveEndMarker := false

	for {
		nextLine++
		if nextLine >= endLine {
			break
		}

		mem = s.BMarks[nextLine] + s.TShift[nextLine]
		pos = mem
		max = s.EMarks[nextLine]

		if pos >= max {
			continue
		}

		if s.TShift[nextLine] < s.BlkIndent {
			break
		}

		if src[pos] != marker {
			continue
		}

		if s.TShift[nextLine]-s.BlkIndent > 3 {
			continue
		}

		pos = s.SkipBytes(pos, marker)

		if pos-mem < len {
			continue
		}

		pos = s.SkipSpaces(pos)
		if pos < max {
			continue
		}

		haveEndMarker = true

		break
	}

	s.Line = nextLine
	if haveEndMarker {
		s.Line++
	}

	s.PushToken(&Fence{
		Params:  params,
		Content: s.Lines(startLine+1, nextLine, s.TShift[startLine], true),
		Map:     [2]int{startLine, nextLine},
	})

	return true
}
