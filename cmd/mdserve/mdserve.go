// mdserve serves Markdown files as HTML on localhost.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/frankbraun/kitchensink/markup"
	"gopkg.in/russross/blackfriday.v2"
)

type markdownRenderer struct {
	toc bool
}

func (r markdownRenderer) Render(filename string) ([]byte, error) {
	fmt.Println("rendering", filename)
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	flags := blackfriday.CommonHTMLFlags
	if r.toc {
		flags |= blackfriday.TOC
	}
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: flags,
	})
	md := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
	node := md.Parse(input)
	var outbuf bytes.Buffer
	renderer.RenderHeader(&outbuf, node)
	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return renderer.RenderNode(&outbuf, node, entering)
	})
	renderer.RenderFooter(&outbuf, node)
	return outbuf.Bytes(), nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options] markdown_file\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	toc := flag.Bool("toc", false, "generate table of contents (TOC)")
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}
	mdr := markdownRenderer{toc: *toc}
	err := markup.Serve(mdr, flag.Arg(0))
	if err != nil {
		fatal(err)
	}
}
