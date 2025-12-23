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

// Sign creates a digital signature for the given message.
//
// Deprecated: DSA is cryptographically weak and should not be used.
// Use Ed25519 or ECDSA for new applications.
//
// This method generates a DSA signature using the private key. The message
// is automatically hashed with SHA-256 before signing.
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
//   - Vulnerable to weak RNG attacks (nonce reuse reveals private key)
//
// Example:
//
//	// For legacy compatibility only
//	privateKey := &dsa.PrivateKey{}
//	privateKey.SetSizes(dsa.L2048N256)
//	signature, err := privateKey.Sign("legacy message")
func (pk *PrivateKey) Sign(message string) (Signature, error) {
	hash := sha256.Sum256([]byte(message))
	if r, s, err := dsa.Sign(rand.Reader, pk.privateKey, hash[:]); err != nil {
		return Signature{}, err
	} else {
		return Signature{R: r, S: s}, nil
	}
}

// Verify verifies a digital signature against the message.
//
// Deprecated: DSA is cryptographically weak and should not be used.
// Use Ed25519 or ECDSA for new applications.
//
// This method verifies that a signature was created by the corresponding private key.
// The message is hashed with SHA-256 before verification.
//
// Parameters:
//   - message: The original message that was signed
//   - signature: The signature to verify
//
// Returns:
//   - bool: true if signature is valid, false otherwise
//
// Example:
//
//	// For legacy compatibility only
//	valid := privateKey.Verify("legacy message", signature)
func (pk *PrivateKey) Verify(message string, signature Signature) bool {
	hash := sha256.Sum256([]byte(message))

	return dsa.Verify(&pk.privateKey.PublicKey, hash[:], signature.R, signature.S)
}

// Get returns the underlying *dsa.PrivateKey.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// This method provides direct access to the Go standard library dsa.PrivateKey.
//
// Returns:
//   - *dsa.PrivateKey: The underlying private key
func (pk *PrivateKey) Get() *dsa.PrivateKey {
	return pk.privateKey
}

// Set assigns an existing *dsa.PrivateKey to this PrivateKey instance.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// Parameters:
//   - privateKey: The dsa.PrivateKey to assign
func (pk *PrivateKey) Set(privateKey *dsa.PrivateKey) {
	pk.privateKey = privateKey
}

// SetSizes generates a new DSA private key with the specified parameter sizes.
//
// Deprecated: DSA is cryptographically weak and should not be used.
// Use Ed25519 or ECDSA for new applications.
//
// This method generates DSA parameters and a key pair. L1024N160 is weak and
// should never be used.
//
// Parameters:
//   - parameterSizes: The L and N parameter sizes (L1024N160, L2048N224, L2048N256, L3072N256)
//
// Returns:
//   - error: Error if key generation fails
//
// Security Warning:
//   - L1024N160: Cryptographically broken, do not use
//   - L2048N256: Minimum for legacy compatibility
//   - L3072N256: Maximum DSA size, still deprecated
//
// Example:
//
//	// For legacy compatibility only
//	privateKey := &dsa.PrivateKey{}
//	err := privateKey.SetSizes(dsa.L2048N256)
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

// GetPemAsn1 encodes the private key as a PEM-encoded ASN.1 string.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// This method converts the DSA private key to PEM format using ASN.1 encoding.
//
// Returns:
//   - string: PEM-encoded DSA private key
//   - error: Error if encoding fails
//
// Behavior:
//   - Output format: "-----BEGIN DSA PRIVATE KEY-----\n...\n-----END DSA PRIVATE KEY-----"
//
// Example:
//
//	// For legacy compatibility only
//	pemString, err := privateKey.GetPemAsn1()
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

// SetPemAsn1 loads a private key from a PEM-encoded ASN.1 string.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// This method decodes a PEM-encoded DSA private key and sets it as the current key.
//
// Parameters:
//   - pemAsn1: PEM-encoded ASN.1 DSA private key string
//
// Returns:
//   - error: Error if decoding fails
//
// Example:
//
//	// For legacy compatibility only
//	pemData, _ := os.ReadFile("dsa_private.pem")
//	privateKey := &dsa.PrivateKey{}
//	err := privateKey.SetPemAsn1(string(pemData))
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

// GetPublicKey extracts the public key from the private key.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// Every DSA private key contains the corresponding public key. This method
// creates a PublicKey instance containing the public key portion.
//
// Returns:
//   - PublicKey: The corresponding public key
//
// Example:
//
//	// For legacy compatibility only
//	privateKey := &dsa.PrivateKey{}
//	privateKey.SetSizes(dsa.L2048N256)
//	publicKey := privateKey.GetPublicKey()
func (pk *PrivateKey) GetPublicKey() PublicKey {
	publicKey := PublicKey{}

	publicKey.Set(pk.privateKey.PublicKey)

	return publicKey
}
