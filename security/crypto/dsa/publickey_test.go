package dsa_test

import (
	"testing"

	"github.com/heaven-chp/common-library-go/security/crypto/dsa"
)

func publickey_instance(t *testing.T) (dsa.PrivateKey, dsa.PublicKey) {
	keyPair := keypair_instance(t)

	return keyPair.GetKeyPair()
}

func publickey_test(t *testing.T, privateKey dsa.PrivateKey, publicKey dsa.PublicKey) {
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

func TestPublicKey_GetPemAsn1(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	if pemAsn1, err := publicKey.GetPemAsn1(); err != nil {
		t.Log(pemAsn1)
		t.Fatal(err)
	}

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_SetPemAsn1(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	if pemAsn1, err := publicKey.GetPemAsn1(); err != nil {
		t.Log(pemAsn1)
		t.Fatal(err)
	} else if err := publicKey.SetPemAsn1(pemAsn1); err != nil {
		t.Log(pemAsn1)
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
