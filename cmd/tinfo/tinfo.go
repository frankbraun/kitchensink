// tinfo prints infos about the terminal it is running in.
package main

import (
	"fmt"
	"os"

	"github.com/frankbraun/tcell"
)

func tinfo() error {
	fmt.Printf("TERM=%s\n", os.Getenv("TERM"))
	fmt.Printf("COLORTERM=%s\n", os.Getenv("COLORTERM"))
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err = s.Init(); err != nil {
		return err
	}
	w, h := s.Size()
	colors := s.Colors()
	s.Fini()
	fmt.Printf("screen.Size(): width=%d, height=%d\n", w, h)
	fmt.Printf("screen.Colors(): %d\n", colors)
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	if err := tinfo(); err != nil {
		fatal(err)
	}
}
