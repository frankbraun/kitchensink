// Copyright (c) 2017 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// mdserve serves Markdown files as HTML on localhost.
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/mutecomm/mute/util/browser"
	"github.com/russross/blackfriday"
)

var html []byte

func serveMarkdown(filename string) error {
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	html = blackfriday.MarkdownCommon(md)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(html)
	})
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf("http://localhost:%d/", port)
	fmt.Println("serving at", addr)
	go browser.Open(addr)
	return http.Serve(listener, nil)
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
	if err := serveMarkdown(os.Args[1]); err != nil {
		fatal(err)
	}
}
