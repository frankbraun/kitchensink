// totaltime shows the total play time of the video content rooted at current
// directory.
package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	ignoreSuffixes = []string{
		".srt",
	}
	videoSuffixes = []string{
		".avi",
		".m4v",
		".mkv",
		".mp4",
		".webm",
		".wmv",
	}
)

func determineDuration(path string) (time.Duration, error) {
	cmd := exec.Command("ffprobe", "-show_format", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "$ ffprobe -show_format %s\n", path)
		fmt.Fprintln(os.Stderr, string(out))
		return 0, err
	}
	start := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Duration") {
			fields := strings.Split(strings.TrimSpace(line), ",")
			durString := strings.TrimPrefix(fields[0], "Duration: ")
			t, err := time.Parse("15:04:05.00", durString)
			if err != nil {
				return 0, err
			}
			return t.Sub(start), nil
		}
	}
	return 0, fmt.Errorf("%s: could not parse duration", path)
}

func totalTime() error {
	var total time.Duration
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !info.Mode().IsRegular() {
			return fmt.Errorf("%s: neither directory nor normal file", path)
		}
		// ignore directories
		if info.IsDir() {
			return nil
		}
		// ignore some suffixes
		for _, suffix := range ignoreSuffixes {
			if strings.HasSuffix(path, suffix) {
				return nil
			}
		}
		// process video suffixes
		for _, suffix := range videoSuffixes {
			if strings.HasSuffix(path, suffix) {
				dur, err := determineDuration(path)
				if err != nil {
					return err
				}
				if dur > math.MaxInt64-total {
					return errors.New("more than 290 years of material")
				}
				total += dur
				return nil
			}
		}
		return fmt.Errorf("%s: unknown suffix", path)
	})
	if err != nil {
		return err
	}
	fmt.Println(total)
	return nil
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Show the total play time of the video content rooted at current directory.\n")
	os.Exit(2)
}

func main() {
	if len(os.Args) != 1 {
		usage()
	}
	if err := totalTime(); err != nil {
		fatal(err)
	}
}
