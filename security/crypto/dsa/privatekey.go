// Package dsa provides dsa crypto related implementations.
//
// Deprecated: DSA is a legacy algorithm and should not be used for new applications.
// Use Ed25519 (crypto/ed25519) or other modern alternatives instead.
// DSA keys with 1024-bit moduli are cryptographically weak.
package dsa

import (
	//lint:ignore SA1019 DSA is deprecated but kept for compatibility
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
func (pk *PrivateKey) Sign(message string) (Signature, error) {
	hash := sha256.Sum256([]byte(message))
	if r, s, err := dsa.Sign(rand.Reader, pk.privateKey, hash[:]); err != nil {
		return Signature{}, err
	} else {
		return Signature{R: r, S: s}, nil
	}
}

// Verify is verifies the signature.
//
// ex) result := privateKey.Verify(message, signature)
func (pk *PrivateKey) Verify(message string, signature Signature) bool {
	hash := sha256.Sum256([]byte(message))

	return dsa.Verify(&pk.privateKey.PublicKey, hash[:], signature.R, signature.S)
}

// Get is to get a *dsa.PrivateKey.
//
// ex) key := privateKey.Get()
func (pk *PrivateKey) Get() *dsa.PrivateKey {
	return pk.privateKey
}

// Set is to set a *dsa.PrivateKey.
//
// ex) privateKey.Set(key)
func (pk *PrivateKey) Set(privateKey *dsa.PrivateKey) {
	pk.privateKey = privateKey
}

// SetSizes is to set the primary key using sizes.
//
// ex) err := privateKey.SetSizes(crypto_dsa.L1024N160)
func (pk *PrivateKey) SetSizes(parameterSizes dsa.ParameterSizes) error {
	params := new(dsa.Parameters)
	if err := dsa.GenerateParameters(params, rand.Reader, parameterSizes); err != nil {
		return err
	}

	privateKey := &dsa.PrivateKey{}
	privateKey.PublicKey.Parameters = *params

	if err := dsa.GenerateKey(privateKey, rand.Reader); err != nil {
		return err
	} else {
		pk.Set(privateKey)
		return nil
	}
}

// GetPemAsn1 is to get a string in Pem/Asn1 format.
//
// ex) pemAsn1, err := privateKey.GetPemAsn1()
func (pk *PrivateKey) GetPemAsn1() (string, error) {
	if blockBytes, err := asn1.Marshal(*pk.privateKey); err != nil {
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
func (pk *PrivateKey) SetPemAsn1(pemAsn1 string) error {
	pk.privateKey = &dsa.PrivateKey{}

	block, _ := pem.Decode([]byte(pemAsn1))
	if _, err := asn1.Unmarshal(block.Bytes, pk.privateKey); err != nil {
		pk.privateKey = nil
		return err
	} else {
		return nil
	}

}

// GetPublicKey is to get a PublicKey.
//
// ex) key := privateKey.GetPublicKey()
func (pk *PrivateKey) GetPublicKey() PublicKey {
	publicKey := PublicKey{}

	publicKey.Set(pk.privateKey.PublicKey)

	return publicKey
}
