// secretbox is a simple tool to encrypt and decrypt files with NaCL's secretbox.
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/terminal"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

func seal(plainFile, cryptFile string) error {
	exists, err := file.Exists(cryptFile)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("file '%s' exists already", cryptFile)
	}
	passphrase, err := terminal.ReadPassphrase(syscall.Stdin, true)
	if err != nil {
		return err
	}
	msg, err := ioutil.ReadFile(plainFile)
	if err != nil {
		return err
	}
	var (
		salt  [32]byte
		nonce [24]byte
		key   [32]byte
	)
	if _, err := io.ReadFull(rand.Reader, salt[:]); err != nil {
		return err
	}
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return err
	}
	derivedKey := argon2.IDKey(passphrase, salt[:], 1, 64*1024, 4, 32)
	copy(key[:], derivedKey)
	enc := secretbox.Seal(append(salt[:], nonce[:]...), msg, &nonce, &key)
	return ioutil.WriteFile(cryptFile, enc, 0644)
}

func open(cryptFile, plainFile string) error {
	exists, err := file.Exists(plainFile)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("file '%s' exists already", plainFile)
	}
	passphrase, err := terminal.ReadPassphrase(syscall.Stdin, false)
	if err != nil {
		return err
	}
	enc, err := ioutil.ReadFile(cryptFile)
	if err != nil {
		return err
	}
	var (
		salt  [32]byte
		nonce [24]byte
		key   [32]byte
	)
	copy(salt[:], enc[:32])
	copy(nonce[:], enc[32:56])
	derivedKey := argon2.IDKey(passphrase, salt[:], 1, 64*1024, 4, 32)
	copy(key[:], derivedKey)
	msg, verify := secretbox.Open(nil, enc[56:], &nonce, &key)
	if !verify {
		return fmt.Errorf("cannot decrypt '%s'", plainFile)
	}
	return ioutil.WriteFile(plainFile, msg, 0644)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%s: error: %s\n", os.Args[0], err)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s seal plain_file crypt_file\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "       %s open crypt_file plain_file\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()
	if flag.NArg() != 3 {
		usage()
	}
	switch flag.Arg(0) {
	case "seal":
		if err := seal(flag.Arg(1), flag.Arg(2)); err != nil {
			fatal(err)
		}
	case "open":
		if err := open(flag.Arg(1), flag.Arg(2)); err != nil {
			fatal(err)
		}
	default:
		usage()
	}
}
