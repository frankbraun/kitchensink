// liner demoes command line editor with history.
package main

import (
	"fmt"
	"os"

	"github.com/peterh/liner"
)

func readline() error {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	for {
		ln, err := line.Prompt("> ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				fmt.Fprintf(os.Stderr, "aborting...\n")
				return nil
			}
			return err
		}
		fmt.Println(ln)
		line.AppendHistory(ln)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	if err := readline(); err != nil {
		fatal(err)
	}
}
