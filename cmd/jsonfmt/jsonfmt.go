// jsonfmt formats JSON files in place.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func formatJSON(filename string) error {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var jsn interface{}
	if err := json.Unmarshal(in, &jsn); err != nil {
		return err
	}
	out, err := json.MarshalIndent(jsn, "", "  ")
	if err != nil {
		return err
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
	if err := os.Rename(fp.Name(), filename); err != nil {
		return err
	}
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s JSON_file\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	if err := formatJSON(os.Args[1]); err != nil {
		fatal(err)
	}
}
