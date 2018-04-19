// Copyright (c) 2018 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// randpass prints a random password (128-bit, 192-bit, or 256-bit) in base64
// encoding to stdout.
package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
)

func randPass(size int, url bool) (string, error) {
	pass := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, pass[:]); err != nil {
		return "", err
	}
	if url {
		return base64.RawURLEncoding.EncodeToString(pass[:]), nil
	}
	return base64.RawStdEncoding.EncodeToString(pass[:]), nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	low := flag.Bool("l", false, "use low security (128-bit)")
	high := flag.Bool("h", false, "use high security (256-bit)")
	url := flag.Bool("url", false, "use URL encoding")
	flag.Parse()
	size := 24 // 192-bit
	if *low {
		size = 16 // 128-bit
	}
	if *high {
		size = 32 // 256-bit
	}
	pw, err := randPass(size, *url)
	if err != nil {
		fatal(err)
	}
	fmt.Println(pw)
}
