// Package dsa provides dsa crypto related implementations.
package dsa

import "math/big"

// Signature is dsa signature
type Signature struct {
	R *big.Int
	S *big.Int
}
