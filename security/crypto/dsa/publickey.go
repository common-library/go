// Package dsa provides dsa crypto related implementations.
package dsa

import (
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
func (this *PublicKey) Verify(message string, signature Signature) bool {
	hash := sha256.Sum256([]byte(message))

	return dsa.Verify(&this.publicKey, hash[:], signature.R, signature.S)
}

// Get is to get a dsa.PublicKey.
//
// ex) key := publicKey.Get()
func (this *PublicKey) Get() dsa.PublicKey {
	return this.publicKey
}

// Set is to set a dsa.PublicKey.
//
// ex) publicKey.Set(key)
func (this *PublicKey) Set(publicKey dsa.PublicKey) {
	this.publicKey = publicKey
}

// GetPemAsn1 is to get a string in Pem/Asn1 format.
//
// ex) pemAsn1, err := publicKey.GetPemAsn1()
func (this *PublicKey) GetPemAsn1() (string, error) {
	if blockBytes, err := asn1.Marshal(this.publicKey); err != nil {
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
func (this *PublicKey) SetPemAsn1(pemAsn1 string) error {
	block, _ := pem.Decode([]byte(pemAsn1))
	if _, err := asn1.Unmarshal(block.Bytes, &this.publicKey); err != nil {
		return err
	} else {
		return nil
	}

}

// GetSsh is to get a string in ssh format.
//
// ex) sshKey, err := publicKey.GetSsh()
func (this *PublicKey) GetSsh() (string, error) {
	if publicKey, err := ssh.NewPublicKey(&this.publicKey); err != nil {
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
	return ssh.NewPublicKey(&this.publicKey)
}

// SetSshPublicKey is to set the public key using ssh.PublicKey.
//
// ex) err := publicKey.SetSshPublicKey(key)
func (this *PublicKey) SetSshPublicKey(publicKey ssh.PublicKey) error {
	if key, err := ssh.ParsePublicKey(publicKey.Marshal()); err != nil {
		return err
	} else {
		this.publicKey = *key.(ssh.CryptoPublicKey).CryptoPublicKey().(*dsa.PublicKey)
		return nil
	}
}
