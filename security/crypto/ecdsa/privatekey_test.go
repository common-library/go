package ecdsa_test

import (
	"crypto/elliptic"
	"testing"

	"github.com/heaven-chp/common-library-go/security/crypto/ecdsa"
)

func privatekey_instance(t *testing.T) ecdsa.PrivateKey {
	privateKey := ecdsa.PrivateKey{}

	if err := privateKey.SetCurve(elliptic.P384()); err != nil {
		t.Fatal(err)
	}

	return privateKey
}

func privatekey_test(t *testing.T, privateKey ecdsa.PrivateKey) {
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

	privatekey_test(t, privateKey)
}

func TestPrivateKey_Set(t *testing.T) {
	privateKey := privatekey_instance(t)

	privateKey.Set(privateKey.Get())

	privatekey_test(t, privateKey)
}

func TestPrivateKey_GetCurve(t *testing.T) {
	privateKey := privatekey_instance(t)

	if privateKey.GetCurve() != elliptic.P384() {
		t.Fatal("privateKey.GetCurve() != elliptic.P384()")
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_SetCurve(t *testing.T) {
	privateKey := ecdsa.PrivateKey{}

	if err := privateKey.SetCurve(elliptic.P384()); err != nil {
		t.Fatal(err)
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_GetPemEC(t *testing.T) {
	privateKey := privatekey_instance(t)

	if pemEC, err := privateKey.GetPemEC(); err != nil {
		t.Log(pemEC)
		t.Fatal(err)
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_SetPemEC(t *testing.T) {
	privateKey := privatekey_instance(t)

	if pemEC, err := privateKey.GetPemEC(); err != nil {
		t.Log(pemEC)
		t.Fatal(err)
	} else if err := privateKey.SetPemEC(pemEC); err != nil {
		t.Log(pemEC)
		t.Fatal(err)
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_GetPublicKey(t *testing.T) {
	privateKey := privatekey_instance(t)

	_ = privateKey.GetPublicKey()

	privatekey_test(t, privateKey)
}
