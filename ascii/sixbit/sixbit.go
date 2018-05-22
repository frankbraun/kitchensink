// Package sixbit implements a 6-bit ASCII encoding.
package sixbit

import (
	"fmt"
)

// EncodedLen returns the encoded length (in bytes) of text in 6-bit ASCII.
func EncodedLen(text []byte) (int, error) {
	var chars int
	for _, c := range text {
		if 32 <= c && c <= 95 {
			chars++
		} else {
			return 0, fmt.Errorf("sixbit: cannot encode '%c' as 6-bit ASCII", c)
		}
	}
	bytes := chars * 6 / 8
	if chars*6%8 > 0 {
		return bytes + 1, nil
	}
	return bytes, nil
}

// Encode text in 6-bit ASCII (one byte per character).
func Encode(text []byte) ([]byte, error) {
	buf := make([]byte, len(text))
	for i, c := range text {
		if 32 <= c && c < 64 {
			buf[i] = c
		} else if 64 <= c && c <= 95 {
			buf[i] = c - 64
		} else {
			return nil, fmt.Errorf("sixbit: cannot encode '%c' as 6-bit ASCII", c)
		}
	}
	return buf, nil
}

// DecodeChar decodes a single 6-bit encode character.
func DecodeChar(c byte) (byte, error) {
	switch {
	case c < 32:
		return c + 64, nil
	case c < 64:
		return c, nil
	default:
		return 0, fmt.Errorf("sixbit: %d is not a 6-bit character", c)
	}
}
