// enctest allows to dynamically test different text encodings.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/frankbraun/kitchensink/ascii/fivebit"
	"github.com/frankbraun/kitchensink/ascii/huffman"
	"github.com/frankbraun/kitchensink/ascii/sixbit"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func testEncodings() error {
	textView := tview.NewTextView()
	inputField := tview.NewInputField()
	inputField.SetChangedFunc(func(text string) {
		textView.Clear()
		t := strings.ToUpper(text)
		fmt.Fprintln(textView)
		fmt.Fprintf(textView, "%s\n", t)
		fmt.Fprintln(textView)
		fmt.Fprintf(textView, "8-bit ASCII length in bytes: %d\n", len(text))
		bytes, err := sixbit.EncodedLen(t)
		if err != nil {
			fmt.Fprintf(textView, err.Error()+"\n")
		} else {
			fmt.Fprintf(textView, "6-bit ASCII length in bytes: %d\n", bytes)
		}
		bytes, err = fivebit.EncodedLen(t)
		if err != nil {
			fmt.Fprintf(textView, err.Error()+"\n")
		} else {
			fmt.Fprintf(textView, "5-bit ASCII length in bytes: %d\n", bytes)
		}
		fmt.Fprintln(textView)
		bytes, err = huffman.EncodedLen8Bit(t)
		if err != nil {
			fmt.Fprintf(textView, err.Error()+"\n")
		} else {
			fmt.Fprintf(textView, "8-bit Huffman encoding in bytes: %d\n", bytes)
		}
		bytes, err = huffman.EncodedLen6Bit(t)
		if err != nil {
			fmt.Fprintf(textView, err.Error()+"\n")
		} else {
			fmt.Fprintf(textView, "6-bit Huffman encoding in bytes: %d\n", bytes)
		}
		bytes, err = huffman.EncodedLen5Bit(t)
		if err != nil {
			fmt.Fprintf(textView, err.Error()+"\n")
		} else {
			fmt.Fprintf(textView, "5-bit Huffman encoding in bytes: %d\n", bytes)
		}
	})
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter || key == tcell.KeyEscape {
			inputField.SetText("")
		}
	})
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).AddItem(textView, 0, 1, false).AddItem(inputField, 1, 0, true)
	app := tview.NewApplication()
	return app.SetRoot(flex, true).Run()
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
	os.Exit(2)
}

func main() {
	if len(os.Args) != 1 {
		usage()
	}
	if err := testEncodings(); err != nil {
		fatal(err)
	}
}
