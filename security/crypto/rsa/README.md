# RSA

RSA public-key cryptography for encryption and digital signatures.

## Overview

The rsa package provides RSA encryption and digital signature operations with support for multiple key sizes, PKCS1v15 and OAEP padding schemes, and PEM format key storage.

## Features

- **Key Generation** - 2048, 3072, 4096-bit keys
- **Encryption** - PKCS1v15 and OAEP padding
- **Digital Signatures** - PSS padding support
- **Key Formats** - PEM/PKCS1/PKCS8 support
- **SSH Keys** - Public key format conversion

## Installation

```bash
go get -u github.com/common-library/go/security/crypto/rsa
```

## Quick Start

```go
import "github.com/common-library/go/security/crypto/rsa"

// Generate key pair
keyPair := &rsa.KeyPair{}
keyPair.Generate(2048)

// Encrypt (OAEP recommended)
ciphertext, _ := keyPair.EncryptOAEP("secret message")

// Decrypt
plaintext, _ := keyPair.DecryptOAEP(ciphertext)
```

## API Reference

### KeyPair Methods

#### Generate

```go
func (kp *KeyPair) Generate(bits int) error
```

Generates RSA key pair (2048, 3072, or 4096 bits recommended).

#### EncryptPKCS1v15 / DecryptPKCS1v15

```go
func (kp *KeyPair) EncryptPKCS1v15(plaintext string) ([]byte, error)
func (kp *KeyPair) DecryptPKCS1v15(ciphertext []byte) (string, error)
```

PKCS#1 v1.5 encryption (legacy compatibility).

#### EncryptOAEP / DecryptOAEP

```go
func (kp *KeyPair) EncryptOAEP(plaintext string) ([]byte, error)
func (kp *KeyPair) DecryptOAEP(ciphertext []byte) (string, error)
```

OAEP encryption (recommended for new applications).

## Examples

### Basic Encryption

```go
package main

import (
    "fmt"
    "log"
    "github.com/common-library/go/security/crypto/rsa"
)

func main() {
    keyPair := &rsa.KeyPair{}
    err := keyPair.Generate(2048)
    if err != nil {
        log.Fatal(err)
    }
    
    // Encrypt with OAEP (recommended)
    message := "Secret message"
    ciphertext, err := keyPair.EncryptOAEP(message)
    if err != nil {
        log.Fatal(err)
    }
    
    // Decrypt
    plaintext, err := keyPair.DecryptOAEP(ciphertext)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Original: %s\n", message)
    fmt.Printf("Decrypted: %s\n", plaintext)
}
```

### Key Export and Import

```go
package main

import (
    "log"
    "os"
    "github.com/common-library/go/security/crypto/rsa"
)

func main() {
    // Generate and save keys
    keyPair := &rsa.KeyPair{}
    keyPair.Generate(2048)
    
    privateKey, publicKey := keyPair.GetKeyPair()
    
    privatePem, _ := privateKey.GetPemPKCS1()
    publicPem, _ := publicKey.GetPemPKIX()
    
    os.WriteFile("private.pem", []byte(privatePem), 0600)
    os.WriteFile("public.pem", []byte(publicPem), 0644)
    
    // Load keys
    privatePem2, _ := os.ReadFile("private.pem")
    newPrivate := &rsa.PrivateKey{}
    newPrivate.SetPemPKCS1(string(privatePem2))
    
    publicPem2, _ := os.ReadFile("public.pem")
    newPublic := &rsa.PublicKey{}
    newPublic.SetPemPKIX(string(publicPem2))
    
    // Use loaded keys
    newKeyPair := &rsa.KeyPair{}
    newKeyPair.SetKeyPair(*newPrivate, *newPublic)
    
    ciphertext, _ := newKeyPair.EncryptOAEP("test")
    plaintext, _ := newKeyPair.DecryptOAEP(ciphertext)
    log.Printf("Decrypted: %s", plaintext)
}
```

## Best Practices

### 1. Use OAEP Padding

```go
// Good: OAEP is more secure
ciphertext, _ := keyPair.EncryptOAEP(message)

// Avoid: PKCS1v15 is legacy
ciphertext, _ := keyPair.EncryptPKCS1v15(message)
```

### 2. Use Adequate Key Sizes

```go
// Good: 2048 bits minimum, 3072+ for long-term security
keyPair.Generate(3072)

// Avoid: 1024 bits is insecure
keyPair.Generate(1024)
```

### 3. Protect Private Keys

```go
// Good: Restrict permissions
os.WriteFile("private.pem", privatePem, 0600)

// Avoid: World-readable
os.WriteFile("private.pem", privatePem, 0644)
```

## Security Considerations

- **Key Size**: Use 2048 bits minimum (3072+ for sensitive data)
- **Padding**: Prefer OAEP over PKCS1v15
- **Random Number Generation**: Uses crypto/rand automatically
- **Message Size**: Limited by key size (e.g., 2048-bit key ≈ 245 bytes with OAEP)
- **Performance**: RSA is slower than symmetric encryption; use for key exchange

## Dependencies

- `crypto/rsa` - Go standard library
- `crypto/rand` - Cryptographic random generator
- `crypto/x509` - X.509 encoding
- `encoding/pem` - PEM encoding

## Further Reading

- [RSA Cryptosystem](https://en.wikipedia.org/wiki/RSA_(cryptosystem))
- [OAEP Padding](https://en.wikipedia.org/wiki/Optimal_asymmetric_encryption_padding)
- [Go crypto/rsa Package](https://pkg.go.dev/crypto/rsa)
