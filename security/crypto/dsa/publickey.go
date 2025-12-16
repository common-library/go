// Package dsa provides dsa crypto related implementations.
//
// Deprecated: DSA is a legacy algorithm and should not be used for new applications.
// Use Ed25519 (crypto/ed25519) or other modern alternatives instead.
// DSA keys with 1024-bit moduli are cryptographically weak.
package dsa

import (
	//lint:ignore SA1019 DSA is deprecated but kept for compatibility
	"crypto/dsa"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// PublicKey is struct that provides public key related methods.
type PublicKey struct {
	publicKey dsa.PublicKey
}

// Verify is verifies the signature.
//
// ex) result := publicKey.Verify(message, signature)
func (pub *PublicKey) Verify(message string, signature Signature) bool {
	hash := sha256.Sum256([]byte(message))

	return dsa.Verify(&pub.publicKey, hash[:], signature.R, signature.S)
}

// Get is to get a dsa.PublicKey.
//
// ex) key := publicKey.Get()
func (pub *PublicKey) Get() dsa.PublicKey {
	return pub.publicKey
}

// Set is to set a dsa.PublicKey.
//
// ex) publicKey.Set(key)
func (pub *PublicKey) Set(publicKey dsa.PublicKey) {
	pub.publicKey = publicKey
}

// GetPemAsn1 is to get a string in Pem/Asn1 format.
//
// ex) pemAsn1, err := publicKey.GetPemAsn1()
func (pub *PublicKey) GetPemAsn1() (string, error) {
	if blockBytes, err := asn1.Marshal(pub.publicKey); err != nil {
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

// SetPemAsn1 is to set the public key using a string in Pem/Asn1 format.
//
// ex) err := publicKey.SetPemAsn1(pemAsn1)
func (pub *PublicKey) SetPemAsn1(pemAsn1 string) error {
	block, _ := pem.Decode([]byte(pemAsn1))
	if _, err := asn1.Unmarshal(block.Bytes, &pub.publicKey); err != nil {
		return err
	} else {
		return nil
	}

}

// GetSsh is to get a string in ssh format.
//
// ex) sshKey, err := publicKey.GetSsh()
func (pub *PublicKey) GetSsh() (string, error) {
	if publicKey, err := ssh.NewPublicKey(&pub.publicKey); err != nil {
		return "", err
	} else {
		return string(ssh.MarshalAuthorizedKey(publicKey)), nil
	}
}

// SetSsh is to set the public key using a string in ssh format.
//
// ex) err := publicKey.SetSsh(sshKey)
func (pub *PublicKey) SetSsh(sshKey string) error {
	if key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshKey)); err != nil {
		return err
	} else {
		return pub.SetSshPublicKey(key)
	}
}

// GetSshPublicKey is to get a ssh.PublicKey.
//
// ex) key, err := publicKey.GetSshPublicKey()
func (pub *PublicKey) GetSshPublicKey() (ssh.PublicKey, error) {
	return ssh.NewPublicKey(&pub.publicKey)
}

// SetSshPublicKey is to set the public key using ssh.PublicKey.
//
// ex) err := publicKey.SetSshPublicKey(key)
func (pub *PublicKey) SetSshPublicKey(publicKey ssh.PublicKey) error {
	if key, err := ssh.ParsePublicKey(publicKey.Marshal()); err != nil {
		return err
	} else {
		pub.publicKey = *key.(ssh.CryptoPublicKey).CryptoPublicKey().(*dsa.PublicKey)
		return nil
	}
}
