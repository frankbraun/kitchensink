// jsonfmt formats JSON files in place.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/tidwall/pretty"
)

func renameAccrossFilesystem(src, dst string) error {
	// open source file
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	// get mode of source file
	fi, err := s.Stat()
	if err != nil {
		return err
	}
	// make sure source file is a regular file
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("source file '%s' is not a regular file", src)
	}
	mode := fi.Mode() & os.ModePerm // only keep standard UNIX permission bits
	// create destination file
	d, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer d.Close()
	// copy content
	if _, err := io.Copy(d, s); err != nil {
		return err
	}
	// remove src file
	defer os.Remove(src)
	return nil
}

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
	err = os.Rename(fp.Name(), filename)
	if err != nil {
		err = renameAccrossFilesystem(fp.Name(), filename)
	}
	return err
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
