package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/spf13/hugo/parser"
	"gopkg.in/russross/blackfriday.v2"
	"gopkg.in/yaml.v2"
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
	// parse the file with hugo/parser to extract front matter
	fp, err := os.Open(mdFile)
	if err != nil {
		return err
	}
	defer fp.Close()
	page, err := parser.ReadFrom(fp)
	if err != nil {
		return err
	}

	// parse YAML frontmatter
	yml, err := page.Metadata()
	if err != nil {
		return err
	}

	// parse title (h1) from markdown
	opt := blackfriday.WithExtensions(blackfriday.CommonExtensions)
	mdParser := blackfriday.New(opt)
	ast := mdParser.Parse(page.Content())
	var title string
	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.Heading {
			if node.HeadingData.Level == 1 {
				title = string(node.FirstChild.Literal)
				return blackfriday.Terminate
			}
		}
		return blackfriday.GoToNext
	})

	// add title to frontmatter and remove it from markdown
	md, err := yaml.Marshal(yml)
	if err != nil {
		return err
	}
	title = fmt.Sprintf("---\ntitle: %s\n", title)
	md = append([]byte(title), md...)
	md = append(md, []byte("---\n\n")...)
	re, err := regexp.Compile("^.+\n=+\n")
	if err != nil {
		return err
	}
	content := re.ReplaceAll(page.Content(), nil)
	md = append(md, content...)

	//fmt.Println(string(md))

	if err := pandoc(md, pdfFile, toc); err != nil {
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
