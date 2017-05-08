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

	"github.com/fsnotify/fsnotify"
	"github.com/mutecomm/mute/util/browser"
	"github.com/russross/blackfriday"
)

var html []byte

func render(filename string) ([]byte, error) {
	md, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	html := blackfriday.MarkdownCommon(md)
	fmt.Println("rendering", filename)
	return html, nil
}

func serveMarkdown(filename string) error {
	// render markdown
	var err error
	html, err = render(filename)
	if err != nil {
		return err
	}
	// watch markdown file for changes and rerender if necessary
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	if err := watcher.Add(filename); err != nil {
		return err
	}
	go func() {
		for {
			var err error
			select {
			case event := <-watcher.Events:
				fmt.Println("event:", event)
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					if err := watcher.Add(filename); err != nil {
						fmt.Println("error:", err)
					}
					html, err = render(filename)
					if err != nil {
						fmt.Println("error:", err)
					}
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					html, err = render(filename)
					if err != nil {
						fmt.Println("error:", err)
					}
				}
			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()
	// serve markdown
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
