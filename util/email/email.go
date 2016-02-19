// Copyright (c) 2016 Frank Braun <frank@cryptogroup.net>
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Package email contains email related functions.
package email

import (
	"net/mail"
)

// IsValid reports whether email is a valid email address (user@domain).
func IsValid(email string) bool {
	a, err := mail.ParseAddress("TEST NAME <" + email + ">")
	if err != nil || a.Name != "TEST NAME" {
		return false
	}
	return true
}
