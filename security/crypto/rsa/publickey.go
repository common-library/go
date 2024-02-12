// Package rsa provides rsa crypto related implementations.
package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// PublicKey is struct that provides public key related methods.
type PublicKey struct {
	publicKey rsa.PublicKey
}

// EncryptPKCS1v15 is encrypt plaintext.
//
// ex) ciphertext, err := publicKey.EncryptPKCS1v15(plaintext)
func (this *PublicKey) EncryptPKCS1v15(plaintext string) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, &this.publicKey, []byte(plaintext))
}

// EncryptOAEP is encrypt plaintext.
//
// ex) ciphertext, err := publicKey.EncryptOAEP(plaintext)
func (this *PublicKey) EncryptOAEP(plaintext string) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, &this.publicKey, []byte(plaintext), nil)
}

// Get is to get a rsa.PublicKey.
//
// ex) key := publicKey.Get()
func (this *PublicKey) Get() rsa.PublicKey {
	return this.publicKey
}

// Set is to set a rsa.PublicKey.
//
// ex) publicKey.Set(key)
func (this *PublicKey) Set(publicKey rsa.PublicKey) {
	this.publicKey = publicKey
}

// GetPemPKCS1 is to get a string in Pem/PKCS1 format.
//
// ex) pemPKCS1 := publicKey.GetPemPKCS1()
func (this *PublicKey) GetPemPKCS1() string {
	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:    "RSA PUBLIC KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PublicKey(&this.publicKey),
		}))
}

// SetPemPKCS1 is to set the public key using a string in Pem/PKCS1 format.
//
// ex) err := publicKey.SetPemPKCS1(pemPKCS1)
func (this *PublicKey) SetPemPKCS1(pemPKCS1 string) error {
	block, _ := pem.Decode([]byte(pemPKCS1))

	if key, err := x509.ParsePKCS1PublicKey(block.Bytes); err != nil {
		return err
	} else {
		this.Set(*key)
		return nil
	}
}

// GetPemPKIX is to get a string in Pem/PKIX format.
//
// ex) pemPKIX, err := publicKey.GetPemPKIX()
func (this *PublicKey) GetPemPKIX() (string, error) {
	if bytes, err := x509.MarshalPKIXPublicKey(&this.publicKey); err != nil {
		return "", err
	} else {
		return string(pem.EncodeToMemory(
			&pem.Block{
				Type:    "PUBLIC KEY",
				Headers: nil,
				Bytes:   bytes,
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
		this.Set(*key.(*rsa.PublicKey))
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
		this.publicKey = *key.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
		return nil
	}
}
