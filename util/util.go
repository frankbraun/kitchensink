// Package util contains utility functions.
package util

import (
	crand "crypto/rand"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
)

// Fatal prints err to stderr and exits the process with exit code 1.
func Fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

// Usage prints the usage of the running command with synopsis and the defined
// options from the flag package to stderr and exits with error code 1.
func Usage(synopsis string) {
	fmt.Fprintf(os.Stderr, "usage: %s %s\n", os.Args[0], synopsis)
	flag.PrintDefaults()
	os.Exit(1)
}

// Rand returns a uniform random value in [0, max). It panics if max <= 0.
func Rand(rand io.Reader, max int64) (int64, error) {
	n, err := crand.Int(rand, big.NewInt(max))
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}
