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
