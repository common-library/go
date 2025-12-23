# ECDSA

Elliptic Curve Digital Signature Algorithm for efficient digital signatures.

## Overview

The ecdsa package provides ECDSA digital signatures using NIST-approved elliptic curves. It offers smaller key sizes and faster operations compared to RSA while maintaining equivalent security levels.

## Features

- **Multiple Curves** - P-256, P-384, P-521 support
- **Fast Signatures** - Faster than RSA
- **Small Keys** - 256-bit key ≈ 128-bit security (vs 3072-bit RSA)
- **NIST Standard** - FIPS 186-4 compliant
- **Key Formats** - PEM/PKCS8/PKIX support

## Installation

```bash
go get -u github.com/common-library/go/security/crypto/ecdsa
```

## Quick Start

```go
import (
    "crypto/elliptic"
    "github.com/common-library/go/security/crypto/ecdsa"
)

// Generate key pair
keyPair := &ecdsa.KeyPair{}
keyPair.Generate(elliptic.P256())

// Sign message
signature, _ := keyPair.Sign("message")

// Verify signature
valid := keyPair.Verify("message", signature)
```

## API Reference

### KeyPair Methods

#### Generate

```go
func (kp *KeyPair) Generate(curve elliptic.Curve) error
```

Generates ECDSA key pair on specified curve.

**Curves:**
- `elliptic.P256()` - 256-bit (recommended for most applications)
- `elliptic.P384()` - 384-bit (higher security)
- `elliptic.P521()` - 521-bit (maximum security)

#### Sign

```go
func (kp *KeyPair) Sign(message string) (Signature, error)
```

Creates digital signature.

#### Verify

```go
func (kp *KeyPair) Verify(message string, signature Signature) bool
```

Verifies digital signature.

## Examples

### Basic Signing

```go
package main

import (
    "crypto/elliptic"
    "fmt"
    "log"
    "github.com/common-library/go/security/crypto/ecdsa"
)

func main() {
    // Generate P-256 key pair
    keyPair := &ecdsa.KeyPair{}
    err := keyPair.Generate(elliptic.P256())
    if err != nil {
        log.Fatal(err)
    }
    
    // Sign message
    message := "Important message"
    signature, err := keyPair.Sign(message)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Signature R: %x\n", signature.R)
    fmt.Printf("Signature S: %x\n", signature.S)
    
    // Verify signature
    valid := keyPair.Verify(message, signature)
    if valid {
        fmt.Println("✓ Signature is valid")
    } else {
        fmt.Println("✗ Signature is invalid")
    }
}
```

### Different Curve Sizes

```go
package main

import (
    "crypto/elliptic"
    "fmt"
    "github.com/common-library/go/security/crypto/ecdsa"
)

func main() {
    curves := []struct {
        name  string
        curve elliptic.Curve
    }{
        {"P-256", elliptic.P256()},
        {"P-384", elliptic.P384()},
        {"P-521", elliptic.P521()},
    }
    
    for _, c := range curves {
        keyPair := &ecdsa.KeyPair{}
        keyPair.Generate(c.curve)
        
        privateKey, publicKey := keyPair.GetKeyPair()
        privKey := privateKey.Get()
        pubKey := publicKey.Get()
        
        fmt.Printf("%s:\n", c.name)
        fmt.Printf("  Private key size: %d bits\n", privKey.Params().BitSize)
        fmt.Printf("  Public key X: %d bytes\n", len(pubKey.X.Bytes()))
        fmt.Printf("  Public key Y: %d bytes\n", len(pubKey.Y.Bytes()))
    }
}
```

### Key Export and Import

```go
package main

import (
    "crypto/elliptic"
    "log"
    "os"
    "github.com/common-library/go/security/crypto/ecdsa"
)

func main() {
    // Generate and save
    keyPair := &ecdsa.KeyPair{}
    keyPair.Generate(elliptic.P256())
    
    privateKey, publicKey := keyPair.GetKeyPair()
    
    privatePem, _ := privateKey.GetPemPKCS8()
    publicPem, _ := publicKey.GetPemPKIX()
    
    os.WriteFile("ecdsa_private.pem", []byte(privatePem), 0600)
    os.WriteFile("ecdsa_public.pem", []byte(publicPem), 0644)
    
    // Load and use
    privatePem2, _ := os.ReadFile("ecdsa_private.pem")
    newPrivate := &ecdsa.PrivateKey{}
    newPrivate.SetPemPKCS8(string(privatePem2))
    
    signature, _ := newPrivate.Sign("test message")
    log.Printf("Signature created with loaded key: %v", signature)
}
```

