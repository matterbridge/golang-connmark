// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import (
	"io/ioutil"
	"testing"

	"github.com/russross/blackfriday"
)

func BenchmarkRenderSpecNoHTML(b *testing.B) {
	b.StopTimer()
	data, err := ioutil.ReadFile("spec/spec-0.28.txt")
	if err != nil {
		b.Fatal(err)
	}

	md := New(HTML(false), XHTMLOutput(true))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		md.RenderToString(data)
	}
}

func BenchmarkRenderSpec(b *testing.B) {
	b.StopTimer()
	data, err := ioutil.ReadFile("spec/spec-0.28.txt")
	if err != nil {
		b.Fatal(err)
	}

	md := New(HTML(true), XHTMLOutput(true))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		md.RenderToString(data)
	}
}

func BenchmarkRenderSpecBlackFriday(b *testing.B) {
	b.StopTimer()
	data, err := ioutil.ReadFile("spec/spec-0.28.txt")
	if err != nil {
		panic(err)
	}

	renderer := blackfriday.HtmlRenderer(blackfriday.HTML_USE_XHTML|blackfriday.HTML_USE_SMARTYPANTS|blackfriday.HTML_SMARTYPANTS_LATEX_DASHES, "", "")
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		blackfriday.Markdown(data, renderer,
			blackfriday.EXTENSION_NO_INTRA_EMPHASIS|
				blackfriday.EXTENSION_TABLES|
				blackfriday.EXTENSION_FENCED_CODE|
				blackfriday.EXTENSION_AUTOLINK|
				blackfriday.EXTENSION_STRIKETHROUGH,
		)
	}
}
