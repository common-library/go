package rsa_test

import (
	"sync"
	"testing"

	"github.com/common-library/go/security/crypto/rsa"
)

var onceForPublicKey sync.Once
var keyPairForPublicKey rsa.KeyPair

func publickey_instance(t *testing.T) (rsa.PrivateKey, rsa.PublicKey) {
	onceForPublicKey.Do(func() {
		keyPairForPublicKey = keypair_instance(t)
	})

	return keyPairForPublicKey.GetKeyPair()
}

func publickey_test(t *testing.T, privateKey rsa.PrivateKey, publicKey rsa.PublicKey) {
	const plaintext = "abcdefg12345"

	if ciphertext, err := publicKey.EncryptPKCS1v15(plaintext); err != nil {
		t.Fatal(err)
	} else if result, err := privateKey.DecryptPKCS1v15(ciphertext); err != nil {
		t.Fatal(err)
	} else if result != plaintext {
		t.Fatal("result != plaintext")
	}

	if ciphertext, err := publicKey.EncryptOAEP(plaintext); err != nil {
		t.Fatal(err)
	} else if result, err := privateKey.DecryptOAEP(ciphertext); err != nil {
		t.Fatal(err)
	} else if result != plaintext {
		t.Fatal("result != plaintext")
	}

}

func TestPublicKey_EncryptPKCS1v15(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_EncryptOAEP(t *testing.T) {
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

func TestPublicKey_GetPemPKCS1(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	if len(publicKey.GetPemPKCS1()) == 0 {
		t.Fatal("invalid privateKey.GetPemPKCS1()")
	}

	publickey_test(t, privateKey, publicKey)
}

func TestPublicKey_SetPemPKCS1(t *testing.T) {
	privateKey, publicKey := publickey_instance(t)

	if pemPKCS1 := publicKey.GetPemPKCS1(); len(pemPKCS1) == 0 {
		t.Fatal("invalid publicKey.GetPemPKCS1()")
	} else if err := publicKey.SetPemPKCS1(pemPKCS1); err != nil {
		t.Fatal(err)
	}

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
