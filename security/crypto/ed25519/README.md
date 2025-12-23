# Ed25519

Ed25519 digital signature cryptography implementation.

## Overview

The ed25519 package provides a high-level interface for Ed25519 public-key signatures based on the elliptic curve Ed25519. It simplifies key generation, signing, and verification operations with support for multiple key formats including PEM, PKCS8, PKIX, and SSH.

## Features

- **Key Pair Generation** - Cryptographically secure random key generation
- **Digital Signatures** - Sign and verify messages with 128-bit security
- **Multiple Formats** - PEM/PKCS8/PKIX/SSH format support
- **Type Safety** - Separate types for private and public keys
- **Key Derivation** - Automatic public key derivation from private key

## Installation

```bash
go get -u github.com/common-library/go/security/crypto/ed25519
```

## Quick Start

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

## API Reference

### KeyPair Type

```go
type KeyPair struct {
    privateKey PrivateKey
    publicKey  PublicKey
}
```

Container for Ed25519 private and public key pair.

### KeyPair Methods

#### Generate

```go
func (kp *KeyPair) Generate() error
```

Generates a new random Ed25519 key pair.

#### Sign

```go
func (kp *KeyPair) Sign(message string) []byte
```

Creates a 64-byte signature for the message.

#### Verify

```go
func (kp *KeyPair) Verify(message string, signature []byte) bool
```

Verifies a signature matches the message.

#### GetKeyPair

```go
func (kp *KeyPair) GetKeyPair() (privateKey PrivateKey, publicKey PublicKey)
```

Retrieves the private and public keys.

#### SetKeyPair

```go
func (kp *KeyPair) SetKeyPair(privateKey PrivateKey, publicKey PublicKey)
```

Sets the private and public keys.

### PrivateKey Type

```go
type PrivateKey struct {
    privateKey ed25519.PrivateKey
    publicKey  ed25519.PublicKey
}
```

Ed25519 private key (64 bytes: 32-byte seed + 32-byte public key).

### PrivateKey Methods

#### Sign

```go
func (pk *PrivateKey) Sign(message string) []byte
```

Creates a digital signature.

#### Verify

```go
func (pk *PrivateKey) Verify(message string, signature []byte) bool
```

Verifies a signature using the associated public key.

#### Get

```go
func (pk *PrivateKey) Get() ed25519.PrivateKey
```

Gets the underlying ed25519.PrivateKey.

#### Set

```go
func (pk *PrivateKey) Set(privateKey ed25519.PrivateKey)
```

Sets the private key.

#### SetDefault

```go
func (pk *PrivateKey) SetDefault() error
```

Generates a new random private key.

#### GetPemPKCS8

```go
func (pk *PrivateKey) GetPemPKCS8() (string, error)
```

Returns the private key in PEM-encoded PKCS#8 format.

#### SetPemPKCS8

```go
func (pk *PrivateKey) SetPemPKCS8(pemPKCS8 string) error
```

Sets the private key from PEM-encoded PKCS#8.

#### GetPublicKey

```go
func (pk *PrivateKey) GetPublicKey() PublicKey
```

Derives the public key from the private key.

### PublicKey Type

```go
type PublicKey struct {
    publicKey ed25519.PublicKey
}
```

Ed25519 public key (32 bytes).

### PublicKey Methods

#### Verify

```go
func (pk *PublicKey) Verify(message string, signature []byte) bool
```

Verifies a digital signature.

#### Get

```go
func (pk *PublicKey) Get() ed25519.PublicKey
```

Gets the underlying ed25519.PublicKey.

#### Set

```go
func (pk *PublicKey) Set(publicKey ed25519.PublicKey)
```

Sets the public key.

#### GetPemPKIX

```go
func (pk *PublicKey) GetPemPKIX() (string, error)
```

Returns the public key in PEM-encoded PKIX format.

#### SetPemPKIX

```go
func (pk *PublicKey) SetPemPKIX(pemPKIX string) error
```

Sets the public key from PEM-encoded PKIX.

#### GetSsh

```go
func (pk *PublicKey) GetSsh() (string, error)
```

Returns the public key in SSH authorized_keys format.

#### SetSsh

```go
func (pk *PublicKey) SetSsh(sshKey string) error
```

Sets the public key from SSH format.

#### GetSshPublicKey

```go
func (pk *PublicKey) GetSshPublicKey() (ssh.PublicKey, error)
```

Gets the key as ssh.PublicKey.

#### SetSshPublicKey

```go
func (pk *PublicKey) SetSshPublicKey(publicKey ssh.PublicKey) error
```

Sets the key from ssh.PublicKey.

## Complete Examples

### Basic Signing and Verification

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/security/crypto/ed25519"
)

