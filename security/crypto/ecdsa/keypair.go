// Package ecdsa provides ecdsa crypto related implementations.
package ecdsa

import (
	"crypto/elliptic"
)

// KeyPair is struct that provides key pair related methods.
type KeyPair struct {
	privateKey PrivateKey
	publicKey  PublicKey
}

// Generate is to generate a key pair.
//
// ex) err := keyPair.Generate(elliptic.P384())
func (this *KeyPair) Generate(curve elliptic.Curve) error {
	if err := this.privateKey.SetCurve(curve); err != nil {
		return err
	} else {
		this.publicKey = this.privateKey.GetPublicKey()
		return nil
	}
}

// Sign is create a signature for message.
//
// ex) signature, err := keyPair.Sign(message)
func (this *KeyPair) Sign(message string) (Signature, error) {
	return this.privateKey.Sign(message)
}

// Verify is verifies the signature.
//
// ex) result := keyPair.Verify(message, signature)
func (this *KeyPair) Verify(message string, signature Signature) bool {
	return this.publicKey.Verify(message, signature)
}

// GetKeyPair is to get a key pair.
//
// ex) privateKey, publicKey := keyPair.GetKeyPair()
func (this *KeyPair) GetKeyPair() (privateKey PrivateKey, publicKey PublicKey) {
	return this.privateKey, this.publicKey
}

// SetKeyPair is to set a key pair.
//
// ex) keyPair.SetKeyPair(privateKey, publicKey)
func (this *KeyPair) SetKeyPair(privateKey PrivateKey, publicKey PublicKey) {
	this.privateKey = privateKey
	this.publicKey = publicKey
}
