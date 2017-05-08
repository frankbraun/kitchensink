// Copyright (c) 2017 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// markup serves markup files as HTML on localhost.
package markup

import (
	"fmt"
	"net"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/mutecomm/mute/util/browser"
)

var html []byte

type Renderer interface {
	Render(filename string) ([]byte, error)
}

func Serve(r Renderer, filename string) error {
	// render markdown
	var err error
	html, err = r.Render(filename)
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
					html, err = r.Render(filename)
					if err != nil {
						fmt.Println("error:", err)
					}
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					html, err = r.Render(filename)
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
