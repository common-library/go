// Package ed25519 provides ed25519 crypto related implementations.
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

// Verify is verifies the signature.
//
// ex) result := publicKey.Verify(message, signature)
func (pk *PublicKey) Verify(message string, signature []byte) bool {
	return ed25519.Verify(pk.publicKey, []byte(message), signature)
}

// Get is to get a ed25519.PublicKey.
//
// ex) key := publicKey.Get()
func (pk *PublicKey) Get() ed25519.PublicKey {
	return pk.publicKey
}

// Set is to set a ed25519.PublicKey.
//
// ex) publicKey.Set(key)
func (pk *PublicKey) Set(publicKey ed25519.PublicKey) {
	pk.publicKey = publicKey
}

// GetPemPKIX is to get a string in Pem/PKIX format.
//
// ex) pemPKIX, err := publicKey.GetPemPKIX()
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

// SetPemPKIX is to set the public key using a string in Pem/PKIX format.
//
// ex) err := publicKey.SetPemPKIX(pemPKIX)
func (pk *PublicKey) SetPemPKIX(pemPKIX string) error {
	block, _ := pem.Decode([]byte(pemPKIX))

	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return err
	} else {
		pk.publicKey = key.(ed25519.PublicKey)
		return nil
	}
}

// GetSsh is to get a string in ssh format.
//
// ex) sshKey, err := publicKey.GetSsh()
func (pk *PublicKey) GetSsh() (string, error) {
	if publicKey, err := ssh.NewPublicKey(pk.publicKey); err != nil {
		return "", err
	} else {
		return string(ssh.MarshalAuthorizedKey(publicKey)), nil
	}
}

// SetSsh is to set the public key using a string in ssh format.
//
// ex) err := publicKey.SetSsh(sshKey)
func (pk *PublicKey) SetSsh(sshKey string) error {
	if key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshKey)); err != nil {
		return err
	} else {
		return pk.SetSshPublicKey(key)
	}
}

// GetSshPublicKey is to get a ssh.PublicKey.
//
// ex) key, err := publicKey.GetSshPublicKey()
func (pk *PublicKey) GetSshPublicKey() (ssh.PublicKey, error) {
	return ssh.NewPublicKey(pk.publicKey)
}

// SetSshPublicKey is to set the public key using ssh.PublicKey.
//
// ex) err := publicKey.SetSshPublicKey(key)
func (pk *PublicKey) SetSshPublicKey(publicKey ssh.PublicKey) error {
	if key, err := ssh.ParsePublicKey(publicKey.Marshal()); err != nil {
		return err
	} else {
		pk.publicKey = key.(ssh.CryptoPublicKey).CryptoPublicKey().(ed25519.PublicKey)
		return nil
	}
}
