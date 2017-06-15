// Copyright (c) 2017 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package textbuffer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloWorld(t *testing.T) {
	s := "Hello, 世界"
	tb, err := NewString(s)
	require.NoError(t, err)
	assert.Equal(t, 9, tb.LineLenRune(0))
	assert.Equal(t, 9, tb.LineLenChar(0))
	assert.Equal(t, 11, tb.LineLenCell(0))
	var b bytes.Buffer
	err = tb.Write(&b)
	if assert.NoError(t, err) {
		assert.Equal(t, s, b.String())
	}
}
