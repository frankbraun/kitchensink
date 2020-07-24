// jsonfmt formats JSON files in place.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tidwall/pretty"
)

func formatJSON(filename string, bePretty bool) error {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var out []byte
	if bePretty {
		out = pretty.Pretty(in)
	} else {
		var jsn interface{}
		if err := json.Unmarshal(in, &jsn); err != nil {
			return err
		}
		out, err = json.MarshalIndent(jsn, "", "  ")
		if err != nil {
			return err
		}
	}
	fp, err := ioutil.TempFile("", "jsonfmt")
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(fp, string(out)); err != nil {
		fp.Close()
		return err
	}
	if err := fp.Close(); err != nil {
		return err
	}
	return os.Rename(fp.Name(), filename)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s JSON_file\n", os.Args[0])
	os.Exit(2)
}

func main() {
	bePretty := flag.Bool("p", false, "Use pretty printer")
	flag.Parse()
	if flag.NArg() > 1 {
		usage()
	}
	if err := formatJSON(flag.Arg(0), *bePretty); err != nil {
		fatal(err)
	}
}
