// Package rsa provides rsa crypto related implementations.
package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// PublicKey is struct that provides public key related methods.
type PublicKey struct {
	publicKey rsa.PublicKey
}

// EncryptPKCS1v15 encrypts plaintext using PKCS#1 v1.5 padding.
//
// This method encrypts data with the public key using PKCS#1 v1.5 padding.
// This padding scheme is legacy and has known vulnerabilities. Use EncryptOAEP
// for new applications.
//
// Parameters:
//   - plaintext: The string to encrypt
//
// Returns:
//   - []byte: The encrypted ciphertext
//   - error: Error if encryption fails or plaintext is too long
//
// Behavior:
//   - Maximum plaintext size depends on key size (e.g., ~245 bytes for 2048-bit key)
//   - Same plaintext produces different ciphertext each time
//   - Only the corresponding private key can decrypt
//
// Example:
//
//	publicKey := &rsa.PublicKey{}
//	// ... load public key ...
//	ciphertext, err := publicKey.EncryptPKCS1v15("secret message")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PublicKey) EncryptPKCS1v15(plaintext string) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, &pk.publicKey, []byte(plaintext))
}

// EncryptOAEP encrypts plaintext using OAEP padding (recommended).
//
// This method uses Optimal Asymmetric Encryption Padding (OAEP) with SHA-256,
// which provides better security than PKCS#1 v1.5. This is the recommended
// padding scheme for RSA encryption.
//
// Parameters:
//   - plaintext: The string to encrypt
//
// Returns:
//   - []byte: The encrypted ciphertext
//   - error: Error if encryption fails or plaintext is too long
//
// Behavior:
//   - Maximum plaintext size = key_size - 2*hash_size - 2 (e.g., 190 bytes for 2048-bit key)
//   - Provides provable security against chosen-ciphertext attacks
//   - Non-deterministic encryption
//
// Example:
//
//	publicKey := &rsa.PublicKey{}
//	// ... load public key ...
//	ciphertext, err := publicKey.EncryptOAEP("confidential data")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PublicKey) EncryptOAEP(plaintext string) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, &pk.publicKey, []byte(plaintext), nil)
}

// Get returns the underlying rsa.PublicKey.
//
// This method provides direct access to the Go standard library rsa.PublicKey
// for use with other crypto functions that require the native type.
//
// Returns:
//   - rsa.PublicKey: The underlying public key (returned by value)
//
// Example:
//
//	nativeKey := publicKey.Get()
//	// Use with Go crypto functions
//	x509.MarshalPKIXPublicKey(&nativeKey)
func (pk *PublicKey) Get() rsa.PublicKey {
	return pk.publicKey
}

// Set assigns an existing rsa.PublicKey to this PublicKey instance.
//
// This method allows setting the public key from an existing Go standard library
// rsa.PublicKey, useful when loading keys from other sources.
//
// Parameters:
//   - publicKey: The rsa.PublicKey to assign
//
// Example:
//
//	nativeKey := rsa.PublicKey{N: n, E: e}
//	publicKey := &rsa.PublicKey{}
//	publicKey.Set(nativeKey)
func (pk *PublicKey) Set(publicKey rsa.PublicKey) {
	pk.publicKey = publicKey
}

// GetPemPKCS1 encodes the public key as a PEM-encoded PKCS#1 string.
//
// PKCS#1 is the traditional RSA public key format. The PEM encoding allows
// the key to be stored in text files. The format uses "RSA PUBLIC KEY" as
// the PEM block type.
//
// Returns:
//   - string: PEM-encoded PKCS#1 public key
//
// Behavior:
//   - Output format: "-----BEGIN RSA PUBLIC KEY-----\n...\n-----END RSA PUBLIC KEY-----"
//   - This format is RSA-specific
//   - Suitable for sharing with systems expecting PKCS#1 format
//
// Example:
//
//	pemString := publicKey.GetPemPKCS1()
//	os.WriteFile("public_key.pem", []byte(pemString), 0644)
func (pk *PublicKey) GetPemPKCS1() string {
	return string(pem.EncodeToMemory(
		&pem.Block{
			Type:    "RSA PUBLIC KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PublicKey(&pk.publicKey),
		}))
}

// SetPemPKCS1 loads a public key from a PEM-encoded PKCS#1 string.
//
// This method decodes a PEM-encoded PKCS#1 public key and sets it as the
// current public key. The input should have the "RSA PUBLIC KEY" PEM block type.
//
// Parameters:
//   - pemPKCS1: PEM-encoded PKCS#1 public key string
//
// Returns:
//   - error: Error if decoding fails or format is invalid
//
// Behavior:
//   - Expects "-----BEGIN RSA PUBLIC KEY-----" header
//   - Validates the RSA public key structure
//   - Replaces the current key if successful
//
// Example:
//
//	pemData, _ := os.ReadFile("public_key.pem")
//	publicKey := &rsa.PublicKey{}
//	err := publicKey.SetPemPKCS1(string(pemData))
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PublicKey) SetPemPKCS1(pemPKCS1 string) error {
	block, _ := pem.Decode([]byte(pemPKCS1))

	if key, err := x509.ParsePKCS1PublicKey(block.Bytes); err != nil {
		return err
	} else {
		pk.Set(*key)
		return nil
	}
}

