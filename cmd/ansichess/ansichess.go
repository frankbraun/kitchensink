// Copyright (c) 2017 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// ansichess shows a chess board in an ANSI terminal.
package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

var (
	chessPawn   = '♟'
	chessPieces = []rune{'♜', '♞', '♝', '♛', '♚', '♝', '♞', '♜'}
)

// markColumns marks the columns (called files) for the chess board
// with letters from 'a' to 'h'.
func markColumns(s tcell.Screen, y int) {
	s.SetContent(1, y, 'a', nil, tcell.StyleDefault)
	s.SetContent(2, y, 'b', nil, tcell.StyleDefault)
	s.SetContent(3, y, 'c', nil, tcell.StyleDefault)
	s.SetContent(4, y, 'd', nil, tcell.StyleDefault)
	s.SetContent(5, y, 'e', nil, tcell.StyleDefault)
	s.SetContent(6, y, 'f', nil, tcell.StyleDefault)
	s.SetContent(7, y, 'g', nil, tcell.StyleDefault)
	s.SetContent(8, y, 'h', nil, tcell.StyleDefault)
}

// markRows marks the rows (called ranks) for the chess board
// with numbers from '1' to '8'.
func markRows(s tcell.Screen, x int) {
	s.SetContent(x, 1, '8', nil, tcell.StyleDefault)
	s.SetContent(x, 2, '7', nil, tcell.StyleDefault)
	s.SetContent(x, 3, '6', nil, tcell.StyleDefault)
	s.SetContent(x, 4, '5', nil, tcell.StyleDefault)
	s.SetContent(x, 5, '4', nil, tcell.StyleDefault)
	s.SetContent(x, 6, '3', nil, tcell.StyleDefault)
	s.SetContent(x, 7, '2', nil, tcell.StyleDefault)
	s.SetContent(x, 8, '1', nil, tcell.StyleDefault)
}

func drawBoard(s tcell.Screen) {
	for x := 1; x <= 8; x++ {
		for y := 1; y <= 8; y++ {
			style := tcell.StyleDefault
			if (x+y)%2 == 0 {
				style = style.Background(tcell.ColorYellow)
			} else {
				style = style.Background(tcell.ColorBlue)
			}
			var c rune
			switch y {
			case 1:
				style = style.Foreground(tcell.ColorBlack)
				c = chessPieces[x-1]
			case 2:
				style = style.Foreground(tcell.ColorBlack)
				c = chessPawn
			case 7:
				style = style.Foreground(tcell.ColorWhite)
				c = chessPawn
			case 8:
				style = style.Foreground(tcell.ColorWhite)
				c = chessPieces[x-1]
			}
			s.SetContent(x, y, c, nil, style)
		}
	}
}

func ansiChess() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	defer s.Fini()
	if err = s.Init(); err != nil {
		return err
	}
	w, h := s.Size()
	if w < 10 || h < 10 {
		return fmt.Errorf("screen to small (%dx%d), must be at least 10x10",
			w, h)
	}
	s.Clear()
	markColumns(s, 0)
	markColumns(s, 9)
	markRows(s, 0)
	markRows(s, 9)
	drawBoard(s)
	s.Show()
	for {
		ev := s.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return nil
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	if err := ansiChess(); err != nil {
		fatal(err)
	}
}