func main() {
    // Generate key pair
    keyPair := &ed25519.KeyPair{}
    err := keyPair.Generate()
    if err != nil {
        log.Fatal(err)
    }
    
    // Sign message
    message := "Hello, World!"
    signature := keyPair.Sign(message)
    fmt.Printf("Signature: %x\n", signature)
    
    // Verify signature
    valid := keyPair.Verify(message, signature)
    if valid {
        fmt.Println("✓ Signature is valid")
    } else {
        fmt.Println("✗ Signature is invalid")
    }
}
```

### Exporting and Importing Keys (PEM Format)

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/security/crypto/ed25519"
)

func main() {
    // Generate key pair
    keyPair := &ed25519.KeyPair{}
    keyPair.Generate()
    
    // Extract keys
    privateKey, publicKey := keyPair.GetKeyPair()
    
    // Export to PEM
    privatePem, err := privateKey.GetPemPKCS8()
    if err != nil {
        log.Fatal(err)
    }
    
    publicPem, err := publicKey.GetPemPKIX()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Private Key (PEM):")
    fmt.Println(privatePem)
    
    fmt.Println("Public Key (PEM):")
    fmt.Println(publicPem)
    
    // Import from PEM
    newPrivateKey := &ed25519.PrivateKey{}
    err = newPrivateKey.SetPemPKCS8(privatePem)
    if err != nil {
        log.Fatal(err)
    }
    
    newPublicKey := &ed25519.PublicKey{}
    err = newPublicKey.SetPemPKIX(publicPem)
    if err != nil {
        log.Fatal(err)
    }
    
    // Verify imported keys work
    signature := newPrivateKey.Sign("test")
    valid := newPublicKey.Verify("test", signature)
    fmt.Printf("Imported keys valid: %v\n", valid)
}
```

### SSH Key Format

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/security/crypto/ed25519"
)

func main() {
    // Generate key pair
    keyPair := &ed25519.KeyPair{}
    keyPair.Generate()
    
    _, publicKey := keyPair.GetKeyPair()
    
    // Export to SSH format
    sshKey, err := publicKey.GetSsh()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("SSH Public Key:")
    fmt.Println(sshKey)
    
    // Import from SSH format
    newPublicKey := &ed25519.PublicKey{}
    err = newPublicKey.SetSsh(sshKey)
    if err != nil {
        log.Fatal(err)
    }
    
    // Verify works
    signature := keyPair.Sign("test")
    valid := newPublicKey.Verify("test", signature)
    fmt.Printf("SSH key import successful: %v\n", valid)
}
```

### Separate Private and Public Key Operations

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/security/crypto/ed25519"
)

func main() {
    // Generate private key
    privateKey := &ed25519.PrivateKey{}
    err := privateKey.SetDefault()
    if err != nil {
        log.Fatal(err)
    }
    
    // Derive public key
    publicKey := privateKey.GetPublicKey()
    
    // Sign with private key
    message := "Secure message"
    signature := privateKey.Sign(message)
    
    // Verify with public key
    valid := publicKey.Verify(message, signature)
    fmt.Printf("Signature valid: %v\n", valid)
    
    // Export public key for sharing
    publicPem, _ := publicKey.GetPemPKIX()
    fmt.Println("\nShare this public key:")
    fmt.Println(publicPem)
}
```

### Message Authentication

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/security/crypto/ed25519"
)

func authenticateMessage(message string, signature []byte, publicKeyPem string) bool {
    publicKey := &ed25519.PublicKey{}
    err := publicKey.SetPemPKIX(publicKeyPem)
    if err != nil {
        log.Printf("Invalid public key: %v", err)
        return false
    }
    
    return publicKey.Verify(message, signature)
}

