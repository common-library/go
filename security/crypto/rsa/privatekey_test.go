package rsa_test

import (
	"sync"
	"testing"

	"github.com/heaven-chp/common-library-go/security/crypto/rsa"
)

var onceForPrivateKey sync.Once
var PRIVATEKEY rsa.PrivateKey

func privatekey_instance(t *testing.T) rsa.PrivateKey {
	onceForPrivateKey.Do(func() {
		if err := PRIVATEKEY.SetBits(4096); err != nil {
			t.Fatal(err)
		}
	})

	privateKey := rsa.PrivateKey{}
	privateKey.Set(PRIVATEKEY.Get())
	return privateKey
}

func privatekey_test(t *testing.T, privateKey rsa.PrivateKey) {
	const plaintext = "abcdefg12345"

	if ciphertext, err := privateKey.EncryptPKCS1v15(plaintext); err != nil {
		t.Fatal(err)
	} else if result, err := privateKey.DecryptPKCS1v15(ciphertext); err != nil {
		t.Fatal(err)
	} else if result != plaintext {
		t.Fatal("result != plaintext")
	}

	if ciphertext, err := privateKey.EncryptOAEP(plaintext); err != nil {
		t.Fatal(err)
	} else if result, err := privateKey.DecryptOAEP(ciphertext); err != nil {
		t.Fatal(err)
	} else if result != plaintext {
		t.Fatal("result != plaintext")
	}
}

func TestPrivateKey_EncryptPKCS1v15(t *testing.T) {
	privateKey := privatekey_instance(t)

	privatekey_test(t, privateKey)
}

func TestPrivateKey_DecryptPKCS1v15(t *testing.T) {
	privateKey := privatekey_instance(t)

	privatekey_test(t, privateKey)
}

func TestPrivateKey_EncryptOAEP(t *testing.T) {
	privateKey := privatekey_instance(t)

	privatekey_test(t, privateKey)
}

func TestPrivateKey_DecryptOAEP(t *testing.T) {
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

func TestPrivateKey_SetBits(t *testing.T) {
	privateKey := rsa.PrivateKey{}

	if err := privateKey.SetBits(4096); err != nil {
		t.Fatal(err)
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_GetPemPKCS1(t *testing.T) {
	privateKey := privatekey_instance(t)

	if len(privateKey.GetPemPKCS1()) == 0 {
		t.Fatal("invalid privateKey.GetPemPKCS1()")
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_SetPemPKCS1(t *testing.T) {
	privateKey := privatekey_instance(t)

	if pemPKCS1 := privateKey.GetPemPKCS1(); len(pemPKCS1) == 0 {
		t.Fatal("invalid privateKey.GetPemPKCS1()")
	} else if err := privateKey.SetPemPKCS1(pemPKCS1); err != nil {
		t.Fatal(err)
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_GetPemPKCS8(t *testing.T) {
	privateKey := privatekey_instance(t)

	if pemPKCS8, err := privateKey.GetPemPKCS8(); err != nil {
		t.Log(pemPKCS8)
		t.Fatal(err)
	}

	privatekey_test(t, privateKey)
}

func TestPrivateKey_SetPemPKCS8(t *testing.T) {
	privateKey := privatekey_instance(t)
	if pemPKCS8, err := privateKey.GetPemPKCS8(); err != nil {
		t.Log(pemPKCS8)
		t.Fatal(err)
	} else if err := privateKey.SetPemPKCS8(pemPKCS8); err != nil {
		t.Log(pemPKCS8)
		t.Fatal(err)
	}

	privatekey_test(t, privateKey)
}
