package dsa_test

import (
	//lint:ignore SA1019 DSA is deprecated but kept for compatibility testing
	crypto_dsa "crypto/dsa"
	"testing"

	"github.com/common-library/go/security/crypto/dsa"
)

func keypair_instance(t *testing.T) dsa.KeyPair {
	t.Parallel()

	keyPair := dsa.KeyPair{}

	if err := keyPair.Generate(crypto_dsa.L1024N160); err != nil {
		t.Fatal(err)
	}

	return keyPair
}

func keypair_test(t *testing.T, keyPair dsa.KeyPair) {
	const message = "abcdefg12345"

	if signature, err := keyPair.Sign(message); err != nil {
		t.Fatal(err)
	} else if keyPair.Verify(message, signature) == false {
		t.Fatal("Verify is false")
	}
}

func TestKeyPair_Generate(t *testing.T) {
	keyPair := dsa.KeyPair{}

	if err := keyPair.Generate(crypto_dsa.L1024N160); err != nil {
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
