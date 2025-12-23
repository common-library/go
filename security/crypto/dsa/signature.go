// Package dsa provides dsa crypto related implementations.
//
// Deprecated: DSA is a legacy algorithm and should not be used for new applications.
// Use Ed25519 (crypto/ed25519) or other modern alternatives instead.
// DSA keys with 1024-bit moduli are cryptographically weak.
package dsa

import "math/big"

// Signature represents a DSA digital signature with R and S components.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// A DSA signature consists of two large integers (R and S) that together
// prove the authenticity of a message. Both components are required for
// signature verification.
//
// Security Warning:
//   - DSA signatures are vulnerable to weak random number generation
//   - Reusing the same random value (nonce) reveals the private key
//   - Use Ed25519 or ECDSA for new applications
//
// Example:
//
//	// For legacy compatibility only
//	signature := dsa.Signature{
//	    R: r,  // First component
//	    S: s,  // Second component
//	}
type Signature struct {
	R *big.Int
	S *big.Int
}
