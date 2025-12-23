// Package ecdsa provides ecdsa crypto related implementations.
package ecdsa

import "math/big"

// Signature represents an ECDSA digital signature with R and S components.
//
// An ECDSA signature consists of two large integers (R and S) that together
// prove the authenticity and integrity of a message. Both components are
// required for signature verification.
//
// Features:
//   - Non-deterministic: Same message produces different signatures
//   - Compact: Smaller than RSA signatures for equivalent security
//   - Secure: Based on elliptic curve discrete logarithm problem
//
// Security:
//   - Each signature uses a unique random nonce automatically
//   - Nonce reuse would reveal the private key (prevented by crypto/rand)
//   - R and S are typically 32-48 bytes each depending on curve
//
// Example:
//
//	signature := ecdsa.Signature{
//	    R: r,  // First signature component
//	    S: s,  // Second signature component
//	}
//
//	// Verify signature
//	valid := publicKey.Verify(message, signature)
type Signature struct {
	R *big.Int
	S *big.Int
}
