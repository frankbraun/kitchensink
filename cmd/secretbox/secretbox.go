// secretbox is a simple tool to encrypt and decrypt files with NaCL's secretbox.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/frankbraun/kitchensink/secretbox"
)

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s seal plain_file crypt_file\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "       %s open crypt_file plain_file\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 3 {
		usage()
	}
	switch flag.Arg(0) {
	case "seal":
		if err := secretbox.Seal(flag.Arg(1), flag.Arg(2)); err != nil {
			fatal(err)
		}
	case "open":
		if err := secretbox.Open(flag.Arg(1), flag.Arg(2)); err != nil {
			fatal(err)
		}
	default:
		usage()
	}
}
