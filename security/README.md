# Security

Cryptographic utilities for secure applications.

## Overview

The security package provides cryptographic primitives and digital signature implementations for secure communication and data authentication.

## Subpackages

### crypto/ed25519

Ed25519 digital signature system for high-performance public-key cryptography.

[📖 Documentation](crypto/ed25519/README.md)

**Features:**
- Fast key generation with secure randomness
- 64-byte deterministic signatures
- 128-bit security level
- PEM/PKCS8/PKIX format support
- SSH public key format conversion

**Quick Example:**
```go
import "github.com/common-library/go/security/crypto/ed25519"

// Generate key pair
keyPair := &ed25519.KeyPair{}
keyPair.Generate()

// Sign message
signature := keyPair.Sign("Hello, World!")

// Verify signature
valid := keyPair.Verify("Hello, World!", signature)
```

### crypto/rsa

RSA public-key cryptography for encryption and signatures.

**Features:**
- Configurable key sizes (2048, 3072, 4096 bits)
- PKCS1v15 and OAEP padding
- Digital signatures with PSS
- PEM format import/export

### crypto/ecdsa

Elliptic Curve Digital Signature Algorithm for efficient signatures.

**Features:**
- Multiple curves (P-256, P-384, P-521)
- Smaller key sizes than RSA
- Fast signature generation
- Modern cryptographic standard

### crypto/dsa

Digital Signature Algorithm (legacy support).

**Features:**
- DSA signatures
- Parameter generation
- Legacy system compatibility

## Algorithm Comparison

| Algorithm | Key Size | Signature Size | Security Level | Speed | Use Case |
|-----------|----------|----------------|----------------|-------|----------|
| Ed25519 | 32 bytes | 64 bytes | 128-bit | Very Fast | Modern apps |
| RSA-2048 | 256 bytes | 256 bytes | 112-bit | Moderate | Legacy compatibility |
| RSA-3072 | 384 bytes | 384 bytes | 128-bit | Slow | High security |
| ECDSA P-256 | 32 bytes | 64 bytes | 128-bit | Fast | Modern apps |
| DSA | Variable | 40-64 bytes | Variable | Moderate | Legacy only |

## Installation

```bash
go get -u github.com/common-library/go/security/crypto/ed25519
go get -u github.com/common-library/go/security/crypto/rsa
go get -u github.com/common-library/go/security/crypto/ecdsa
go get -u github.com/common-library/go/security/crypto/dsa
```

## Choosing an Algorithm

### Use Ed25519 when:
- Building new applications
- Need maximum performance
- Want small key/signature sizes
- Require deterministic signatures

### Use RSA when:
- Interoperating with legacy systems
- Need encryption (not just signatures)
- Required by compliance/standards

### Use ECDSA when:
- Need NIST-approved curves
- Interoperating with systems requiring ECDSA
- Want smaller keys than RSA

### Avoid DSA:
- Considered legacy
- Use Ed25519 or ECDSA instead

## Security Best Practices

1. **Key Protection** - Never expose private keys
2. **Secure Storage** - Encrypt private keys at rest
3. **Key Rotation** - Regularly rotate signing keys
4. **Randomness** - Use crypto/rand for key generation
5. **Algorithm Selection** - Prefer Ed25519 for new projects
6. **Signature Verification** - Always verify before trusting data
7. **Key Sizes** - Use minimum 2048-bit RSA, prefer 3072+ for long-term security

## Common Use Cases

### Message Authentication

```go
// Sender
keyPair := &ed25519.KeyPair{}
keyPair.Generate()
signature := keyPair.Sign(message)

_, publicKey := keyPair.GetKeyPair()
publicPem, _ := publicKey.GetPemPKIX()
// Share publicPem with receiver

// Receiver
receivedPublicKey := &ed25519.PublicKey{}
receivedPublicKey.SetPemPKIX(publicPem)
valid := receivedPublicKey.Verify(message, signature)
```

### Key Storage

```go
// Save keys
privateKey, publicKey := keyPair.GetKeyPair()
privatePem, _ := privateKey.GetPemPKCS8()
publicPem, _ := publicKey.GetPemPKIX()

os.WriteFile("private.pem", []byte(privatePem), 0600)
os.WriteFile("public.pem", []byte(publicPem), 0644)

// Load keys
privatePem, _ := os.ReadFile("private.pem")
privateKey := &ed25519.PrivateKey{}
privateKey.SetPemPKCS8(string(privatePem))
```

### SSH Key Generation

```go
keyPair := &ed25519.KeyPair{}
keyPair.Generate()

_, publicKey := keyPair.GetKeyPair()
sshKey, _ := publicKey.GetSsh()

os.WriteFile("id_ed25519.pub", []byte(sshKey), 0644)
```

## Dependencies

- `crypto/ed25519` - Go standard library
- `crypto/rsa` - Go standard library
- `crypto/ecdsa` - Go standard library
- `crypto/dsa` - Go standard library
- `golang.org/x/crypto/ssh` - SSH key formats

## Further Reading

- [Ed25519 Documentation](crypto/ed25519/README.md)
- [Go Cryptography Documentation](https://pkg.go.dev/crypto)
- [NIST Cryptographic Standards](https://csrc.nist.gov/)
