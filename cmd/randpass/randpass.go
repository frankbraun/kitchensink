// Copyright (c) 2018 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// randpass prints a random 192-bit password in base64 encoding to stdout.
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func randPass() (string, error) {
	var pass [24]byte
	if _, err := io.ReadFull(rand.Reader, pass[:]); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(pass[:]), nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func main() {
	pw, err := randPass()
	if err != nil {
		fatal(err)
	}
	fmt.Println(pw)
}
