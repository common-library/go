package ecdsa_test

import (
	"crypto/elliptic"
	"testing"

	"github.com/common-library/go/security/crypto/ecdsa"
)

func keypair_instance(t *testing.T) ecdsa.KeyPair {
	keyPair := ecdsa.KeyPair{}

	if err := keyPair.Generate(elliptic.P384()); err != nil {
		t.Fatal(err)
	}

	return keyPair
}

func keypair_test(t *testing.T, keyPair ecdsa.KeyPair) {
	const message = "abcdefg12345"

	if signature, err := keyPair.Sign(message); err != nil {
		t.Fatal(err)
	} else if keyPair.Verify(message, signature) == false {
		t.Fatal("Verify is false")
	}
}

func TestKeyPair_Generate(t *testing.T) {
	keyPair := ecdsa.KeyPair{}

	if err := keyPair.Generate(elliptic.P384()); err != nil {
		t.Fatal(err)
	}

	keypair_test(t, keyPair)
}

func TestKeyPair_Sign(t *testing.T) {
	keyPair := keypair_instance(t)

	keypair_test(t, keyPair)
}

func TestKeyPair_Verify(t *testing.T) {
	keyPair := keypair_instance(t)

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
