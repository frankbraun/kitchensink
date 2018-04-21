// termecho is a test program that echos lines typed in a terminal.
package main

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func termechoMain() error {
	oldState, err := terminal.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return err
	}
	defer terminal.Restore(int(syscall.Stdin), oldState)
	fmt.Println("enter 'q' to quit")
	t := terminal.NewTerminal(os.Stdin, "> ")
	for {
		line, err := t.ReadLine()
		if err != nil {
			return err
		}
		if line == "q" {
			break
		}
		fmt.Println(line)
	}
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	// work around defer not working after os.Exit()
	if err := termechoMain(); err != nil {
		fatal(err)
	}
}
