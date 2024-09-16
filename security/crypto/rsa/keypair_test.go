package rsa_test

import (
	"sync"
	"testing"

	"github.com/common-library/go/security/crypto/rsa"
)

var onceForKeyPair sync.Once
var KEYPAIR rsa.KeyPair

func keypair_instance(t *testing.T) rsa.KeyPair {
	t.Parallel()

	onceForKeyPair.Do(func() {
		if err := KEYPAIR.Generate(4096); err != nil {
			t.Fatal(err)
		}
	})

	keyPair := rsa.KeyPair{}
	keyPair.SetKeyPair(KEYPAIR.GetKeyPair())
	return keyPair
}

func keypair_test(t *testing.T, keyPair rsa.KeyPair) {
	const plaintext = "abcdefg12345"

	if ciphertext, err := keyPair.EncryptPKCS1v15(plaintext); err != nil {
		t.Fatal(err)
	} else if result, err := keyPair.DecryptPKCS1v15(ciphertext); err != nil {
		t.Fatal(err)
	} else if result != plaintext {
		t.Fatal("result != plaintext")
	}

	if ciphertext, err := keyPair.EncryptOAEP(plaintext); err != nil {
		t.Fatal(err)
	} else if result, err := keyPair.DecryptOAEP(ciphertext); err != nil {
		t.Fatal(err)
	} else if result != plaintext {
		t.Fatal("result != plaintext")
	}

}

func TestKeyPair_Generate(t *testing.T) {
	t.Parallel()

	keyPair := rsa.KeyPair{}

	if err := keyPair.Generate(4096); err != nil {
		t.Fatal(err)
	}

	keypair_test(t, keyPair)
}

func TestKeyPair_GetKeyPair(t *testing.T) {
	keyPair := keypair_instance(t)

	_, _ = keyPair.GetKeyPair()
}

func TestKeyPair_SetKeyPair(t *testing.T) {
	keyPair := keypair_instance(t)

	keyPair.SetKeyPair(keyPair.GetKeyPair())
}
