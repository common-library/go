package dsa_test

import (
	crypto_dsa "crypto/dsa"
	"testing"

	"github.com/common-library/go/security/crypto/dsa"
)

func privatekey_instance(t *testing.T) dsa.PrivateKey {
	t.Parallel()

	privateKey := dsa.PrivateKey{}

	if err := privateKey.SetSizes(crypto_dsa.L1024N160); err != nil {
		t.Fatal(err)
	}

	return privateKey
}

func privatekey_test(t *testing.T, privateKey dsa.PrivateKey) {
	const message = "abcdefg12345"

	if signature, err := privateKey.Sign(message); err != nil {
		t.Fatal(err)
	} else if privateKey.Verify(message, signature) == false {
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
}

func TestPrivateKey_Set(t *testing.T) {
	privateKey := privatekey_instance(t)

	privateKey.Set(privateKey.Get())
}

func TestPrivateKey_SetSizes(t *testing.T) {
	privateKey := dsa.PrivateKey{}

	if err := privateKey.SetSizes(crypto_dsa.L1024N160); err != nil {
		t.Fatal(err)
	}
}

func TestPrivateKey_GetPemAsn1(t *testing.T) {
	privateKey := privatekey_instance(t)

	if pemAsn1, err := privateKey.GetPemAsn1(); err != nil {
		t.Log(pemAsn1)
		t.Fatal(err)
	}
}

func TestPrivateKey_SetPemAsn1(t *testing.T) {
	privateKey := privatekey_instance(t)

	if pemAsn1, err := privateKey.GetPemAsn1(); err != nil {
		t.Log(pemAsn1)
		t.Fatal(err)
	} else if err := privateKey.SetPemAsn1(pemAsn1); err != nil {
		t.Log(pemAsn1)
		t.Fatal(err)
	}
}

func TestPrivateKey_GetPublicKey(t *testing.T) {
	privateKey := privatekey_instance(t)

	_ = privateKey.GetPublicKey()
}
