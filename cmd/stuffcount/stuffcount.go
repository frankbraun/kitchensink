// stuffcoun counts stuff in a stuff.md file.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func printCategory(category string, counter, toPurge int) {
	if category != "" {
		fmt.Printf("%s %d", category, counter)
		if toPurge > 0 {
			fmt.Printf(" (%d to purge)", toPurge)
		}
		fmt.Printf("\n")
	}
}

func countStuff(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(f)
	var (
		category string
		counter  int
		toPurge  int
		total    int
	)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "-   ") && strings.HasSuffix(line, ":") {
			printCategory(category, counter, toPurge)
			category = strings.TrimPrefix(line, "-   ")
			counter = 0
			toPurge = 0
		} else {
			counter++
			if strings.Contains(line, "-   ~~") {
				toPurge++
			}
			total++
		}
	}
	if err := s.Err(); err != nil {
		return err
	}
	printCategory(category, counter, toPurge)
	fmt.Println("----------------------------------------")
	fmt.Printf("total: %d\n", total)
	return err
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s stuff.md\n", os.Args[0])
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	if err := countStuff(os.Args[1]); err != nil {
		fatal(err)
	}
}
