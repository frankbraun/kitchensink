package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gdamore/tcell"
	"github.com/mutecomm/mute/tui/textbuffer"
)

var (
	curX int
	curY int
	maxX int
	maxY int
)

func draw(s tcell.Screen, tb *textbuffer.TextBuffer, w, h int) {
	for y := 0; y+1 < h && y+curY < tb.Lines(); y++ {
		x := 0
		lineLen := tb.LineLenCell(y + curY)
		for x < w && x+curX < lineLen {
			c, cw := tb.GetCell(x+curX, y+curY)
			if x+cw <= lineLen {
				s.SetContent(x, y, c[0], c[1:], tcell.StyleDefault)
			}
			x += cw
		}
	}
	// draw status bar
	bar := fmt.Sprintf("w=%d, h=%d, curX=%d, curY=%d, maxX=%d, maxY=%d",
		w, h, curX, curY, maxX, maxY)
	style := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)
	for x := 0; x < w; x++ {
		var r rune
		if x < len(bar) {
			r = rune(bar[x])
		}
		s.SetContent(x, h-1, r, nil, style)
	}
}

func redraw(s tcell.Screen, tb *textbuffer.TextBuffer, w, h int) {
	s.Clear()
	draw(s, tb, w, h)
	s.Show()
}

func least(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	tb := textbuffer.New(buf)
	maxY = tb.Lines()
	for y := 0; y < maxY; y++ {
		if tb.LineLenCell(y) > maxX {
			maxX = tb.LineLenCell(y)
		}
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
	redraw(s, tb, w, h)
	for {
		ev := s.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			switch ev.(*tcell.EventKey).Rune() {
			case ' ': // page down
				curY += h
				if curY >= maxY {
					curY = maxY - h
				}
				redraw(s, tb, w, h)
			case 'b': // page down
				curY -= h
				if curY < 0 {
					curY = 0
				}
				redraw(s, tb, w, h)
			case 'g': // top
				curY = 0
				redraw(s, tb, w, h)
			case 'G': // bottom
				curY = maxY - h - 1
				if curY < 0 {
					curY = 0
				}
				redraw(s, tb, w, h)
			case 'j': // down
				if curY < maxY {
					curY++
				}
				redraw(s, tb, w, h)
			case 'k': // up
				if curY > 0 {
					curY--
				}
				redraw(s, tb, w, h)
			case 'h': // left
				if curX > 0 {
					curX--
				}
				redraw(s, tb, w, h)
			case 'l': // right
				if curX < maxX {
					curX++
				}
				redraw(s, tb, w, h)
			case 'q': // quit
				return nil
			}
		case *tcell.EventResize:
			w, h = s.Size()
			s.Clear()
			draw(s, tb, w, h)
			s.Sync()
		}
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s filename\n", os.Args[0])
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	if err := least(os.Args[1]); err != nil {
		fatal(err)
	}
}
