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
func (this *PublicKey) Verify(message string, signature []byte) bool {
	return ed25519.Verify(this.publicKey, []byte(message), signature)
}

// Get is to get a ed25519.PublicKey.
//
// ex) key := publicKey.Get()
func (this *PublicKey) Get() ed25519.PublicKey {
	return this.publicKey
}

// Set is to set a ed25519.PublicKey.
//
// ex) publicKey.Set(key)
func (this *PublicKey) Set(publicKey ed25519.PublicKey) {
	this.publicKey = publicKey
}

// GetPemPKIX is to get a string in Pem/PKIX format.
//
// ex) pemPKIX, err := publicKey.GetPemPKIX()
func (this *PublicKey) GetPemPKIX() (string, error) {
	if blockBytes, err := x509.MarshalPKIXPublicKey(this.publicKey); err != nil {
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
func (this *PublicKey) SetPemPKIX(pemPKIX string) error {
	block, _ := pem.Decode([]byte(pemPKIX))

	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return err
	} else {
		this.publicKey = key.(ed25519.PublicKey)
		return nil
	}
}

// GetSsh is to get a string in ssh format.
//
// ex) sshKey, err := publicKey.GetSsh()
func (this *PublicKey) GetSsh() (string, error) {
	if publicKey, err := ssh.NewPublicKey(this.publicKey); err != nil {
		return "", err
	} else {
		return string(ssh.MarshalAuthorizedKey(publicKey)), nil
	}
}

// SetSsh is to set the public key using a string in ssh format.
//
// ex) err := publicKey.SetSsh(sshKey)
func (this *PublicKey) SetSsh(sshKey string) error {
	if key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshKey)); err != nil {
		return err
	} else {
		return this.SetSshPublicKey(key)
	}
}

// GetSshPublicKey is to get a ssh.PublicKey.
//
// ex) key, err := publicKey.GetSshPublicKey()
func (this *PublicKey) GetSshPublicKey() (ssh.PublicKey, error) {
	return ssh.NewPublicKey(this.publicKey)
}

// SetSshPublicKey is to set the public key using ssh.PublicKey.
//
// ex) err := publicKey.SetSshPublicKey(key)
func (this *PublicKey) SetSshPublicKey(publicKey ssh.PublicKey) error {
	if key, err := ssh.ParsePublicKey(publicKey.Marshal()); err != nil {
		return err
	} else {
		this.publicKey = key.(ssh.CryptoPublicKey).CryptoPublicKey().(ed25519.PublicKey)
		return nil
	}
}
