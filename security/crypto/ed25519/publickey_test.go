package ed25519_test

import (
	"testing"

	"github.com/common-library/go/security/crypto/ed25519"
)

func publickey_instance(t *testing.T) (ed25519.PrivateKey, ed25519.PublicKey) {
	keyPair := keypair_instance(t)

	return keyPair.GetKeyPair()
}

func publickey_test(t *testing.T, privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) {
	const message = "abcdefg12345"

	signature := privateKey.Sign(message)
	if publicKey.Verify(message, signature) == false {
		t.Fatal("Verify is false")
	}
}

func TestPublicKey_Verify(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_Get(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	_ = publicKey.Get()

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_Set(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	publicKey.Set(publicKey.Get())

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_GetPemPKIX(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	if pemPKIX, err := publicKey.GetPemPKIX(); err != nil {
		t.Log(pemPKIX)
		t.Fatal(err)
	}

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_SetPemPKIX(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	if pemPKIX, err := publicKey.GetPemPKIX(); err != nil {
		t.Log(pemPKIX)
		t.Fatal(err)
	} else if err := publicKey.SetPemPKIX(pemPKIX); err != nil {
		t.Log(pemPKIX)
		t.Fatal(err)
	}

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_GetSsh(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	if key, err := publicKey.GetSsh(); err != nil {
		t.Log(key)
		t.Fatal(err)
	}

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_SetSsh(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	if key, err := publicKey.GetSsh(); err != nil {
		t.Log(key)
		t.Fatal(err)
	} else if err := publicKey.SetSsh(key); err != nil {
		t.Log(key)
		t.Fatal(err)
	}

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_GetSshPublicKey(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	if key, err := publicKey.GetSshPublicKey(); err != nil {
		t.Log(key)
		t.Fatal(err)
	}

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_SetSshPublicKey(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	if key, err := publicKey.GetSshPublicKey(); err != nil {
		t.Log(key)
		t.Fatal(err)
	} else if err := publicKey.SetSshPublicKey(key); err != nil {
		t.Log(key)
		t.Fatal(err)
	}

	publickey_test(t, privateKey, publicKey)
}
