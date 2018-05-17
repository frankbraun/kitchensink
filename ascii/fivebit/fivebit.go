// Package fivebit implements a 5-bit ASCII encoding.
package fivebit

import (
	"bytes"
	"fmt"
)

// 64-90 "@ABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 32 ' '
// 33 '!'
// 44 ','
// 46 '.'
// 63 '?'
var charSet = []byte("@ABCDEFGHIJKLMNOPQRSTUVWXYZ !,.?")

// EncodedLen returns the encoded length (in bytes) of text in 5-bit ASCII.
func EncodedLen(text []byte) (int, error) {
	var chars int
	for _, c := range text {
		if bytes.ContainsRune(charSet, rune(c)) {
			chars++
		} else {
			return 0, fmt.Errorf("fivebit: cannot encode '%c' as 5-bit ASCII", c)
		}
	}
	bytes := chars * 5 / 8
	if chars*5%8 > 0 {
		return bytes + 1, nil
	}
	return bytes, nil
}

// Encode text in 5-bit ASCII (one byte per character).
func Encode(text []byte) ([]byte, error) {
	buf := make([]byte, len(text))
	for i, c := range text {
		idx := bytes.Index(charSet, []byte{c})
		if idx < 0 {
			return nil, fmt.Errorf("fivebit: cannot encode '%c' as 5-bit ASCII", c)
		}
		buf[i] = byte(idx)
	}
	return buf, nil
}

// DecodeChar decodes a single 5-bit encode character.
func DecodeChar(c byte) (byte, error) {
	switch {
	case c < 27:
		return c + 64, nil
	case c == 27:
		return ' ', nil
	case c == 28:
		return '!', nil
	case c == 29:
		return ',', nil
	case c == 30:
		return '.', nil
	case c == 31:
		return '?', nil
	default:
		return 0, fmt.Errorf("fivebit: %d is not a 5-bit character", c)
	}
}
