// Copyright (c) 2017 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// goast parses Go into an AST.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
)

type visitor struct {
	sAST  []interface{}
	depth int
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		indent := strings.Repeat("  ", v.depth)
		fmt.Printf("%s%s\n", indent, reflect.TypeOf(node).String())

		/*
			    TODO: store nodes
				switch n := node.(type) {
				default:
					fmt.Printf("%s%v\n", indent, n)
				}
		*/

		return &visitor{sAST: v.sAST, depth: v.depth + 1}
	}
	return nil
}

func parseGo(filename string, trace bool) error {
	fset := token.NewFileSet()
	mode := parser.ParseComments
	if trace {
		mode |= parser.Trace
	}
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		// pkgs
		_, err := parser.ParseDir(fset, filename, nil, mode)
		if err != nil {
			return err
		}
	} else {
		astf, err := parser.ParseFile(fset, filename, nil, mode)
		if err != nil {
			return err
		}
		sAST := make([]interface{}, 0)
		ast.Walk(&visitor{sAST: sAST, depth: 0}, astf)
		jsn, err := json.MarshalIndent(sAST, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsn))
	}
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s Go_file\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	trace := flag.Bool("t", false, "print a trace of parsed productions")
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}
	if err := parseGo(flag.Arg(0), *trace); err != nil {
		fatal(err)
	}
}
