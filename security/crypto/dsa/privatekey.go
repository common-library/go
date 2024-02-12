// Package dsa provides dsa crypto related implementations.
package dsa

import (
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/pem"
)

// PrivateKey is struct that provides private key related methods.
type PrivateKey struct {
	privateKey *dsa.PrivateKey
}

// Sign is create a signature for message.
//
// ex) signature, err := privateKey.Sign(message)
func (this *PrivateKey) Sign(message string) (Signature, error) {
	hash := sha256.Sum256([]byte(message))
	if r, s, err := dsa.Sign(rand.Reader, this.privateKey, hash[:]); err != nil {
		return Signature{}, err
	} else {
		return Signature{R: r, S: s}, nil
	}
}

// Verify is verifies the signature.
//
// ex) result := privateKey.Verify(message, signature)
func (this *PrivateKey) Verify(message string, signature Signature) bool {
	hash := sha256.Sum256([]byte(message))

	return dsa.Verify(&this.privateKey.PublicKey, hash[:], signature.R, signature.S)
}

// Get is to get a *dsa.PrivateKey.
//
// ex) key := privateKey.Get()
func (this *PrivateKey) Get() *dsa.PrivateKey {
	return this.privateKey
}

// Set is to set a *dsa.PrivateKey.
//
// ex) privateKey.Set(key)
func (this *PrivateKey) Set(privateKey *dsa.PrivateKey) {
	this.privateKey = privateKey
}

// SetSizes is to set the primary key using sizes.
//
// ex) err := privateKey.SetSizes(crypto_dsa.L1024N160)
func (this *PrivateKey) SetSizes(parameterSizes dsa.ParameterSizes) error {
	params := new(dsa.Parameters)
	if err := dsa.GenerateParameters(params, rand.Reader, parameterSizes); err != nil {
		return err
	}

	privateKey := &dsa.PrivateKey{}
	privateKey.PublicKey.Parameters = *params

	if err := dsa.GenerateKey(privateKey, rand.Reader); err != nil {
		return err
	} else {
		this.Set(privateKey)
		return nil
	}
}

// GetPemAsn1 is to get a string in Pem/Asn1 format.
//
// ex) pemAsn1, err := privateKey.GetPemAsn1()
func (this *PrivateKey) GetPemAsn1() (string, error) {
	if blockBytes, err := asn1.Marshal(*this.privateKey); err != nil {
		return "", err
	} else {
		return string(pem.EncodeToMemory(
			&pem.Block{
				Type:    "DSA PRIVATE KEY",
				Headers: nil,
				Bytes:   blockBytes,
			})), nil
	}
}

// SetPemAsn1 is to set the primary key using a string in Pem/Asn1 format.
//
// ex) err := privateKey.SetPemAsn1(pemAsn1)
func (this *PrivateKey) SetPemAsn1(pemAsn1 string) error {
	this.privateKey = &dsa.PrivateKey{}

	block, _ := pem.Decode([]byte(pemAsn1))
	if _, err := asn1.Unmarshal(block.Bytes, this.privateKey); err != nil {
		this.privateKey = nil
		return err
	} else {
		return nil
	}

}

// GetPublicKey is to get a PublicKey.
//
// ex) key := privateKey.GetPublicKey()
func (this *PrivateKey) GetPublicKey() PublicKey {
	publicKey := PublicKey{}

	publicKey.Set(this.privateKey.PublicKey)

	return publicKey
}
