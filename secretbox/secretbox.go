// Package secretbox implements a layer above NaCL's secretbox to encrypt and
// decrypt files.
package secretbox

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"syscall"

	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/terminal"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

// Seal encrypts the file plainFile with a passphrase read from stdin and
// stores the result in cryptFile (which must not exist):
//
// 1. Load plainFile as MSG.
//
// 2. Generate 32-byte SALT (for Argon2id) and 24-byte NONCE (for secretbox).
//
// 3. Derive 32-byte KEY (for secretbox) from passphrase with Argon2id using
// SALT (with time=1, memory=64MB, and threads=4).
//
// 4. Encrypt MSG to ENC with NaCL's secretbox.Seal using NONCE and KEY.
//
// 5. Save SALT|NONCE|ENC to cryptFile.
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
	derivedKey := argon2.IDKey(passphrase, salt[:], 1, 64*1024, 4, 32)
	copy(key[:], derivedKey)
	enc := secretbox.Seal(append(salt[:], nonce[:]...), msg, &nonce, &key)
	return ioutil.WriteFile(cryptFile, enc, 0644)
}

// Open decrypts the file cryptFile with a passphrase read from stdin and
// stores the result in plainFile (which must not exist):
//
// 1. Load cryptFile as BUF.
//
// 2. Split BUF into SALT|NONCE|ENC, where SALT is 32-byte, NONCE is 24-byte,
// and ENC is the remainder.
//
// 3. Derive 32-byte KEY (for secretbox) from passphrase with Argon2id using
// SALT (with time=1, memory=64MB, and threads=4).
//
// 4. Decrypt ENC to MSG with NaCL's secretbox.Open using NONCE and KEY.
//
// 5. Save MSG to plainFile.
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
	derivedKey := argon2.IDKey(passphrase, salt[:], 1, 64*1024, 4, 32)
	copy(key[:], derivedKey)
	msg, verify := secretbox.Open(nil, enc[56:], &nonce, &key)
	if !verify {
		return fmt.Errorf("cannot decrypt '%s'", plainFile)
	}
	return ioutil.WriteFile(plainFile, msg, 0644)
}
