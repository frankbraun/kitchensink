// Package secretbox implements a layer above NaCL's secretbox to encrypt and
// decrypt files.
package secretbox

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"runtime"
	"syscall"

	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/terminal"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

// maxThreads returns the maximum number of usable threads as uint8.
func maxThreads() uint8 {
	var threads uint8
	numCPU := runtime.NumCPU()
	if numCPU > math.MaxUint8 {
		threads = math.MaxUint8
	} else {
		threads = uint8(numCPU)
	}
	return threads
}

// Seal encrypts the file plainFile with a passphrase read from stdin and
// stores the result in cryptFile (which must not exist).
func Seal(plainFile, cryptFile string) error {
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
	derivedKey := argon2.IDKey(passphrase, salt[:], 1, 64*1024, maxThreads(), 32)
	copy(key[:], derivedKey)
	enc := secretbox.Seal(append(salt[:], nonce[:]...), msg, &nonce, &key)
	return ioutil.WriteFile(cryptFile, enc, 0644)
}

// Open decrypts the file cryptFile with a passphrase read from stdin and
// stores the result in plainFile (which must not exist).
func Open(cryptFile, plainFile string) error {
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
	derivedKey := argon2.IDKey(passphrase, salt[:], 1, 64*1024, maxThreads(), 32)
	copy(key[:], derivedKey)
	msg, verify := secretbox.Open(nil, enc[56:], &nonce, &key)
	if !verify {
		return fmt.Errorf("cannot decrypt '%s'", plainFile)
	}
	return ioutil.WriteFile(plainFile, msg, 0644)
}
