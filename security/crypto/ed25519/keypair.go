// Package ed25519 provides Ed25519 digital signature cryptography.
//
// This package implements Ed25519 public-key signature system based on the
// elliptic curve Ed25519. It provides key generation, signing, and verification
// operations with simplified interfaces for key pair management.
//
// # Features
//
//   - Ed25519 key pair generation
//   - Digital signature creation and verification
//   - PEM/PKCS8/PKIX format support
//   - SSH public key format conversion
//   - Type-safe key management
//
// # Basic Example
//
//	keyPair := &ed25519.KeyPair{}
//	err := keyPair.Generate()
//	signature := keyPair.Sign("message")
//	valid := keyPair.Verify("message", signature)
package ed25519

// KeyPair is struct that provides key pair related methods.
type KeyPair struct {
	privateKey PrivateKey
	publicKey  PublicKey
}

// Generate creates a new Ed25519 key pair.
//
// This method generates a random Ed25519 private key and derives the
// corresponding public key using the system's cryptographic random generator.
//
// # Returns
//
//   - error: Error if key generation fails, nil on success
//
// # Behavior
//
// The generated key pair:
//   - Private key: 64 bytes (32-byte seed + 32-byte public key)
//   - Public key: 32 bytes
//   - Uses crypto/rand for secure randomness
//   - Automatically derives public key from private key
//
// # Examples
//
// Basic usage:
//
//	keyPair := &ed25519.KeyPair{}
//	err := keyPair.Generate()
//	if err != nil {
//	    log.Fatal(err)
//	}
func (kp *KeyPair) Generate() error {
	if err := kp.privateKey.SetDefault(); err != nil {
		return err
	} else {
		kp.publicKey = kp.privateKey.GetPublicKey()
		return nil
	}
}

// Sign creates a digital signature for the message.
//
// This method signs the message using the private key, producing a
// 64-byte signature that can be verified with the corresponding public key.
//
// # Parameters
//
//   - message: Text message to sign
//
// # Returns
//
//   - []byte: 64-byte Ed25519 signature
//
// # Behavior
//
// The signature:
//   - Always 64 bytes regardless of message length
//   - Deterministic for same message and key
//   - Cryptographically secure (128-bit security level)
//   - Cannot be forged without the private key
//
// # Examples
//
// Sign a message:
//
//	signature := keyPair.Sign("Hello, World!")
//	fmt.Printf("Signature: %x\n", signature)
func (kp *KeyPair) Sign(message string) []byte {
	return kp.privateKey.Sign(message)
}

// Verify verifies a digital signature.
//
// This method verifies that the signature was created by the private key
// corresponding to this key pair's public key.
//
// # Parameters
//
//   - message: Original text message
//   - signature: 64-byte signature to verify
//
// # Returns
//
//   - bool: true if signature is valid, false otherwise
//
// # Behavior
//
// Verification checks:
//   - Signature matches message and public key
//   - Signature has not been tampered with
//   - Returns false for invalid signatures (no panic)
//
// # Examples
//
// Verify a signature:
//
//	signature := keyPair.Sign("message")
//	valid := keyPair.Verify("message", signature)
//	if valid {
//	    fmt.Println("Signature is valid")
//	}
func (kp *KeyPair) Verify(message string, signature []byte) bool {
	return kp.publicKey.Verify(message, signature)
}

// GetKeyPair retrieves the private and public keys.
//
// This method returns both keys from the key pair for separate operations
// or storage.
//
// # Returns
//
//   - privateKey: PrivateKey for signing operations
//   - publicKey: PublicKey for verification operations
//
// # Examples
//
// Extract keys:
//
//	privateKey, publicKey := keyPair.GetKeyPair()
//	pemPrivate, _ := privateKey.GetPemPKCS8()
//	pemPublic, _ := publicKey.GetPemPKIX()
func (kp *KeyPair) GetKeyPair() (privateKey PrivateKey, publicKey PublicKey) {
	return kp.privateKey, kp.publicKey
}

// SetKeyPair sets the private and public keys.
//
// This method initializes the key pair with existing keys, useful for
// loading keys from storage or external sources.
//
// # Parameters
//
//   - privateKey: PrivateKey to set
//   - publicKey: PublicKey to set (must correspond to private key)
//
// # Examples
//
// Load existing keys:
//
//	privateKey := ed25519.PrivateKey{}
//	privateKey.SetPemPKCS8(pemString)
//	publicKey := privateKey.GetPublicKey()
//	keyPair.SetKeyPair(privateKey, publicKey)
func (kp *KeyPair) SetKeyPair(privateKey PrivateKey, publicKey PublicKey) {
	kp.privateKey = privateKey
	kp.publicKey = publicKey
}
