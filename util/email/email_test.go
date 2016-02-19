// Copyright (c) 2016 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package email

import (
	"strconv"
	"testing"
)

func TestIsValid(t *testing.T) {
	var emailTests = []struct {
		email string
		valid bool
	}{
		{"john.doe@example.com", true},
		{"john@example.com", true},
		{"@example.com", false},
		{"john.doeexample.com", false},
		{"johnexamplecom", false},
		{"", false},
		{"John <john.doe@example.com>", false}, // IsValid() expects user@doman
		{"John Doe <john.doe@example.com>", false},
	}
	for _, test := range emailTests {
		if IsValid(test.email) != test.valid {
			t.Errorf("IsValid(%s) != %s", test.email,
				strconv.FormatBool(test.valid))
		}
	}
}
