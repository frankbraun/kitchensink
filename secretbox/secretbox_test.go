package secretbox

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSealAndOpen(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "secretbox_test")
	if err != nil {
		t.Fatalf("ioutil.TempDir() failed: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	plainFileTest := filepath.Join("testdata", "test.txt")
	plainFileTemp := filepath.Join(tmpdir, "test.txt")
	cryptFileTemp := filepath.Join(tmpdir, "test.bin")

	testPass = "test"

	msg, err := ioutil.ReadFile(plainFileTest)
	if err != nil {
		t.Fatalf("ioutil.ReadFile(%s) failed: %v", plainFileTest, err)
	}

	// encrypt plain text file
	err = Seal(plainFileTest, cryptFileTemp)
	if err != nil {
		t.Fatalf("Seal() failed: %v", err)
	}
	tmp, err := ioutil.ReadFile(cryptFileTemp)
	if err != nil {
		t.Fatalf("ioutil.ReadFile(%s) failed: %v", cryptFileTemp, err)
	}

	// encrypting to existing file should file
	err = Seal(plainFileTest, cryptFileTemp)
	if err == nil {
		t.Fatal("Seal() should fail")
	}

	// decrypt encrypted file and compare results
	err = Open(cryptFileTemp, plainFileTemp)
	if err != nil {
		t.Fatalf("Seal() failed: %v", err)
	}
	tmp, err = ioutil.ReadFile(plainFileTemp)
	if err != nil {
		t.Fatalf("ioutil.ReadFile(%s) failed: %v", plainFileTemp, err)
	}
	if !bytes.Equal(tmp, msg) {
		t.Fatal("plainFileTemp != plainFileTest")
	}

	// decrypting to existing file should fail
	err = Open(cryptFileTemp, plainFileTemp)
	if err == nil {
		t.Fatal("Open() shoudl fail")
	}

	if err := os.Remove(plainFileTemp); err != nil {
		t.Fatalf("os.Remove(%s) failed: %v", plainFileTemp, err)
	}

	// decrypting with wrong passphrase should fail
	testPass = "fail"
	err = Open(cryptFileTemp, plainFileTemp)
	if err == nil {
		t.Fatalf("Open() with wrong passphrase should fail")
	}
}
