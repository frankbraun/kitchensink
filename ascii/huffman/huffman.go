// Package huffman implements Huffman encodings for various ASCII encodings.
package huffman

import (
	"bytes"
	"io"

	bitHufio "github.com/frankbraun/huffman/hufio"
	"github.com/frankbraun/kitchensink/ascii/fivebit"
	"github.com/frankbraun/kitchensink/ascii/sixbit"
	"github.com/icza/huffman/hufio"
)

// EncodedLen8Bit returns the Huffman encoded length (in bytes) of text in
// 8-bit ASCII.
func EncodedLen8Bit(text string) (int, error) {
	var buf bytes.Buffer
	w := hufio.NewWriter(&buf)
	if _, err := io.WriteString(w, text); err != nil {
		return 0, err
	}
	if err := w.Close(); err != nil {
		return 0, err
	}
	return buf.Len(), nil
}

// EncodedLen6Bit returns the Huffman encoded length (in bytes) of text in
// 6-bit ASCII.
func EncodedLen6Bit(text string) (int, error) {
	enc, err := sixbit.Encode(text)
	if err != nil {
		return 0, err
	}
	var buf bytes.Buffer
	w := bitHufio.NewWriterOptions(&buf, &bitHufio.Options{BitWidth: 6})
	if _, err := w.Write(enc); err != nil {
		return 0, err
	}
	if err := w.Close(); err != nil {
		return 0, err
	}
	return buf.Len(), nil
}

// EncodedLen5Bit returns the Huffman encoded length (in bytes) of text in
// 5-bit ASCII.
func EncodedLen5Bit(text string) (int, error) {
	enc, err := fivebit.Encode(text)
	if err != nil {
		return 0, err
	}
	var buf bytes.Buffer
	w := bitHufio.NewWriterOptions(&buf, &bitHufio.Options{BitWidth: 5})
	if _, err := w.Write(enc); err != nil {
		return 0, err
	}
	if err := w.Close(); err != nil {
		return 0, err
	}
	return buf.Len(), nil
}
