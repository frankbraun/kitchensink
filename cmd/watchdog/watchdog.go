// package watches a set of files and executes a command if one of them changes.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type stringArray []string

func (a *stringArray) String() string {
	return strings.Join(*a, ",")
}

func (a *stringArray) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func execute(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	} else {
		fmt.Fprintln(os.Stderr, "success")
	}
}

func watch(files, args []string) error {
	// execute once in any case
	execute(args)
	if len(files) == 0 {
		return errors.New("no watch files specified with -w")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Fprintln(os.Stderr, "event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Fprintln(os.Stderr, "modified file:", event.Name)
					// file has been written, execute
					execute(args)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Fprintln(os.Stderr, "error:", err)
			}
		}
	}()

	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			return err
		}
	}
	<-done
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-w] command [arguments]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Watch files specified with -w and execute command if one of them changes\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var watchFiles stringArray
	flag.Var(&watchFiles, "w", "Watch file")
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}
	if err := watch(watchFiles, flag.Args()); err != nil {
		fatal(err)
	}
}
