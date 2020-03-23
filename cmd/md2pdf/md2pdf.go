// md2pdf converts Markdown content containing frontmatter into PDFs.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gohugoio/hugo/parser"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
)

func pandoc(content []byte, pdfFile string, toc, xelatex bool) error {
	args := []string{
		"--standalone", "-o", pdfFile,
		//"--variable", "classoption=twocolumn",
		"--variable", "papersize=a4paper",
		"--variable", "links-as-notes",
		"--filter", "pandoc-citeproc",
	}
	if toc {
		args = append(args, "--toc")
	}
	if xelatex {
		args = append(args, "--latex-engine=xelatex")
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
		fmt.Fprintln(os.Stderr, string(out))
		return err
	}
	return nil
}

func md2pdf(mdFile, pdfFile string, toc, hugo, xelatex bool) error {
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
	if hugo {
		base := filepath.Base(mdFile)
		category := filepath.Base(filepath.Dir(mdFile))
		date, ok := yml["date"].(string)
		if ok && date != "" {
			md = append(md, []byte(date+" ")...)
		}
		pdfFile := strings.TrimSuffix(base, ".txt") + ".pdf"
		txtFile := strings.TrimSuffix(base, ".txt") + ".txt"
		download := fmt.Sprintf("[read as [txt](/%s/%s) or [PDF](/%s/%s)]\n\n",
			category, txtFile, category, pdfFile)
		md = append(md, []byte(download)...)
		md = append(md, []byte("<!--more-->\n")...)
	}
	re, err := regexp.Compile("^.+\n=+\n")
	if err != nil {
		return err
	}
	content := re.ReplaceAll(page.Content(), nil)
	md = append(md, content...)

	//fmt.Println(string(md))

	if hugo {
		if err := ioutil.WriteFile(pdfFile, md, 0600); err != nil {
			return err
		}
	} else {
		if err := pandoc(md, pdfFile, toc, xelatex); err != nil {
			return err
		}
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
	os.Exit(2)
}

func main() {
	hugo := flag.Bool("hugo", false, "generate output for Hugo instead of PDF")
	toc := flag.Bool("toc", false, "generate table of contents (TOC)")
	xelatex := flag.Bool("xelatex", false, "use xelatex engine")
	flag.Parse()
	if flag.NArg() != 2 {
		usage()
	}
	err := md2pdf(flag.Arg(0), flag.Arg(1), *toc, *hugo, *xelatex)
	if err != nil {
		fatal(err)
	}
}
