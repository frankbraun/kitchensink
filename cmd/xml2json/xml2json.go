// Copyright (c) 2016 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// xml2json converts XML to JSON.
package main

import (
	"os"

	"github.com/clbanning/mxj"
	"github.com/frankbraun/kitchensink/util"
)

func main() {
	m, err := mxj.NewMapXmlReader(os.Stdin)
	if err != nil {
		util.Fatal(err)
	}
	if err := m.JsonIndentWriter(os.Stdout, "", "  "); err != nil {
		util.Fatal(err)
	}
	if _, err := os.Stdout.WriteString("\n"); err != nil {
		util.Fatal(err)
	}
}
