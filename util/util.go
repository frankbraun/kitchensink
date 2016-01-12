// Copyright (c) 2016 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Package util contains utility functions.
package util

import (
	"flag"
	"fmt"
	"os"
)

// Usage prints the usage of the running command with synopsis and the defined
// options from the flag package to stderr and exits with error code 1.
func Usage(synopsis string) {
	fmt.Fprintf(os.Stderr, "Usage: %s %s\n", os.Args[0], synopsis)
	flag.PrintDefaults()
	os.Exit(1)
}
