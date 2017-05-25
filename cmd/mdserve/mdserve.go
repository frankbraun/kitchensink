// Copyright (c) 2017 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// mdserve serves Markdown files as HTML on localhost.
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/frankbraun/kitchensink/markup"
	"github.com/russross/blackfriday"
)

const (
	commonHtmlFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	commonExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS
)

type markdownRenderer struct{}

func (r markdownRenderer) Render(filename string) ([]byte, error) {
	fmt.Println("rendering", filename)
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	renderer := blackfriday.HtmlRenderer(commonHtmlFlags, "", "")
	options := blackfriday.Options{
		Extensions: commonExtensions,
	}
	return blackfriday.MarkdownOptions(md, renderer, options), nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s markdown_file\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	err := markup.Serve(new(markdownRenderer), os.Args[1])
	if err != nil {
		fatal(err)
	}
}