// GetPemPKIX encodes the public key as a PEM-encoded PKIX string.
//
// PKIX (Public-Key Infrastructure X.509) is the algorithm-agnostic public key
// format defined in X.509. This is the preferred format for modern applications
// as it can represent public keys from various algorithms.
//
// Returns:
//   - string: PEM-encoded PKIX public key
//   - error: Error if encoding fails
//
// Behavior:
//   - Output format: "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----"
//   - Algorithm-agnostic format (includes algorithm identifier)
//   - Better interoperability with other systems
//
// Example:
//
//	pemString, err := publicKey.GetPemPKIX()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("public_key_pkix.pem", []byte(pemString), 0644)
func (pk *PublicKey) GetPemPKIX() (string, error) {
	if bytes, err := x509.MarshalPKIXPublicKey(&pk.publicKey); err != nil {
		return "", err
	} else {
		return string(pem.EncodeToMemory(
			&pem.Block{
				Type:    "PUBLIC KEY",
				Headers: nil,
				Bytes:   bytes,
			})), nil
	}
}

// SetPemPKIX loads a public key from a PEM-encoded PKIX string.
//
// This method decodes a PEM-encoded PKIX public key and sets it as the
// current public key. PKIX format is algorithm-agnostic and widely supported.
//
// Parameters:
//   - pemPKIX: PEM-encoded PKIX public key string
//
// Returns:
//   - error: Error if decoding fails or key is not RSA
//
// Behavior:
//   - Expects "-----BEGIN PUBLIC KEY-----" header
//   - Automatically identifies RSA algorithm from key data
//   - Type assertion ensures the key is RSA
//
// Example:
//
//	pemData, _ := os.ReadFile("public_key_pkix.pem")
//	publicKey := &rsa.PublicKey{}
//	err := publicKey.SetPemPKIX(string(pemData))
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PublicKey) SetPemPKIX(pemPKIX string) error {
	block, _ := pem.Decode([]byte(pemPKIX))

	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return err
	} else {
		pk.Set(*key.(*rsa.PublicKey))
		return nil
	}
}

// GetSsh encodes the public key as an SSH authorized_keys format string.
//
// This method converts the RSA public key to the SSH authorized_keys format,
// which is used for SSH authentication and other SSH-based systems.
//
// Returns:
//   - string: SSH authorized_keys format string
//   - error: Error if encoding fails
//
// Behavior:
//   - Output format: "ssh-rsa AAAAB3NzaC1yc2EA..." (single line)
//   - Suitable for ~/.ssh/authorized_keys files
//   - Compatible with OpenSSH and other SSH implementations
//
// Example:
//
//	sshKey, err := publicKey.GetSsh()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("id_rsa.pub", []byte(sshKey), 0644)
func (pk *PublicKey) GetSsh() (string, error) {
	if publicKey, err := ssh.NewPublicKey(&pk.publicKey); err != nil {
		return "", err
	} else {
		return string(ssh.MarshalAuthorizedKey(publicKey)), nil
	}
}

// SetSsh loads a public key from an SSH authorized_keys format string.
//
// This method parses an SSH authorized_keys format string and sets the RSA
// public key. The format is commonly found in ~/.ssh/authorized_keys files.
//
// Parameters:
//   - sshKey: SSH authorized_keys format string (e.g., "ssh-rsa AAAAB3...")
//
// Returns:
//   - error: Error if parsing fails or key is not RSA
//
// Behavior:
//   - Accepts single-line SSH format
//   - Ignores comments and options if present
//   - Type assertion ensures the key is RSA
//
// Example:
//
//	sshData, _ := os.ReadFile("id_rsa.pub")
//	publicKey := &rsa.PublicKey{}
//	err := publicKey.SetSsh(string(sshData))
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PublicKey) SetSsh(sshKey string) error {
	if key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshKey)); err != nil {
		return err
	} else {
		return pk.SetSshPublicKey(key)
	}
}

// GetSshPublicKey converts the public key to an ssh.PublicKey.
//
// This method creates an ssh.PublicKey instance from the RSA public key,
// which can be used with the golang.org/x/crypto/ssh package for SSH operations.
//
// Returns:
//   - ssh.PublicKey: The SSH public key instance
//   - error: Error if conversion fails
//
// Behavior:
//   - Returns a type that implements ssh.PublicKey interface
//   - Can be used with ssh package functions
//   - Preserves all key material
//
// Example:
//
//	sshPubKey, err := publicKey.GetSshPublicKey()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use with SSH package
//	cert := &ssh.Certificate{Key: sshPubKey}
func (pk *PublicKey) GetSshPublicKey() (ssh.PublicKey, error) {
	return ssh.NewPublicKey(&pk.publicKey)
}

// SetSshPublicKey loads a public key from an ssh.PublicKey.
//
// This method converts an ssh.PublicKey instance to an RSA public key.
// The ssh.PublicKey must contain an RSA key, otherwise an error occurs.
//
// Parameters:
//   - publicKey: The ssh.PublicKey to convert
//
// Returns:
//   - error: Error if conversion fails or key is not RSA
//
// Behavior:
//   - Type assertion ensures the key is RSA
//   - Marshals and unmarshals to ensure proper conversion
//   - Replaces the current key if successful
//
// Example:
//
//	sshPubKey, _ := ssh.ParsePublicKey(sshKeyBytes)
//	publicKey := &rsa.PublicKey{}
//	err := publicKey.SetSshPublicKey(sshPubKey)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PublicKey) SetSshPublicKey(publicKey ssh.PublicKey) error {
	if key, err := ssh.ParsePublicKey(publicKey.Marshal()); err != nil {
		return err
	} else {
		pk.publicKey = *key.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
		return nil
	}
}