func main() {
    // Sender generates key pair
    keyPair := &ed25519.KeyPair{}
    keyPair.Generate()
    
    _, publicKey := keyPair.GetKeyPair()
    publicPem, _ := publicKey.GetPemPKIX()
    
    // Sender signs message
    message := "Authenticated message"
    signature := keyPair.Sign(message)
    
    // Receiver verifies message
    valid := authenticateMessage(message, signature, publicPem)
    if valid {
        fmt.Println("✓ Message authenticated")
    } else {
        fmt.Println("✗ Message authentication failed")
    }
    
    // Try with tampered message
    tamperedMessage := "Tampered message"
    valid = authenticateMessage(tamperedMessage, signature, publicPem)
    if !valid {
        fmt.Println("✓ Tampered message detected")
    }
}
```

### Key Storage and Loading

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/common-library/go/security/crypto/ed25519"
)

func saveKeys(privateKeyPath, publicKeyPath string, keyPair *ed25519.KeyPair) error {
    privateKey, publicKey := keyPair.GetKeyPair()
    
    // Save private key
    privatePem, err := privateKey.GetPemPKCS8()
    if err != nil {
        return err
    }
    
    err = os.WriteFile(privateKeyPath, []byte(privatePem), 0600)
    if err != nil {
        return err
    }
    
    // Save public key
    publicPem, err := publicKey.GetPemPKIX()
    if err != nil {
        return err
    }
    
    err = os.WriteFile(publicKeyPath, []byte(publicPem), 0644)
    if err != nil {
        return err
    }
    
    return nil
}

func loadKeys(privateKeyPath, publicKeyPath string) (*ed25519.KeyPair, error) {
    // Load private key
    privatePem, err := os.ReadFile(privateKeyPath)
    if err != nil {
        return nil, err
    }
    
    privateKey := &ed25519.PrivateKey{}
    err = privateKey.SetPemPKCS8(string(privatePem))
    if err != nil {
        return nil, err
    }
    
    // Load public key
    publicPem, err := os.ReadFile(publicKeyPath)
    if err != nil {
        return nil, err
    }
    
    publicKey := &ed25519.PublicKey{}
    err = publicKey.SetPemPKIX(string(publicPem))
    if err != nil {
        return nil, err
    }
    
    // Create key pair
    keyPair := &ed25519.KeyPair{}
    keyPair.SetKeyPair(*privateKey, *publicKey)
    
    return keyPair, nil
}

func main() {
    // Generate and save keys
    keyPair := &ed25519.KeyPair{}
    keyPair.Generate()
    
    err := saveKeys("private.pem", "public.pem", keyPair)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Keys saved")
    
    // Load keys
    loadedKeyPair, err := loadKeys("private.pem", "public.pem")
    if err != nil {
        log.Fatal(err)
    }
    
    // Verify loaded keys work
    signature := loadedKeyPair.Sign("test")
    valid := loadedKeyPair.Verify("test", signature)
    fmt.Printf("Loaded keys work: %v\n", valid)
}
```

## Best Practices

### 1. Protect Private Keys

```go
// Good: Restrict private key file permissions
os.WriteFile("private.pem", privatePem, 0600) // Owner read/write only

// Avoid: World-readable private keys
os.WriteFile("private.pem", privatePem, 0644) // Insecure
```

### 2. Verify Message Integrity

```go
// Good: Always verify before trusting message
if publicKey.Verify(message, signature) {
    processMessage(message)
}

// Avoid: Trust without verification
processMessage(message) // Dangerous!
```

### 3. Use Deterministic Signatures

```go
// Good: Ed25519 signatures are deterministic
sig1 := privateKey.Sign("message")
sig2 := privateKey.Sign("message")
// sig1 == sig2 for same key and message

// This is a security feature - no random nonces needed
```

### 4. Handle Errors Properly

```go
// Good: Check all error returns
privatePem, err := privateKey.GetPemPKCS8()
if err != nil {
    return fmt.Errorf("failed to export key: %w", err)
}

// Avoid: Ignore errors
privatePem, _ := privateKey.GetPemPKCS8()
```

## Performance Tips

1. **Reuse Key Pairs** - Generate once, reuse many times
2. **Public Key Distribution** - Share public keys freely, signatures are fast to verify
3. **Batch Verification** - Verify multiple signatures in parallel for throughput
4. **Key Storage** - Cache loaded keys in memory during application lifetime

## Security Considerations

1. **Private Key Protection** - Never expose private keys
2. **Message Hashing** - Ed25519 can sign messages of any length directly
3. **Signature Size** - Always 64 bytes regardless of message size
4. **Security Level** - 128-bit security (equivalent to 3072-bit RSA)
5. **Deterministic** - Same message and key always produce same signature

## Testing

```go
func TestEd25519(t *testing.T) {
    keyPair := &ed25519.KeyPair{}
    err := keyPair.Generate()
    if err != nil {
        t.Fatalf("Failed to generate: %v", err)
    }
    
    message := "test message"
    signature := keyPair.Sign(message)
    
    if len(signature) != 64 {
        t.Errorf("Expected 64-byte signature, got %d", len(signature))
    }
    
    valid := keyPair.Verify(message, signature)
    if !valid {
        t.Error("Signature verification failed")
    }
    
    // Test wrong message
    invalid := keyPair.Verify("wrong message", signature)
    if invalid {
        t.Error("Should reject wrong message")
    }
}
```

## Dependencies

- `crypto/ed25519` - Go standard library
- `crypto/rand` - Cryptographic random generator
- `crypto/x509` - X.509 encoding
- `encoding/pem` - PEM encoding
- `golang.org/x/crypto/ssh` - SSH key formats

## Further Reading

- [Ed25519 Specification](https://ed25519.cr.yp.to/)
- [RFC 8032](https://tools.ietf.org/html/rfc8032)
- [Go crypto/ed25519 Package](https://pkg.go.dev/crypto/ed25519)
