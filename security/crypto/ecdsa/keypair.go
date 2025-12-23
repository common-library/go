// Package ecdsa provides ECDSA digital signature cryptography.
//
// This package implements Elliptic Curve Digital Signature Algorithm (ECDSA)
// with support for multiple NIST curves, signature generation and verification,
// and PEM format key storage.
//
// # Features
//
//   - ECDSA key pair generation (P-256, P-384, P-521 curves)
//   - Digital signature creation and verification
//   - PEM/PKCS8/PKIX format support
//   - SSH public key format conversion
//   - Smaller key sizes than RSA
//
// # Basic Example
//
//	keyPair := &ecdsa.KeyPair{}
//	err := keyPair.Generate(elliptic.P256())
//	signature, _ := keyPair.Sign("message")
//	valid := keyPair.Verify("message", signature)
package ecdsa

import (
	"crypto/elliptic"
)

// KeyPair provides ECDSA key pair operations.
type KeyPair struct {
	privateKey PrivateKey
	publicKey  PublicKey
}

// Generate creates a new ECDSA key pair.
//
// # Parameters
//
//   - curve: Elliptic curve (P256, P384, or P521)
//
// # Returns
//
//   - error: Error if generation fails, nil on success
//
// # Examples
//
// P-256 curve (recommended for most applications):
//
//	keyPair := &ecdsa.KeyPair{}
//	err := keyPair.Generate(elliptic.P256())
//
// P-384 curve (higher security):
//
//	err := keyPair.Generate(elliptic.P384())
func (kp *KeyPair) Generate(curve elliptic.Curve) error {
	if err := kp.privateKey.SetCurve(curve); err != nil {
		return err
	} else {
		kp.publicKey = kp.privateKey.GetPublicKey()
		return nil
	}
}

// Sign creates a digital signature for the message.
//
// # Parameters
//
//   - message: Text message to sign
//
// # Returns
//
//   - Signature: ECDSA signature (R and S components)
//   - error: Error if signing fails, nil on success
//
// # Examples
//
//	signature, err := keyPair.Sign("Hello, World!")
func (kp *KeyPair) Sign(message string) (Signature, error) {
	return kp.privateKey.Sign(message)
}

// Verify verifies a digital signature.
//
// # Parameters
//
//   - message: Original text message
//   - signature: ECDSA signature to verify
//
// # Returns
//
//   - bool: true if signature is valid, false otherwise
//
// # Examples
//
//	signature, _ := keyPair.Sign("message")
//	valid := keyPair.Verify("message", signature)
func (kp *KeyPair) Verify(message string, signature Signature) bool {
	return kp.publicKey.Verify(message, signature)
}

// GetKeyPair retrieves the private and public keys.
//
// # Returns
//
//   - privateKey: ECDSA private key
//   - publicKey: ECDSA public key
//
// # Examples
//
//	privateKey, publicKey := keyPair.GetKeyPair()
func (kp *KeyPair) GetKeyPair() (privateKey PrivateKey, publicKey PublicKey) {
	return kp.privateKey, kp.publicKey
}

// SetKeyPair sets the private and public keys.
//
// # Parameters
//
//   - privateKey: ECDSA private key
//   - publicKey: ECDSA public key
//
// # Examples
//
//	keyPair.SetKeyPair(privateKey, publicKey)
func (kp *KeyPair) SetKeyPair(privateKey PrivateKey, publicKey PublicKey) {
	kp.privateKey = privateKey
	kp.publicKey = publicKey
}
