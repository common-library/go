// Package rsa provides rsa crypto related implementations.
package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

// PrivateKey is struct that provides private key related methods.
type PrivateKey struct {
	privateKey *rsa.PrivateKey
}

// EncryptPKCS1v15 is encrypt plaintext.
//
// ex) ciphertext, err := privateKey.EncryptPKCS1v15(plaintext)
func (this *PrivateKey) EncryptPKCS1v15(plaintext string) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, &this.privateKey.PublicKey, []byte(plaintext))
}

// DecryptPKCS1v15 is decrypt ciphertext.
//
// ex) plaintext, err := privateKey.DecryptPKCS1v15(ciphertext)
func (this *PrivateKey) DecryptPKCS1v15(ciphertext []byte) (string, error) {
	if plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, this.privateKey, ciphertext); err != nil {
		return "", err
	} else {
		return string(plaintext), nil
	}
}

// EncryptOAEP is encrypt plaintext.
//
// ex) ciphertext, err := privateKey.EncryptOAEP(plaintext)
func (this *PrivateKey) EncryptOAEP(plaintext string) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, &this.privateKey.PublicKey, []byte(plaintext), nil)
}

// DecryptOAEP is decrypt ciphertext.
//
// ex) plaintext, err := privateKey.DecryptOAEP(ciphertext)
func (this *PrivateKey) DecryptOAEP(ciphertext []byte) (string, error) {
	if plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, this.privateKey, ciphertext, nil); err != nil {
		return "", err
	} else {
		return string(plaintext), nil
	}
}

// Get is to get a *rsa.PrivateKey.
//
// ex) key := privateKey.Get()
func (this *PrivateKey) Get() *rsa.PrivateKey {
	return this.privateKey
}

// Set is to set a *rsa.PrivateKey.
//
// ex) privateKey.Set(key)
func (this *PrivateKey) Set(privateKey *rsa.PrivateKey) {
	this.privateKey = privateKey
}

// SetBits is to set the primary key using bits.
//
// ex) err := privateKey.SetBits(4096)
func (this *PrivateKey) SetBits(bits int) error {
	if privateKey, err := rsa.GenerateKey(rand.Reader, bits); err != nil {
		return err
	} else {
		this.Set(privateKey)
		return nil
	}
}

// GetPemPKCS1 is to get a string in Pem/PKCS1 format.
//
// ex) pemPKCS1 := privateKey.GetPemPKCS1()
func (this *PrivateKey) GetPemPKCS1() string {
	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:    "RSA PRIVATE KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PrivateKey(this.privateKey),
		}))
}

// SetPemPKCS1 is to set the primary key using a string in Pem/PKCS1 format.
//
// ex) err := privateKey.SetPemPKCS1(pemPKCS1)
func (this *PrivateKey) SetPemPKCS1(pemPKCS1 string) error {
	block, _ := pem.Decode([]byte(pemPKCS1))

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		return err
	} else {
		this.Set(key)
		return nil
	}
}

// GetPemPKCS8 is to get a string in Pem/PKCS8 format.
//
// ex) pemPKCS8, err := privateKey.GetPemPKCS8()
func (this *PrivateKey) GetPemPKCS8() (string, error) {
	if bytes, err := x509.MarshalPKCS8PrivateKey(this.privateKey); err != nil {
		return "", err
	} else {
		return string(pem.EncodeToMemory(
			&pem.Block{
				Type:    "PRIVATE KEY",
				Headers: nil,
				Bytes:   bytes,
			})), nil
	}
}

// SetPemPKCS8 is to set the primary key using a string in Pem/PKCS8 format.
//
// ex) err := privateKey.SetPemPKCS8(pemPKCS8)
func (this *PrivateKey) SetPemPKCS8(pemPKCS8 string) error {
	block, _ := pem.Decode([]byte(pemPKCS8))

	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return err
	} else {
		this.Set(key.(*rsa.PrivateKey))
		return nil
	}

}

// GetPublicKey is to get a PublicKey.
//
// ex) key := privateKey.GetPublicKey()
func (this *PrivateKey) GetPublicKey() PublicKey {
	publicKey := PublicKey{}

	publicKey.Set(this.privateKey.PublicKey)

	return publicKey
}
