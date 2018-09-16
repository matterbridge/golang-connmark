// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "unicode/utf8"

type ParserBlock struct{}

type BlockRule func(*StateBlock, int, int, bool) bool

var blockRules []BlockRule

func (b ParserBlock) Parse(src []byte, md *Markdown, env *Environment) []Token {
	n := 1
	buf := make([]byte, 0, len(src))
	for i := 0; i < len(src); i++ {
		switch ch := src[i]; ch {
		case '\n':
			n++
			buf = append(buf, '\n')
		case '\r':
			buf = append(buf, '\n')
			n++
			if i < len(src) && src[i+1] == '\n' {
				i++
			}
		default:
			buf = append(buf, ch)
		}
	}
	if len(buf) == 0 || buf[len(buf)-1] != '\n' {
		n++
	}
	src = buf

	indentFound := false
	start := 0
	pos := 0
	indent := 0
	offset := 0

	mem := make([]int, 0, n*5)
	bMarks := mem[0:0:n]
	eMarks := mem[n : n : n*2]
	tShift := mem[n*2 : n*2 : n*3]
	sCount := mem[n*3 : n*3 : n*4]
	bsCount := mem[n*4 : n*4 : n*5]

	for pos < len(src) {
		r, size := utf8.DecodeRune(src[pos:])

		if !indentFound {
			if runeIsSpace(r) {
				indent++
				if r == '\t' {
					offset += 4 - offset%4
				} else {
					offset++
				}
				pos += size
				continue
			}
			indentFound = true
		}

		if r == '\n' || pos == len(src)-1 {
			if r != '\n' {
				pos++
			}
			bMarks = append(bMarks, start)
			eMarks = append(eMarks, pos)
			tShift = append(tShift, indent)
			sCount = append(sCount, offset)
			bsCount = append(bsCount, 0)

			indentFound = false
			indent = 0
			offset = 0
			start = pos + 1
		}

		pos += size
	}

	bMarks = append(bMarks, len(src))
	eMarks = append(eMarks, len(src))
	tShift = append(tShift, 0)
	sCount = append(sCount, 0)
	bsCount = append(bsCount, 0)

	var s StateBlock
	s.BMarks = bMarks
	s.EMarks = eMarks
	s.TShift = tShift
	s.SCount = sCount
	s.BSCount = bsCount
	s.LineMax = n - 1
	s.Src = string(src)
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

		if s.SCount[line] < s.BlkIndent {
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
			s.Line = line
		}
	}
}
