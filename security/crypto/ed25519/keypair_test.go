package ed25519_test

import (
	"testing"

	"github.com/common-library/go/security/crypto/ed25519"
)

func keypair_instance(t *testing.T) ed25519.KeyPair {
	keyPair := ed25519.KeyPair{}

	if err := keyPair.Generate(); err != nil {
		t.Fatal(err)
	}

	return keyPair
}

func keypair_test(t *testing.T, keyPair ed25519.KeyPair) {
	const message = "abcdefg12345"

	signature := keyPair.Sign(message)
	if keyPair.Verify(message, signature) == false {
		t.Fatal("Verify is false")
	}
}

func TestKeyPair_Generate(t *testing.T) {
	keyPair := ed25519.KeyPair{}

	if err := keyPair.Generate(); err != nil {
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
