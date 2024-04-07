package ed25519_test

import (
	"testing"

	"github.com/common-library/go/security/crypto/ed25519"
)

func privatekey_instance(t *testing.T) ed25519.PrivateKey {
	privateKey := ed25519.PrivateKey{}

	if err := privateKey.SetDefault(); err != nil {
		t.Fatal(err)
	}

	return privateKey
}

func privatekey_test(t *testing.T, privateKey ed25519.PrivateKey) {
	const message = "abcdefg12345"

	signature := privateKey.Sign(message)
	if privateKey.Verify(message, signature) == false {
		t.Fatal("Verify is false")
	}
}

func TestPrivateKey_Sign(t *testing.T) {
	privateKey := privatekey_instance(t)

	privatekey_test(t, privateKey)
}

func TestPrivateKey_Verify(t *testing.T) {
	privateKey := privatekey_instance(t)

	privatekey_test(t, privateKey)
}

func TestPrivateKey_Get(t *testing.T) {
	privateKey := privatekey_instance(t)

	_ = privateKey.Get()

	privatekey_test(t, privateKey)
}

func TestPrivateKey_Set(t *testing.T) {
	privateKey := privatekey_instance(t)

	privateKey.Set(privateKey.Get())

	privatekey_test(t, privateKey)
}

func TestPrivateKey_SetDefault(t *testing.T) {
	privateKey := ed25519.PrivateKey{}

	if err := privateKey.SetDefault(); err != nil {
		t.Fatal(err)
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_GetPemPKCS8(t *testing.T) {
	privateKey := privatekey_instance(t)

	if pemAsn1, err := privateKey.GetPemPKCS8(); err != nil {
		t.Log(pemAsn1)
		t.Fatal(err)
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_SetPemPKCS8(t *testing.T) {
	privateKey := privatekey_instance(t)

	if pemAsn1, err := privateKey.GetPemPKCS8(); err != nil {
		t.Log(pemAsn1)
		t.Fatal(err)
	} else if err := privateKey.SetPemPKCS8(pemAsn1); err != nil {
		t.Log(pemAsn1)
		t.Fatal(err)
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_GetPublicKey(t *testing.T) {
	privateKey := privatekey_instance(t)

	_ = privateKey.GetPublicKey()

	privatekey_test(t, privateKey)
}
