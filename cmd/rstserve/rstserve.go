// rstserve serves reStructuredText files as HTML on localhost.
package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/frankbraun/kitchensink/markup"
)

var rstRenderer = "rst2html.py"

type reStructuredTextRenderer struct{}

func (r reStructuredTextRenderer) Render(filename string) ([]byte, error) {
	fmt.Println("rendering", filename)
	return exec.Command(rstRenderer, filename).Output()
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s rst_file\n", os.Args[0])
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	err := markup.Serve(new(reStructuredTextRenderer), os.Args[1])
	if err != nil {
		fatal(err)
	}
}