### Message Authentication

```go
package main

import (
    "crypto/elliptic"
    "fmt"
    "github.com/common-library/go/security/crypto/ecdsa"
)

func authenticateMessage(message string, signature ecdsa.Signature, publicKeyPem string) bool {
    publicKey := &ecdsa.PublicKey{}
    err := publicKey.SetPemPKIX(publicKeyPem)
    if err != nil {
        return false
    }
    
    return publicKey.Verify(message, signature)
}

func main() {
    // Sender
    keyPair := &ecdsa.KeyPair{}
    keyPair.Generate(elliptic.P256())
    
    message := "Transfer $100 to account 12345"
    signature, _ := keyPair.Sign(message)
    
    _, publicKey := keyPair.GetKeyPair()
    publicPem, _ := publicKey.GetPemPKIX()
    
    // Receiver
    valid := authenticateMessage(message, signature, publicPem)
    if valid {
        fmt.Println("✓ Message authenticated")
    } else {
        fmt.Println("✗ Authentication failed")
    }
    
    // Tampered message
    tamperedValid := authenticateMessage("Transfer $999 to account 12345", signature, publicPem)
    if !tamperedValid {
        fmt.Println("✓ Tampered message detected")
    }
}
```

## Best Practices

### 1. Choose Appropriate Curve

```go
// Good: P-256 for most applications
keyPair.Generate(elliptic.P256())

// Good: P-384 for higher security requirements
keyPair.Generate(elliptic.P384())

// Consider: P-521 only if maximum security needed
keyPair.Generate(elliptic.P521())
```

### 2. Verify All Signatures

```go
// Good: Always verify before trusting
if keyPair.Verify(message, signature) {
    processMessage(message)
}

// Avoid: Process without verification
processMessage(message)
```

### 3. Use Non-Deterministic Signatures

```go
// ECDSA signatures are naturally non-deterministic
// Each signature for the same message will be different
// This is a security feature
```

## Curve Comparison

| Curve | Key Size | Security Level | Performance |
|-------|----------|----------------|-------------|
| P-256 | 256 bits | 128-bit | Fast |
| P-384 | 384 bits | 192-bit | Moderate |
| P-521 | 521 bits | 256-bit | Slower |

**Equivalent RSA Key Sizes:**
- P-256 ≈ RSA 3072-bit
- P-384 ≈ RSA 7680-bit
- P-521 ≈ RSA 15360-bit

## Security Considerations

- **Randomness**: Each signature uses random nonce (automatically handled)
- **Signature Verification**: Always required before trusting data
- **Curve Choice**: P-256 suitable for most applications
- **Key Protection**: Protect private keys with file permissions
- **Message Hashing**: Messages are automatically hashed before signing

## Performance Tips

1. **Reuse Keys** - Generate once, use many times
2. **P-256 for Speed** - Fastest curve with good security
3. **Batch Verification** - Verify multiple signatures in parallel
4. **Cache Public Keys** - Store frequently used public keys

## ECDSA vs Other Algorithms

| Feature | ECDSA | RSA | Ed25519 |
|---------|-------|-----|---------|
| Key Size | Small | Large | Very Small |
| Speed | Fast | Slow | Very Fast |
| Signatures | Non-deterministic | Deterministic (PSS) | Deterministic |
| NIST Approved | Yes | Yes | No |
| Use Case | General | Legacy/Compatibility | Modern Apps |

## Dependencies

- `crypto/ecdsa` - Go standard library
- `crypto/elliptic` - Elliptic curve operations
- `crypto/rand` - Cryptographic random generator
- `crypto/x509` - X.509 encoding

## Further Reading

- [ECDSA](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm)
- [NIST Curves](https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.186-4.pdf)
- [Go crypto/ecdsa Package](https://pkg.go.dev/crypto/ecdsa)
