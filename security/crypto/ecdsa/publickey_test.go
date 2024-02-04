package ecdsa_test

import (
	"testing"

	"github.com/heaven-chp/common-library-go/security/crypto/ecdsa"
)

func publickey_instance(t *testing.T) (ecdsa.PrivateKey, ecdsa.PublicKey) {
	keyPair := keypair_instance(t)

	return keyPair.GetKeyPair()
}

func publickey_test(t *testing.T, privateKey ecdsa.PrivateKey, publicKey ecdsa.PublicKey) {
	const message = "abcdefg12345"

	if signature, err := privateKey.Sign(message); err != nil {
		t.Fatal(err)
	} else if publicKey.Verify(message, signature) == false {
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
