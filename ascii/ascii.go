// Package ascii implements an efficient ASCII encoding.
package ascii

import (
	"bytes"
	"fmt"

	"github.com/frankbraun/kitchensink/ascii/fivebit"
	"github.com/frankbraun/kitchensink/ascii/sixbit"
	"github.com/icza/bitio"
)

// Encode an ASCII text efficiently (either as 5-bit or 6-bit ASCII).
// The text can be maximally 127 characters long.
func Encode(text []byte) ([]byte, error) {
	var buf bytes.Buffer
	var sixBit bool
	enc, err := fivebit.Encode(text)
	if err != nil {
		enc, err = sixbit.Encode(text)
		if err != nil {
			return nil, err
		}
		sixBit = true
	}
	l := len(enc)
	if l > 127 {
		return nil, fmt.Errorf("ascii: text has too many characters (%d > 127)", l)
	}
	w := bitio.NewWriter(&buf)
	// write bit indicating 5-bit for 6-bit encoding
	if err := w.WriteBool(sixBit); err != nil {
		return nil, err
	}
	// write length of encoding
	if err := w.WriteBits(uint64(l), 7); err != nil {
		return nil, err
	}
	// write encoding
	bits := byte(5)
	if sixBit {
		bits = 6
	}
	for _, c := range enc {
		if err := w.WriteBits(uint64(c), bits); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode an efficiently encoded ASCII text
func Decode(text []byte) ([]byte, error) {
	r := bitio.NewReader(bytes.NewBuffer(text))
	sixBit, err := r.ReadBool()
	if err != nil {
		return nil, err
	}
	l, err := r.ReadBits(7)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, l)
	for i := 0; i < int(l); i++ {
		var c byte
		if sixBit {
			u, err := r.ReadBits(6)
			if err != nil {
				return nil, err
			}
			c, err = sixbit.DecodeChar(byte(u))
			if err != nil {
				return nil, err
			}
		} else {
			u, err := r.ReadBits(5)
			if err != nil {
				return nil, err
			}
			c, err = fivebit.DecodeChar(byte(u))
			if err != nil {
				return nil, err
			}
		}
		buf[i] = c
	}
	return buf, nil
}
