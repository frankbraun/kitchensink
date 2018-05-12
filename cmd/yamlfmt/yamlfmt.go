// yamlfmt formats YAML files in place.
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func formatYAML(filename string) error {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var yml interface{}
	if err := yaml.Unmarshal(in, &yml); err != nil {
		return err
	}
	out, err := yaml.Marshal(yml)
	if err != nil {
		return err
	}
	fp, err := ioutil.TempFile("", "yamlfmt")
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
	fmt.Fprintf(os.Stderr, "usage: %s YAML_file\n", os.Args[0])
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	if err := formatYAML(os.Args[1]); err != nil {
		fatal(err)
	}
}
