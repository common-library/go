// Package ecdsa provides ecdsa crypto related implementations.
package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

// PrivateKey is struct that provides private key related methods.
type PrivateKey struct {
	privateKey *ecdsa.PrivateKey
}

// Sign is create a signature for message.
//
// ex) signature, err := privateKey.Sign(message)
func (pk *PrivateKey) Sign(message string) (Signature, error) {
	hash := sha256.Sum256([]byte(message))
	if r, s, err := ecdsa.Sign(rand.Reader, pk.privateKey, hash[:]); err != nil {
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

	return ecdsa.Verify(&pk.privateKey.PublicKey, hash[:], signature.R, signature.S)
}

// Get is to get a *ecdsa.PrivateKey.
//
// ex) key := privateKey.Get()
func (pk *PrivateKey) Get() *ecdsa.PrivateKey {
	return pk.privateKey
}

// Set is to set a *ecdsa.PrivateKey.
//
// ex) privateKey.Set(key)
func (pk *PrivateKey) Set(privateKey *ecdsa.PrivateKey) {
	pk.privateKey = privateKey
}

// GetCurve is to get a elliptic.Curve.
//
// ex) curve := privateKey.GetCurve()
func (pk *PrivateKey) GetCurve() elliptic.Curve {
	return pk.privateKey.PublicKey.Curve
}

// SetCurve is to set the primary key using Curve.
//
// ex) err := privateKey.SetCurve(elliptic.P384())
func (pk *PrivateKey) SetCurve(curve elliptic.Curve) error {
	if privateKey, err := ecdsa.GenerateKey(curve, rand.Reader); err != nil {
		return err
	} else {
		pk.Set(privateKey)
		return nil
	}
}

// GetPemEC is to get a string in Pem/EC format.
//
// ex) pemEC, err := privateKey.GetPemEC()
func (pk *PrivateKey) GetPemEC() (string, error) {
	if blockBytes, err := x509.MarshalECPrivateKey(pk.privateKey); err != nil {
		return "", err
	} else {
		return string(pem.EncodeToMemory(
			&pem.Block{
				Type:    "ECDSA PRIVATE KEY",
				Headers: nil,
				Bytes:   blockBytes,
			})), nil
	}
}

// SetPemEC is to set the primary key using a string in Pem/EC format.
//
// ex) err := privateKey.SetPemEC(pemEC)
func (pk *PrivateKey) SetPemEC(pemEC string) error {
	block, _ := pem.Decode([]byte(pemEC))

	if key, err := x509.ParseECPrivateKey(block.Bytes); err != nil {
		return err
	} else {
		pk.Set(key)
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
