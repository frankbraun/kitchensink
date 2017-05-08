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

type markdownRenderer struct{}

func (r markdownRenderer) Render(filename string) ([]byte, error) {
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fmt.Println("rendering", filename)
	return blackfriday.MarkdownCommon(md), nil
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
