package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func pandoc(content []byte, pdfFile string, toc bool) error {
	args := []string{"--standalone", "-o", pdfFile}
	if toc {
		args = append(args, "--toc")
	}
	cmd := exec.Command("pandoc", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		stdin.Write(content)
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, out)
		return err
	}
	return nil
}

func md2pdf(mdFile, pdfFile string, toc bool) error {
	src, err := ioutil.ReadFile(mdFile)
	if err != nil {
		return err
	}
	if err := pandoc(src, pdfFile, toc); err != nil {
		return err
	}
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options] markdown_file pdf_file\n",
		os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	toc := flag.Bool("toc", false, "generate table of contents (TOC)")
	flag.Parse()
	if flag.NArg() != 2 {
		usage()
	}
	if err := md2pdf(flag.Arg(0), flag.Arg(1), *toc); err != nil {
		fatal(err)
	}
}
