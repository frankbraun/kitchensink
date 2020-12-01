// urldecode decodes URL parameters
package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

func decodeURL(s string) error {
	u, err := url.QueryUnescape(s)
	if err != nil {
		return err
	}
	fmt.Println(u)
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s URL\n", os.Args[0])
	os.Exit(2)
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}
	if err := decodeURL(flag.Arg(0)); err != nil {
		fatal(err)
	}
}
