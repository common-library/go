// Package rsa provides rsa crypto related implementations.
package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

// PrivateKey is struct that provides private key related methods.
type PrivateKey struct {
	privateKey *rsa.PrivateKey
}

// EncryptPKCS1v15 encrypts plaintext using the public key with PKCS#1 v1.5 padding.
//
// This method uses the public key portion of the private key to encrypt data.
// PKCS#1 v1.5 padding is a legacy scheme and should be avoided for new applications.
// Use EncryptOAEP instead for better security.
//
// Parameters:
//   - plaintext: The string to encrypt
//
// Returns:
//   - []byte: The encrypted ciphertext
//   - error: Error if encryption fails or plaintext is too long
//
// Behavior:
//   - Maximum plaintext size depends on key size and padding overhead
//   - Uses cryptographically secure random padding
//   - Same plaintext produces different ciphertext each time (non-deterministic)
//
// Example:
//
//	privateKey := &rsa.PrivateKey{}
//	privateKey.SetBits(2048)
//	ciphertext, err := privateKey.EncryptPKCS1v15("secret message")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PrivateKey) EncryptPKCS1v15(plaintext string) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, &pk.privateKey.PublicKey, []byte(plaintext))
}

// DecryptPKCS1v15 decrypts ciphertext using PKCS#1 v1.5 padding.
//
// This method decrypts data that was encrypted with the corresponding public key
// using PKCS#1 v1.5 padding. This padding scheme has known vulnerabilities and
// should only be used for compatibility with legacy systems.
//
// Parameters:
//   - ciphertext: The encrypted data to decrypt
//
// Returns:
//   - string: The decrypted plaintext
//   - error: Error if decryption fails or ciphertext is invalid
//
// Behavior:
//   - Ciphertext must match the key size (e.g., 256 bytes for 2048-bit key)
//   - Vulnerable to padding oracle attacks in some scenarios
//   - Consider using DecryptOAEP for new applications
//
// Example:
//
//	plaintext, err := privateKey.DecryptPKCS1v15(ciphertext)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Decrypted:", plaintext)
func (pk *PrivateKey) DecryptPKCS1v15(ciphertext []byte) (string, error) {
	if plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, pk.privateKey, ciphertext); err != nil {
		return "", err
	} else {
		return string(plaintext), nil
	}
}

// EncryptOAEP encrypts plaintext using the public key with OAEP padding.
//
// This method uses Optimal Asymmetric Encryption Padding (OAEP), which provides
// better security than PKCS#1 v1.5. OAEP is the recommended padding scheme for
// RSA encryption in new applications. Uses SHA-256 for hashing.
//
// Parameters:
//   - plaintext: The string to encrypt
//
// Returns:
//   - []byte: The encrypted ciphertext
//   - error: Error if encryption fails or plaintext is too long
//
// Behavior:
//   - Maximum plaintext size = key_size - 2*hash_size - 2 (e.g., 190 bytes for 2048-bit key with SHA-256)
//   - Provides provable security against chosen-ciphertext attacks
//   - Non-deterministic (same plaintext produces different ciphertext)
//
// Example:
//
//	privateKey := &rsa.PrivateKey{}
//	privateKey.SetBits(2048)
//	ciphertext, err := privateKey.EncryptOAEP("confidential data")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PrivateKey) EncryptOAEP(plaintext string) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, &pk.privateKey.PublicKey, []byte(plaintext), nil)
}

// DecryptOAEP decrypts ciphertext using OAEP padding.
//
// This method decrypts data that was encrypted with the corresponding public key
// using OAEP padding. OAEP is the recommended padding scheme for RSA encryption
// due to its strong security properties. Uses SHA-256 for hashing.
//
// Parameters:
//   - ciphertext: The encrypted data to decrypt
//
// Returns:
//   - string: The decrypted plaintext
//   - error: Error if decryption fails or ciphertext is invalid
//
// Behavior:
//   - Ciphertext must match the key size (e.g., 256 bytes for 2048-bit key)
//   - Provides protection against chosen-ciphertext attacks
//   - Automatically validates padding during decryption
//
// Example:
//
//	plaintext, err := privateKey.DecryptOAEP(ciphertext)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Decrypted:", plaintext)
func (pk *PrivateKey) DecryptOAEP(ciphertext []byte) (string, error) {
	if plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pk.privateKey, ciphertext, nil); err != nil {
		return "", err
	} else {
		return string(plaintext), nil
	}
}

// Get returns the underlying *rsa.PrivateKey.
//
// This method provides direct access to the Go standard library rsa.PrivateKey
// for use with other crypto functions that require the native type.
//
// Returns:
//   - *rsa.PrivateKey: The underlying private key
//
// Example:
//
//	nativeKey := privateKey.Get()
//	// Use with Go crypto functions
//	x509.MarshalPKCS1PrivateKey(nativeKey)
func (pk *PrivateKey) Get() *rsa.PrivateKey {
	return pk.privateKey
}

// Set assigns an existing *rsa.PrivateKey to this PrivateKey instance.
//
// This method allows setting the private key from an existing Go standard library
// rsa.PrivateKey, useful when loading keys from other sources or libraries.
//
// Parameters:
//   - privateKey: The rsa.PrivateKey to assign
//
// Example:
//
//	nativeKey, _ := rsa.GenerateKey(rand.Reader, 2048)
//	privateKey := &rsa.PrivateKey{}
//	privateKey.Set(nativeKey)
func (pk *PrivateKey) Set(privateKey *rsa.PrivateKey) {
	pk.privateKey = privateKey
}

