package crypt_test

import (
	"testing"

	"github.com/devshansharma/tools/crypt"
)

func TestAESGCM(t *testing.T) {
	key, err := crypt.GenerateAES256Key()
	if err != nil {
		t.Fatal(err)
	}

	plaintext := "This is a secret message!"
	encrypted, err := crypt.EncryptAESGCM(key, []byte(plaintext))
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := crypt.DecryptAESGCM(key, encrypted)
	if err != nil {
		t.Fatal(err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original plaintext: %s != %s", decrypted, plaintext)
	}
}
