// Copyright (c) 2017 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// lesser — less than less. Work in progress.
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type root struct {
	app *views.Application

	views.BoxLayout
}

func (r *root) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlL:
			r.app.Refresh()
			return true
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'Q', 'q':
				r.app.Quit()
				return true
			}
		}
	}
	return r.BoxLayout.HandleEvent(ev)
}

func (r *root) Draw() {
	r.BoxLayout.Draw()
}

func lesser(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	app := &views.Application{}
	app.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))

	root := &root{app: app}
	root.SetOrientation(views.Vertical)

	ta := views.NewTextArea()
	ta.SetContent(string(buf))

	root.InsertWidget(0, ta, 1.0)

	app.SetRootWidget(root)
	if err := app.Run(); err != nil {
		return err
	}
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s filename\n", os.Args[0])
	os.Exit(1)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	if err := lesser(os.Args[1]); err != nil {
		fatal(err)
	}
}
