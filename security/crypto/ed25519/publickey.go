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
//	publicKey := &ed25519.PublicKey{}
//	publicKey.SetPemPKIX(pemString)
//	valid := publicKey.Verify("message", signature)
package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// PublicKey is struct that provides public key related methods.
type PublicKey struct {
	publicKey ed25519.PublicKey
}

// Verify verifies a digital signature.
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
//	valid := publicKey.Verify("message", signature)
func (pk *PublicKey) Verify(message string, signature []byte) bool {
	return ed25519.Verify(pk.publicKey, []byte(message), signature)
}

// Get retrieves the underlying ed25519.PublicKey.
//
// # Returns
//
//   - ed25519.PublicKey: Raw public key (32 bytes)
//
// # Examples
//
//	key := publicKey.Get()
func (pk *PublicKey) Get() ed25519.PublicKey {
	return pk.publicKey
}

// Set sets the public key from an ed25519.PublicKey.
//
// # Parameters
//
//   - publicKey: Ed25519 public key to set
//
// # Examples
//
//	publicKey.Set(key)
func (pk *PublicKey) Set(publicKey ed25519.PublicKey) {
	pk.publicKey = publicKey
}

// GetPemPKIX returns the public key in PEM-encoded PKIX format.
//
// # Returns
//
//   - string: PEM-encoded public key
//   - error: Error if encoding fails, nil on success
//
// # Examples
//
//	pemString, err := publicKey.GetPemPKIX()
func (pk *PublicKey) GetPemPKIX() (string, error) {
	if blockBytes, err := x509.MarshalPKIXPublicKey(pk.publicKey); err != nil {
		return "", err
	} else {
		return string(pem.EncodeToMemory(
			&pem.Block{
				Type:    "PUBLIC KEY",
				Headers: nil,
				Bytes:   blockBytes,
			})), nil
	}
}

// SetPemPKIX sets the public key from a PEM-encoded PKIX string.
//
// # Parameters
//
//   - pemPKIX: PEM-encoded public key string
//
// # Returns
//
//   - error: Error if decoding or parsing fails, nil on success
//
// # Examples
//
//	err := publicKey.SetPemPKIX(pemString)
func (pk *PublicKey) SetPemPKIX(pemPKIX string) error {
	block, _ := pem.Decode([]byte(pemPKIX))

	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return err
	} else {
		pk.publicKey = key.(ed25519.PublicKey)
		return nil
	}
}

// GetSsh returns the public key in SSH authorized_keys format.
//
// # Returns
//
//   - string: SSH public key string
//   - error: Error if encoding fails, nil on success
//
// # Examples
//
//	sshKey, err := publicKey.GetSsh()
func (pk *PublicKey) GetSsh() (string, error) {
	if publicKey, err := ssh.NewPublicKey(pk.publicKey); err != nil {
		return "", err
	} else {
		return string(ssh.MarshalAuthorizedKey(publicKey)), nil
	}
}

// SetSsh sets the public key from an SSH authorized_keys format string.
//
// # Parameters
//
//   - sshKey: SSH public key string
//
// # Returns
//
//   - error: Error if parsing fails, nil on success
//
// # Examples
//
//	err := publicKey.SetSsh(sshKey)
func (pk *PublicKey) SetSsh(sshKey string) error {
	if key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshKey)); err != nil {
		return err
	} else {
		return pk.SetSshPublicKey(key)
	}
}

// GetSshPublicKey returns the key as an ssh.PublicKey.
//
// # Returns
//
//   - ssh.PublicKey: SSH public key interface
//   - error: Error if conversion fails, nil on success
//
// # Examples
//
//	sshKey, err := publicKey.GetSshPublicKey()
func (pk *PublicKey) GetSshPublicKey() (ssh.PublicKey, error) {
	return ssh.NewPublicKey(pk.publicKey)
}

// SetSshPublicKey sets the public key from an ssh.PublicKey.
//
// # Parameters
//
//   - publicKey: SSH public key to set
//
// # Returns
//
//   - error: Error if conversion fails, nil on success
//
// # Examples
//
//	err := publicKey.SetSshPublicKey(sshKey)
func (pk *PublicKey) SetSshPublicKey(publicKey ssh.PublicKey) error {
	if key, err := ssh.ParsePublicKey(publicKey.Marshal()); err != nil {
		return err
	} else {
		pk.publicKey = key.(ssh.CryptoPublicKey).CryptoPublicKey().(ed25519.PublicKey)
		return nil
	}
}
