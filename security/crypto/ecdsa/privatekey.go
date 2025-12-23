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

// Sign creates a digital signature for the given message.
//
// This method generates an ECDSA signature using the private key. The message
// is automatically hashed with SHA-256 before signing. Signatures are non-deterministic
// (each call produces a different signature for the same message).
//
// Parameters:
//   - message: The string to sign
//
// Returns:
//   - Signature: The signature containing R and S components
//   - error: Error if signing fails
//
// Behavior:
//   - Message is hashed with SHA-256 automatically
//   - Uses cryptographically secure random number generator
//   - Signature format: {R, S} where both are big.Int
//
// Example:
//
//	privateKey := &ecdsa.PrivateKey{}
//	privateKey.SetCurve(elliptic.P256())
//	signature, err := privateKey.Sign("message to sign")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PrivateKey) Sign(message string) (Signature, error) {
	hash := sha256.Sum256([]byte(message))
	if r, s, err := ecdsa.Sign(rand.Reader, pk.privateKey, hash[:]); err != nil {
		return Signature{}, err
	} else {
		return Signature{R: r, S: s}, nil
	}
}

// Verify verifies a digital signature against the message.
//
// This method verifies that a signature was created by the private key corresponding
// to this public key. The message is hashed with SHA-256 before verification.
//
// Parameters:
//   - message: The original message that was signed
//   - signature: The signature to verify
//
// Returns:
//   - bool: true if signature is valid, false otherwise
//
// Behavior:
//   - Message is hashed with SHA-256 automatically
//   - Timing-safe comparison
//   - Returns false on any error (no error return)
//
// Example:
//
//	valid := privateKey.Verify("message to sign", signature)
//	if valid {
//	    fmt.Println("Signature is valid")
//	}
func (pk *PrivateKey) Verify(message string, signature Signature) bool {
	hash := sha256.Sum256([]byte(message))

	return ecdsa.Verify(&pk.privateKey.PublicKey, hash[:], signature.R, signature.S)
}

// Get returns the underlying *ecdsa.PrivateKey.
//
// This method provides direct access to the Go standard library ecdsa.PrivateKey
// for use with other crypto functions that require the native type.
//
// Returns:
//   - *ecdsa.PrivateKey: The underlying private key
//
// Example:
//
//	nativeKey := privateKey.Get()
//	// Use with Go crypto functions
//	x509.MarshalECPrivateKey(nativeKey)
func (pk *PrivateKey) Get() *ecdsa.PrivateKey {
	return pk.privateKey
}

// Set assigns an existing *ecdsa.PrivateKey to this PrivateKey instance.
//
// This method allows setting the private key from an existing Go standard library
// ecdsa.PrivateKey, useful when loading keys from other sources.
//
// Parameters:
//   - privateKey: The ecdsa.PrivateKey to assign
//
// Example:
//
//	nativeKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
//	privateKey := &ecdsa.PrivateKey{}
//	privateKey.Set(nativeKey)
func (pk *PrivateKey) Set(privateKey *ecdsa.PrivateKey) {
	pk.privateKey = privateKey
}

// GetCurve returns the elliptic curve used by this key.
//
// This method retrieves the elliptic curve (P-256, P-384, or P-521) that
// was used to generate this key pair.
//
// Returns:
//   - elliptic.Curve: The elliptic curve of this key
//
// Behavior:
//   - Common curves: elliptic.P256(), P384(), P521()
//   - The curve determines key size and security level
//
// Example:
//
//	curve := privateKey.GetCurve()
//	fmt.Printf("Key uses curve: %s\n", curve.Params().Name)
func (pk *PrivateKey) GetCurve() elliptic.Curve {
	return pk.privateKey.PublicKey.Curve
}

// SetCurve generates a new ECDSA private key using the specified curve.
//
// This is a convenience method that generates a new key pair on the specified
// elliptic curve and sets it as the current private key.
//
// Parameters:
//   - curve: The elliptic curve to use (P-256, P-384, or P-521)
//
// Returns:
//   - error: Error if key generation fails
//
// Behavior:
//   - Generates both private and public keys
//   - Uses cryptographically secure random number generator
//   - P-256 recommended for most applications
//
// Example:
//
//	privateKey := &ecdsa.PrivateKey{}
//	err := privateKey.SetCurve(elliptic.P256())
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PrivateKey) SetCurve(curve elliptic.Curve) error {
	if privateKey, err := ecdsa.GenerateKey(curve, rand.Reader); err != nil {
		return err
	} else {
		pk.Set(privateKey)
		return nil
	}
}

// GetPemEC encodes the private key as a PEM-encoded EC private key string.
//
// This method converts the ECDSA private key to PEM format using the EC private
// key structure. The format uses "ECDSA PRIVATE KEY" as the PEM block type.
//
// Returns:
//   - string: PEM-encoded EC private key
//   - error: Error if encoding fails
//
// Behavior:
//   - Output format: "-----BEGIN ECDSA PRIVATE KEY-----\n...\n-----END ECDSA PRIVATE KEY-----"
//   - Includes the curve parameters
//   - ECDSA-specific format
//
// Example:
//
//	pemString, err := privateKey.GetPemEC()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("ecdsa_private.pem", []byte(pemString), 0600)
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

// SetPemEC loads a private key from a PEM-encoded EC private key string.
//
// This method decodes a PEM-encoded EC private key and sets it as the current
// private key. The input should have the "ECDSA PRIVATE KEY" PEM block type.
//
// Parameters:
//   - pemEC: PEM-encoded EC private key string
//
// Returns:
//   - error: Error if decoding fails or format is invalid
//
// Behavior:
//   - Expects "-----BEGIN ECDSA PRIVATE KEY-----" header
//   - Automatically detects the curve from key data
//   - Validates the EC key structure
//
// Example:
//
//	pemData, _ := os.ReadFile("ecdsa_private.pem")
//	privateKey := &ecdsa.PrivateKey{}
//	err := privateKey.SetPemEC(string(pemData))
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PrivateKey) SetPemEC(pemEC string) error {
	block, _ := pem.Decode([]byte(pemEC))

	if key, err := x509.ParseECPrivateKey(block.Bytes); err != nil {
		return err
	} else {
		pk.Set(key)
		return nil
	}
}

// GetPublicKey extracts the public key from the private key.
//
// Every ECDSA private key contains the corresponding public key. This method
// creates a PublicKey instance containing the public key portion.
//
// Returns:
//   - PublicKey: The corresponding public key
//
// Behavior:
//   - The public key is derived from the private key's curve point
//   - The returned public key can be safely shared with others
//   - Useful for distributing public keys for signature verification
//
// Example:
//
//	privateKey := &ecdsa.PrivateKey{}
//	privateKey.SetCurve(elliptic.P256())
//	publicKey := privateKey.GetPublicKey()
//	// Share publicKey with others for verification
func (pk *PrivateKey) GetPublicKey() PublicKey {
	publicKey := PublicKey{}

	publicKey.Set(pk.privateKey.PublicKey)

	return publicKey
}
