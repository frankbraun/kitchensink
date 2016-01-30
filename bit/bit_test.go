// Copyright (c) 2016 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package bit

import (
	"testing"
)

func TestCount(t *testing.T) {
	if Count(0) != 0 {
		t.Error("Count(0) should be 0")
	}
	if Count(1) != 1 {
		t.Error("Count(1) should be 1")
	}
	if Count(7) != 3 {
		t.Error("Count(7) should be 3")
	}
	if Count(8) != 1 {
		t.Error("Count(8) should be 1")
	}
	if Count(255) != 8 {
		t.Error("Count(255) should be 8")
	}
}
