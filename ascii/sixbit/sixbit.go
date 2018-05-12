package sixbit

import (
	"fmt"
)

// EncodedLen returns the encoded length in bytes of text in 6-bit ASCII.
func EncodedLen(text string) (int, error) {
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
