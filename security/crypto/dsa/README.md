# DSA (Deprecated)

Digital Signature Algorithm - Legacy support only.

## ⚠️ Deprecation Notice

**DSA is deprecated and should not be used for new applications.**

Use modern alternatives:
- **Ed25519** - Recommended for new applications
- **ECDSA** - NIST-approved alternative

This package is maintained for legacy compatibility only.

## Overview

The dsa package provides Digital Signature Algorithm signatures as specified in FIPS 186-4. DSA is considered legacy technology and has been superseded by more secure and efficient algorithms.

## Migration Guide

### From DSA to Ed25519 (Recommended)

```go
// Old: DSA
import "github.com/common-library/go/security/crypto/dsa"

dsaKey := &dsa.KeyPair{}
dsaKey.Generate(dsa.L2048N256)
signature, _ := dsaKey.Sign(message)

// New: Ed25519
import "github.com/common-library/go/security/crypto/ed25519"

ed25519Key := &ed25519.KeyPair{}
ed25519Key.Generate()
signature, _ := ed25519Key.Sign(message)
```

### From DSA to ECDSA (NIST Alternative)

```go
// Old: DSA
import "github.com/common-library/go/security/crypto/dsa"

dsaKey := &dsa.KeyPair{}
dsaKey.Generate(dsa.L2048N256)

// New: ECDSA
import (
    "crypto/elliptic"
    "github.com/common-library/go/security/crypto/ecdsa"
)

ecdsaKey := &ecdsa.KeyPair{}
ecdsaKey.Generate(elliptic.P256())
```

## Why DSA is Deprecated

1. **Security Concerns** - Vulnerable to weak random number generation
2. **Performance** - Slower than modern alternatives
3. **Key Size** - Large keys for equivalent security
4. **Industry Movement** - NIST recommends ECDSA/Ed25519
5. **Limited Support** - Decreasing support in modern systems

## Legacy Documentation

For maintaining existing DSA systems only.

### Installation

```bash
go get -u github.com/common-library/go/security/crypto/dsa
```

### API Reference

#### Generate

```go
func (kp *KeyPair) Generate(sizes dsa.ParameterSizes) error
```

**Deprecated**: Use Ed25519 or ECDSA instead.

Generates DSA key pair with specified parameter sizes:
- `dsa.L1024N160` - 1024-bit L, 160-bit N (weak, not recommended)
- `dsa.L2048N224` - 2048-bit L, 224-bit N
- `dsa.L2048N256` - 2048-bit L, 256-bit N
- `dsa.L3072N256` - 3072-bit L, 256-bit N (maximum)

#### Sign

```go
func (kp *KeyPair) Sign(message string) (Signature, error)
```

**Deprecated**: Use Ed25519 or ECDSA instead.

Creates DSA signature.

#### Verify

```go
func (kp *KeyPair) Verify(message string, signature Signature) bool
```

**Deprecated**: Use Ed25519 or ECDSA instead.

Verifies DSA signature.

## Legacy Example

```go
package main

import (
    "crypto/dsa"
    "fmt"
    "log"
    "github.com/common-library/go/security/crypto/dsa"
)

func main() {
    // ⚠️ Legacy code - do not use for new applications
    keyPair := &dsa.KeyPair{}
    err := keyPair.Generate(dsa.L2048N256)
    if err != nil {
        log.Fatal(err)
    }
    
    message := "Legacy message"
    signature, err := keyPair.Sign(message)
    if err != nil {
        log.Fatal(err)
    }
    
    valid := keyPair.Verify(message, signature)
    fmt.Printf("Signature valid: %v\n", valid)
}
```

## Parameter Size Comparison

| Size | L (bits) | N (bits) | Security | Status |
|------|----------|----------|----------|--------|
| L1024N160 | 1024 | 160 | ~80-bit | ⚠️ Weak |
| L2048N224 | 2048 | 224 | ~112-bit | Deprecated |
| L2048N256 | 2048 | 256 | ~128-bit | Deprecated |
| L3072N256 | 3072 | 256 | ~128-bit | Deprecated |

## Modern Alternatives Comparison

| Algorithm | Key Size | Speed | Security | NIST | Recommendation |
|-----------|----------|-------|----------|------|----------------|
| DSA | 2048-3072 bits | Slow | Weak RNG risk | Yes (legacy) | ❌ Deprecated |
| ECDSA | 256-521 bits | Fast | Strong | Yes | ✅ Good |
| Ed25519 | 256 bits | Very Fast | Strong | No | ✅✅ Best |

## Security Warnings

1. **Weak RNG** - DSA is critically vulnerable to weak random number generators
2. **Nonce Reuse** - Reusing nonce reveals private key
3. **L1024N160** - Should never be used (broken security)
4. **Timing Attacks** - More vulnerable than modern algorithms
5. **Limited Lifespan** - Support being removed from many systems

## Migration Timeline

**Recommended Migration Schedule:**

1. **Immediate**: Stop using DSA for new applications
2. **Short-term** (< 6 months): Plan migration to Ed25519/ECDSA
3. **Medium-term** (6-12 months): Begin migration of existing systems
4. **Long-term** (12-24 months): Complete migration, remove DSA

## Dependencies

- `crypto/dsa` - Go standard library
- `crypto/rand` - Cryptographic random generator
- `crypto/sha256` - Message hashing

## Further Reading

- [FIPS 186-4](https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.186-4.pdf)
- [DSA Deprecation](https://csrc.nist.gov/projects/digital-signature-standard)
- [Ed25519 Migration](https://ed25519.cr.yp.to/)
- [ECDSA Alternative](https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.186-5.pdf)
