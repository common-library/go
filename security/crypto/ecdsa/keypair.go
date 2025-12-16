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
func (kp *KeyPair) Generate(curve elliptic.Curve) error {
	if err := kp.privateKey.SetCurve(curve); err != nil {
		return err
	} else {
		kp.publicKey = kp.privateKey.GetPublicKey()
		return nil
	}
}

// Sign is create a signature for message.
//
// ex) signature, err := keyPair.Sign(message)
func (kp *KeyPair) Sign(message string) (Signature, error) {
	return kp.privateKey.Sign(message)
}

// Verify is verifies the signature.
//
// ex) result := keyPair.Verify(message, signature)
func (kp *KeyPair) Verify(message string, signature Signature) bool {
	return kp.publicKey.Verify(message, signature)
}

// GetKeyPair is to get a key pair.
//
// ex) privateKey, publicKey := keyPair.GetKeyPair()
func (kp *KeyPair) GetKeyPair() (privateKey PrivateKey, publicKey PublicKey) {
	return kp.privateKey, kp.publicKey
}

// SetKeyPair is to set a key pair.
//
// ex) keyPair.SetKeyPair(privateKey, publicKey)
func (kp *KeyPair) SetKeyPair(privateKey PrivateKey, publicKey PublicKey) {
	kp.privateKey = privateKey
	kp.publicKey = publicKey
}
