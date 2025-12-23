// Package ed25519 provides Ed25519 digital signature cryptography.
//
// This package implements Ed25519 public-key signature system based on the
// elliptic curve Ed25519. It provides key generation, signing, and verification
// operations with simplified interfaces for key pair management.
//
// # Features
//
//   - Ed25519 key pair generation
//   - Digital signature creation and verification
//   - PEM/PKCS8/PKIX format support
//   - SSH public key format conversion
//   - Type-safe key management
//
// # Basic Example
//
//	privateKey := &ed25519.PrivateKey{}
//	err := privateKey.SetDefault()
//	signature := privateKey.Sign("message")
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

// Sign creates a digital signature for the message.
//
// # Parameters
//
//   - message: Text message to sign
//
// # Returns
//
//   - []byte: 64-byte Ed25519 signature
//
// # Examples
//
//	signature := privateKey.Sign("Hello, World!")
func (pk *PrivateKey) Sign(message string) []byte {
	return ed25519.Sign(pk.privateKey, []byte(message))
}

// Verify verifies a digital signature using the public key.
//
// # Parameters
//
//   - message: Original text message
//   - signature: 64-byte signature to verify
//
// # Returns
//
//   - bool: true if signature is valid, false otherwise
//
// # Examples
//
//	valid := privateKey.Verify("message", signature)
func (pk *PrivateKey) Verify(message string, signature []byte) bool {
	return ed25519.Verify(pk.publicKey, []byte(message), signature)
}

// Get retrieves the underlying ed25519.PrivateKey.
//
// # Returns
//
//   - ed25519.PrivateKey: Raw private key (64 bytes)
//
// # Examples
//
//	key := privateKey.Get()
func (pk *PrivateKey) Get() ed25519.PrivateKey {
	return pk.privateKey
}

// Set sets the private key from an ed25519.PrivateKey.
//
// # Parameters
//
//   - privateKey: Ed25519 private key to set
//
// # Examples
//
//	privateKey.Set(key)
func (pk *PrivateKey) Set(privateKey ed25519.PrivateKey) {
	pk.privateKey = privateKey
	pk.publicKey = privateKey.Public().(ed25519.PublicKey)
}

// SetDefault generates a new random private key.
//
// # Returns
//
//   - error: Error if key generation fails, nil on success
//
// # Examples
//
//	err := privateKey.SetDefault()
func (pk *PrivateKey) SetDefault() error {
	if _, privateKey, err := ed25519.GenerateKey(rand.Reader); err != nil {
		return err
	} else {
		pk.Set(privateKey)
		return nil
	}
}

// GetPemPKCS8 returns the private key in PEM-encoded PKCS#8 format.
//
// # Returns
//
//   - string: PEM-encoded private key
//   - error: Error if encoding fails, nil on success
//
// # Examples
//
//	pemString, err := privateKey.GetPemPKCS8()
func (pk *PrivateKey) GetPemPKCS8() (string, error) {
	if blockBytes, err := x509.MarshalPKCS8PrivateKey(pk.privateKey); err != nil {
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

// SetPemPKCS8 sets the private key from a PEM-encoded PKCS#8 string.
//
// # Parameters
//
//   - pemPKCS8: PEM-encoded private key string
//
// # Returns
//
//   - error: Error if decoding or parsing fails, nil on success
//
// # Examples
//
//	err := privateKey.SetPemPKCS8(pemString)
func (pk *PrivateKey) SetPemPKCS8(pemPKCS8 string) error {
	block, _ := pem.Decode([]byte(pemPKCS8))

	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return err
	} else {
		pk.Set(key.(ed25519.PrivateKey))
		return nil
	}
}

// GetPublicKey derives the public key from the private key.
//
// # Returns
//
//   - PublicKey: Corresponding public key
//
// # Examples
//
//	publicKey := privateKey.GetPublicKey()
func (pk *PrivateKey) GetPublicKey() PublicKey {
	publicKey := PublicKey{}

	publicKey.Set(pk.publicKey)

	return publicKey
}
