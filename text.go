// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "strings"

func isTerminatorChar(ch byte) bool {
	return strings.IndexByte("\n!#$%&*+-:<=>@[\\]^_`{}~", ch) != -1
}

func ruleText(s *StateInline, silent bool) bool {
	pos := s.Pos
	max := s.PosMax
	src := s.Src

	for pos < max && !isTerminatorChar(src[pos]) {
		pos++
	}
	if pos == s.Pos {
		return false
	}

	if !silent {
		s.Pending.WriteString(src[s.Pos:pos])
	}

	s.Pos = pos

	return true
}
