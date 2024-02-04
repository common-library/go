// Package ecdsa provides ecdsa crypto related implementations.
package ecdsa

import "math/big"

// Signature is ecdsa signature
type Signature struct {
	R *big.Int
	S *big.Int
}
