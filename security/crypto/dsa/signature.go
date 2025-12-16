// Package dsa provides dsa crypto related implementations.
//
// Deprecated: DSA is a legacy algorithm and should not be used for new applications.
// Use Ed25519 (crypto/ed25519) or other modern alternatives instead.
// DSA keys with 1024-bit moduli are cryptographically weak.
package dsa

import "math/big"

// Signature is dsa signature
type Signature struct {
	R *big.Int
	S *big.Int
}
