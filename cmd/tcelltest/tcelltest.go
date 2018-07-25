// tcelltest takes the tcell package for a test drive.
package main

import (
	"flag"
	"fmt"

	//"io"
	"os"
	//"os/signal"
	//"syscall"
	//ptmx "github.com/frankbraun/pty"
	//"github.com/frankbraun/tcell"
	//"golang.org/x/crypto/ssh/terminal"
)

/*
func drawBorder(s tcell.Screen) {
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

func tcelltest(pty bool) error {
	var (
		s   tcell.Screen
		err error
	)
	if pty {
		// run tcell through pseudoterminal (pty) master/slave-layer
		stdinfd := int(os.Stdin.Fd())
		state, err := terminal.MakeRaw(stdinfd)
		if err != nil {
			return err
		}
		defer terminal.Restore(stdinfd, state)
		pty, tty, err := ptmx.Open()
		if err != nil {
			return err
		}
		defer pty.Close()
		defer tty.Close()
		sigChan := make(chan os.Signal, 2)
		go func() {
			ptmx.InheritSize(os.Stdin, pty)
			for range sigChan {
				ptmx.InheritSize(os.Stdin, pty)
			}
		}()
		signal.Notify(sigChan, syscall.SIGWINCH)
		go func() {
			io.Copy(os.Stdout, pty)
		}()
		go func() {
			io.Copy(pty, os.Stdin)
		}()
		s, err = tcell.NewTerminfoScreenWithTTY(os.Getenv("TERM"), tty)
		if err != nil {
			return err
		}
	} else {
		// run tcell directly
		s, err = tcell.NewScreen()
		if err != nil {
			return err
		}
	}
	defer s.Fini()
	if err = s.Init(); err != nil {
		return err
	}
	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()
	drawBorder(s)
	s.Show()
	for {
		ev := s.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return nil
		case *tcell.EventResize:
			s.Clear()
			drawBorder(s)
			s.Sync()
		}
	}
}
*/

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s Go_file\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	/*
		pty := flag.Bool("pty", false,
			"run tcell through pseudoterminal (pty) master/slave-layer ")
	*/
	flag.Parse()
	if flag.NArg() != 0 {
		usage()
	}
	/*
		if err := tcelltest(*pty); err != nil {
			fatal(err)
		}
	*/
}
