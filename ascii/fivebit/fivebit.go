package fivebit

import (
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

// EncodedLen returns the encoded length in bytes of text in 5-bit ASCII.
func EncodedLen(text string) (int, error) {
	var chars int
	for _, c := range text {
		if strings.ContainsRune(charSet, c) {
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
