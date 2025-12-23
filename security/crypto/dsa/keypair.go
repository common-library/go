// Package dsa provides DSA digital signature cryptography.
//
// # Deprecated
//
// DSA is a legacy algorithm and should not be used for new applications.
// Use Ed25519 (crypto/ed25519) or ECDSA (crypto/ecdsa) instead.
// DSA keys with 1024-bit moduli are cryptographically weak.
//
// This package is maintained only for compatibility with legacy systems.
//
// # Features
//
//   - DSA key pair generation (L1024N160, L2048N224, L2048N256, L3072N256)
//   - Digital signature creation and verification
//   - PEM format support
//
// # Migration Recommendation
//
//	// Instead of DSA:
//	import "github.com/common-library/go/security/crypto/ed25519"
//	keyPair := &ed25519.KeyPair{}
//	keyPair.Generate()
package dsa

//lint:ignore SA1019 DSA is deprecated but kept for compatibility
import "crypto/dsa"

// KeyPair provides DSA key pair operations.
//
// Deprecated: Use Ed25519 or ECDSA instead.
type KeyPair struct {
	privateKey PrivateKey
	publicKey  PublicKey
}

// Generate creates a new DSA key pair.
//
// Deprecated: Use Ed25519 or ECDSA for new applications.
//
// # Parameters
//
//   - parameterSizes: Key size (L1024N160, L2048N224, L2048N256, L3072N256)
//
// # Returns
//
//   - error: Error if generation fails, nil on success
//
// # Examples
//
//	keyPair := &dsa.KeyPair{}
//	err := keyPair.Generate(dsa.L2048N256)
func (kp *KeyPair) Generate(parameterSizes dsa.ParameterSizes) error {
	if err := kp.privateKey.SetSizes(parameterSizes); err != nil {
		return err
	} else {
		kp.publicKey = kp.privateKey.GetPublicKey()
		return nil
	}
}

// Sign creates a digital signature for the message.
//
// Deprecated: Use Ed25519 or ECDSA for new applications.
//
// # Parameters
//
//   - message: Text message to sign
//
// # Returns
//
//   - Signature: DSA signature (R and S components)
//   - error: Error if signing fails, nil on success
//
// # Examples
//
//	signature, err := keyPair.Sign("message")
func (kp *KeyPair) Sign(message string) (Signature, error) {
	return kp.privateKey.Sign(message)
}

// Verify verifies a digital signature.
//
// Deprecated: Use Ed25519 or ECDSA for new applications.
//
// # Parameters
//
//   - message: Original text message
//   - signature: DSA signature to verify
//
// # Returns
//
//   - bool: true if signature is valid, false otherwise
//
// # Examples
//
//	valid := keyPair.Verify("message", signature)
func (kp *KeyPair) Verify(message string, signature Signature) bool {
	return kp.publicKey.Verify(message, signature)
}

// GetKeyPair retrieves the private and public keys.
//
// Deprecated: Use Ed25519 or ECDSA for new applications.
//
// # Returns
//
//   - privateKey: DSA private key
//   - publicKey: DSA public key
//
// # Examples
//
//	privateKey, publicKey := keyPair.GetKeyPair()
func (kp *KeyPair) GetKeyPair() (privateKey PrivateKey, publicKey PublicKey) {
	return kp.privateKey, kp.publicKey
}

// SetKeyPair sets the private and public keys.
//
// Deprecated: Use Ed25519 or ECDSA for new applications.
//
// # Parameters
//
//   - privateKey: DSA private key
//   - publicKey: DSA public key
//
// # Examples
//
//	keyPair.SetKeyPair(privateKey, publicKey)
func (kp *KeyPair) SetKeyPair(privateKey PrivateKey, publicKey PublicKey) {
	kp.privateKey = privateKey
	kp.publicKey = publicKey
}
