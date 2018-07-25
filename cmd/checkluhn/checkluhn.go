// checkluhn runs the Luhn algorithm on a purported credit card number.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func checkLuhn(purportedCC string) bool {
	var sum = 0
	var nDigits = len(purportedCC)
	var parity = nDigits % 2

	for i := 0; i < nDigits; i++ {
		var digit = int(purportedCC[i] - 48)
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s purported_cc\n", os.Args[0])
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	valid := checkLuhn(strings.Replace(os.Args[1], " ", "", -1))
	if !valid {
		fatal(errors.New("number not valid"))
	}
	fmt.Println("number valid")
}
