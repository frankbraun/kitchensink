// Package fivebit implements a 5-bit ASCII encoding.
package fivebit

import (
	"bytes"
	"fmt"
	"strings"
)

// 64-90 "@ABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 32 ' '
// 33 '!'
// 44 ','
// 46 '.'
// 63 '?'
const charSet = "@ABCDEFGHIJKLMNOPQRSTUVWXYZ !,.?"

// EncodedLen returns the encoded length (in bytes) of text in 5-bit ASCII.
func EncodedLen(text string) (int, error) {
	var chars int
	for _, c := range []byte(text) {
		if strings.ContainsRune(charSet, rune(c)) {
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
func Encode(text string) ([]byte, error) {
	buf := make([]byte, len(text))
	for i, c := range []byte(text) {
		idx := bytes.Index([]byte(charSet), []byte{c})
		if idx < 0 {
			return nil, fmt.Errorf("fivebit: cannot encode '%c' as 5-bit ASCII", c)
		}
		buf[i] = byte(idx)
	}
	return buf, nil
}
