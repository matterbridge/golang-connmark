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

import "bytes"

type StateInline struct {
	StateCore

	Pos          int
	PosMax       int
	Level        int
	Pending      bytes.Buffer
	PendingLevel int

	Cache map[int]int
}

func (s *StateInline) PushToken(tok Token) {
	if s.Pending.Len() > 0 {
		s.PushPending()
	}
	tok.SetLevel(s.Level)
	s.PendingLevel = s.Level
	s.Tokens = append(s.Tokens, tok)
}

func (s *StateInline) PushOpeningToken(tok Token) {
	if s.Pending.Len() > 0 {
		s.PushPending()
	}
	tok.SetLevel(s.Level)
	s.Level++
	s.PendingLevel = s.Level
	s.Tokens = append(s.Tokens, tok)
}

func (s *StateInline) PushClosingToken(tok Token) {
	if s.Pending.Len() > 0 {
		s.PushPending()
	}
	s.Level--
	tok.SetLevel(s.Level)
	s.PendingLevel = s.Level
	s.Tokens = append(s.Tokens, tok)
}

func (s *StateInline) PushPending() {
	s.Tokens = append(s.Tokens, &Text{
		Content: s.Pending.String(),
		Lvl:     s.PendingLevel,
	})
	s.Pending.Reset()
}