// SetBits generates a new RSA private key with the specified bit size.
//
// This is a convenience method that generates a new key pair and sets it as
// the current private key. Common bit sizes are 2048, 3072, and 4096.
//
// Parameters:
//   - bits: The size of the key in bits (recommended: 2048 minimum, 3072+ for sensitive data)
//
// Returns:
//   - error: Error if key generation fails
//
// Behavior:
//   - Generates both private and public keys
//   - Uses cryptographically secure random number generator
//   - Larger key sizes provide more security but slower operations
//
// Example:
//
//	privateKey := &rsa.PrivateKey{}
//	err := privateKey.SetBits(2048)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PrivateKey) SetBits(bits int) error {
	if privateKey, err := rsa.GenerateKey(rand.Reader, bits); err != nil {
		return err
	} else {
		pk.Set(privateKey)
		return nil
	}
}

// GetPemPKCS1 encodes the private key as a PEM-encoded PKCS#1 string.
//
// PKCS#1 is the traditional RSA private key format. The PEM encoding allows
// the key to be stored in text files or transmitted as text. The format uses
// "RSA PRIVATE KEY" as the PEM block type.
//
// Returns:
//   - string: PEM-encoded PKCS#1 private key
//
// Behavior:
//   - Output format: "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----"
//   - This format is RSA-specific (not algorithm-agnostic)
//   - Widely supported but less flexible than PKCS#8
//
// Example:
//
//	pemString := privateKey.GetPemPKCS1()
//	os.WriteFile("private_key.pem", []byte(pemString), 0600)
func (pk *PrivateKey) GetPemPKCS1() string {
	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:    "RSA PRIVATE KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PrivateKey(pk.privateKey),
		}))
}

// SetPemPKCS1 loads a private key from a PEM-encoded PKCS#1 string.
//
// This method decodes a PEM-encoded PKCS#1 private key and sets it as the
// current private key. The input should have the "RSA PRIVATE KEY" PEM block type.
//
// Parameters:
//   - pemPKCS1: PEM-encoded PKCS#1 private key string
//
// Returns:
//   - error: Error if decoding fails or format is invalid
//
// Behavior:
//   - Expects "-----BEGIN RSA PRIVATE KEY-----" header
//   - Validates the RSA key structure
//   - Replaces the current key if successful
//
// Example:
//
//	pemData, _ := os.ReadFile("private_key.pem")
//	privateKey := &rsa.PrivateKey{}
//	err := privateKey.SetPemPKCS1(string(pemData))
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PrivateKey) SetPemPKCS1(pemPKCS1 string) error {
	block, _ := pem.Decode([]byte(pemPKCS1))

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		return err
	} else {
		pk.Set(key)
		return nil
	}
}

// GetPemPKCS8 encodes the private key as a PEM-encoded PKCS#8 string.
//
// PKCS#8 is a more flexible, algorithm-agnostic private key format that can
// represent keys from various algorithms (RSA, ECDSA, Ed25519, etc.).
// This is the preferred format for modern applications.
//
// Returns:
//   - string: PEM-encoded PKCS#8 private key
//   - error: Error if encoding fails
//
// Behavior:
//   - Output format: "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"
//   - Algorithm-agnostic format (stores algorithm identifier)
//   - Better interoperability with other systems
//
// Example:
//
//	pemString, err := privateKey.GetPemPKCS8()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("private_key_pkcs8.pem", []byte(pemString), 0600)
func (pk *PrivateKey) GetPemPKCS8() (string, error) {
	if bytes, err := x509.MarshalPKCS8PrivateKey(pk.privateKey); err != nil {
		return "", err
	} else {
		return string(pem.EncodeToMemory(
			&pem.Block{
				Type:    "PRIVATE KEY",
				Headers: nil,
				Bytes:   bytes,
			})), nil
	}
}

// SetPemPKCS8 loads a private key from a PEM-encoded PKCS#8 string.
//
// This method decodes a PEM-encoded PKCS#8 private key and sets it as the
// current private key. PKCS#8 is the algorithm-agnostic format that can
// represent various key types.
//
// Parameters:
//   - pemPKCS8: PEM-encoded PKCS#8 private key string
//
// Returns:
//   - error: Error if decoding fails or key is not RSA
//
// Behavior:
//   - Expects "-----BEGIN PRIVATE KEY-----" header
//   - Automatically identifies RSA algorithm from key data
//   - Type assertion ensures the key is RSA
//
// Example:
//
//	pemData, _ := os.ReadFile("private_key_pkcs8.pem")
//	privateKey := &rsa.PrivateKey{}
//	err := privateKey.SetPemPKCS8(string(pemData))
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PrivateKey) SetPemPKCS8(pemPKCS8 string) error {
	block, _ := pem.Decode([]byte(pemPKCS8))

	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
		return err
	} else {
		pk.Set(key.(*rsa.PrivateKey))
		return nil
	}

}

// GetPublicKey extracts the public key from the private key.
//
// Every RSA private key contains the corresponding public key. This method
// creates a PublicKey instance containing the public key portion.
//
// Returns:
//   - PublicKey: The corresponding public key
//
// Behavior:
//   - The public key is derived from the private key's N and E values
//   - The returned public key can be safely shared with others
//   - Useful for creating matched key pairs for encryption/decryption
//
// Example:
//
//	privateKey := &rsa.PrivateKey{}
//	privateKey.SetBits(2048)
//	publicKey := privateKey.GetPublicKey()
//	// Share publicKey with others for encryption
func (pk *PrivateKey) GetPublicKey() PublicKey {
	publicKey := PublicKey{}

	publicKey.Set(pk.privateKey.PublicKey)

	return publicKey
}
