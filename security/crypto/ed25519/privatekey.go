// Package ed25519 provides ed25519 crypto related implementations.
package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)

// PrivateKey is struct that provides private key related methods.
type PrivateKey struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

// Sign is create a signature for message.
//
// ex) signature := privateKey.Sign(message)
func (this *PrivateKey) Sign(message string) []byte {
	return ed25519.Sign(this.privateKey, []byte(message))
}

// Verify is verifies the signature.
//
// ex) result := privateKey.Verify(message, signature)
func (this *PrivateKey) Verify(message string, signature []byte) bool {
	return ed25519.Verify(this.publicKey, []byte(message), signature)
}

// Get is to get a ed25519.PrivateKey.
//
// ex) key := privateKey.Get()
func (this *PrivateKey) Get() ed25519.PrivateKey {
	return this.privateKey
}

// Set is to set a ed25519.PrivateKey.
//
// ex) privateKey.Set(key)
func (this *PrivateKey) Set(privateKey ed25519.PrivateKey) {
	this.privateKey = privateKey
	this.publicKey = privateKey.Public().(ed25519.PublicKey)
}

// SetDefault is to set the primary key.
//
// ex) err := privateKey.SetDefault()
func (this *PrivateKey) SetDefault() error {
	if _, privateKey, err := ed25519.GenerateKey(rand.Reader); err != nil {
		return err
	} else {
		this.Set(privateKey)
		return nil
	}
}

// GetPemPKCS8 is to get a string in Pem/PKCS8 format.
//
// ex) pemPKCS8, err := privateKey.GetPemPKCS8()
func (this *PrivateKey) GetPemPKCS8() (string, error) {
	if blockBytes, err := x509.MarshalPKCS8PrivateKey(this.privateKey); err != nil {
		return "", err
	} else {
		return string(pem.EncodeToMemory(
			&pem.Block{
				Type:    "ED25519 PRIVATE KEY",
				Headers: nil,
				Bytes:   blockBytes,
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
		this.Set(key.(ed25519.PrivateKey))
		return nil
	}
}

// GetPublicKey is to get a PublicKey.
//
// ex) key := privateKey.GetPublicKey()
func (this *PrivateKey) GetPublicKey() PublicKey {
	publicKey := PublicKey{}

	publicKey.Set(this.publicKey)

	return publicKey
}
