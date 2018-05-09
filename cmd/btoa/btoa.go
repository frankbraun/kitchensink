// btoa computes the ascii85 encoding of a binary file.
package main

import (
	"encoding/ascii85"
	"fmt"
	"io/ioutil"
	"os"
)

func btoa(filename string) error {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	dst := make([]byte, ascii85.MaxEncodedLen(len(src)))
	n := ascii85.Encode(dst, src)
	fmt.Println(string(dst[:n]))

	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s binary_file\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	if err := btoa(os.Args[1]); err != nil {
		fatal(err)
	}
}
