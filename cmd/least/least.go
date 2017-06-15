package main

import (
	"fmt"
	"os"

	"github.com/frankbraun/kitchensink/textbuffer"
	"github.com/gdamore/tcell"
)

func draw(s tcell.Screen, tb *textbuffer.TextBuffer, w, h int) {
	for y := 0; y < h && y < tb.Lines(); y++ {
		x := 0
		lineLen := tb.LineLenCell(y)
		for x < w && x < lineLen {
			c, cw := tb.GetCell(x, y)
			if x+cw <= lineLen {
				s.SetContent(x, y, c[0], c[1:], tcell.StyleDefault)
			}
			x += cw
		}
	}
}

func least(filename string) error {
	tb, err := textbuffer.NewFile(filename)
	if err != nil {
		return err
	}
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
	w, h := s.Size()
	s.Clear()
	draw(s, tb, w, h)
	s.Show()
	for {
		ev := s.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return nil
		case *tcell.EventResize:
			w, h = s.Size()
			s.Clear()
			draw(s, tb, w, h)
			s.Sync()
		}
	}
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
	if err := least(os.Args[1]); err != nil {
		fatal(err)
	}
}
