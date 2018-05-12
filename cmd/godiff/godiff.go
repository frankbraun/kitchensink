// godiff computs a diff between two files
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func diff(fileA, fileB string) error {
	a, err := ioutil.ReadFile(fileA)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(fileB)
	if err != nil {
		return err
	}

	dmp := diffmatchpatch.New()
	textA, textB, lineArray := dmp.DiffLinesToRunes(string(a), string(b))
	diffs := dmp.DiffMainRunes(textA, textB, true)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)
	patches := dmp.PatchMake(string(a), diffs)
	text := dmp.PatchToText(patches)

	fmt.Print(text)
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s file_a file_b\n", os.Args[0])
	os.Exit(2)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}
	if err := diff(os.Args[1], os.Args[2]); err != nil {
		fatal(err)
	}
}
