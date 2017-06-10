// Copyright (c) 2017 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// tcelltest takes the tcell package for a test drive.
package main

import (
	"fmt"
	"os"

	"github.com/frankbraun/tcell"
)

func draw(s tcell.Screen) {
	w, h := s.Size()
	for x := 0; x < w; x++ {
		s.SetContent(x, 0, tcell.RuneBlock, nil, tcell.StyleDefault)
		s.SetContent(x, h-1, tcell.RuneBlock, nil, tcell.StyleDefault)
	}
	for y := 1; y+1 < h; y++ {
		s.SetContent(0, y, tcell.RuneBlock, nil, tcell.StyleDefault)
		s.SetContent(w-1, y, tcell.RuneBlock, nil, tcell.StyleDefault)
	}
}

func tcelltest() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	defer s.Fini()
	if err = s.Init(); err != nil {
		return err
	}
	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()
	draw(s)
	s.Show()
	for {
		ev := s.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return nil
		case *tcell.EventResize:
			s.Clear()
			draw(s)
			s.Sync()
		}
	}
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	if err := tcelltest(); err != nil {
		fatal(err)
	}
}
