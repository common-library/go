// Package ecdsa provides ecdsa crypto related implementations.
package ecdsa

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// PublicKey is struct that provides public key related methods.
type PublicKey struct {
	publicKey ecdsa.PublicKey
}

// Verify verifies a digital signature against the message.
//
// This method verifies that a signature was created by the corresponding private key.
// The message is hashed with SHA-256 before verification, matching the Sign operation.
//
// Parameters:
//   - message: The original message that was signed
//   - signature: The signature to verify (containing R and S components)
//
// Returns:
//   - bool: true if signature is valid, false otherwise
//
// Behavior:
//   - Message is hashed with SHA-256 automatically
//   - Constant-time comparison (timing-safe)
//   - Returns false on any error (no error return)
//
// Example:
//
//	valid := publicKey.Verify("message to sign", signature)
//	if valid {
//	    fmt.Println("✓ Signature is valid")
//	} else {
//	    fmt.Println("✗ Signature is invalid")
//	}
func (pk *PublicKey) Verify(message string, signature Signature) bool {
	hash := sha256.Sum256([]byte(message))

	return ecdsa.Verify(&pk.publicKey, hash[:], signature.R, signature.S)
}

// Get returns the underlying ecdsa.PublicKey.
//
// This method provides direct access to the Go standard library ecdsa.PublicKey
// for use with other crypto functions that require the native type.
//
// Returns:
//   - ecdsa.PublicKey: The underlying public key (returned by value)
//
// Example:
//
//	nativeKey := publicKey.Get()
//	// Use with Go crypto functions
//	x509.MarshalPKIXPublicKey(&nativeKey)
func (pk *PublicKey) Get() ecdsa.PublicKey {
	return pk.publicKey
}

// Set assigns an existing ecdsa.PublicKey to this PublicKey instance.
//
// This method allows setting the public key from an existing Go standard library
// ecdsa.PublicKey, useful when loading keys from other sources.
//
// Parameters:
//   - publicKey: The ecdsa.PublicKey to assign
//
// Example:
//
//	nativeKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
//	publicKey := &ecdsa.PublicKey{}
//	publicKey.Set(nativeKey)
func (pk *PublicKey) Set(publicKey ecdsa.PublicKey) {
	pk.publicKey = publicKey
}

// GetPemPKIX encodes the public key as a PEM-encoded PKIX string.
//
// PKIX (Public-Key Infrastructure X.509) is the algorithm-agnostic public key
// format. This is the preferred format for modern applications as it can
// represent public keys from various algorithms (RSA, ECDSA, Ed25519, etc.).
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
//	os.WriteFile("ecdsa_public.pem", []byte(pemString), 0644)
func (pk *PublicKey) GetPemPKIX() (string, error) {
	if blockBytes, err := x509.MarshalPKIXPublicKey(&pk.publicKey); err != nil {
		return "", err
	} else {
		return string(pem.EncodeToMemory(
			&pem.Block{
				Type:    "PUBLIC KEY",
				Headers: nil,
				Bytes:   blockBytes,
			})), nil
	}
}

// SetPemPKIX loads a public key from a PEM-encoded PKIX string.
//
// This method decodes a PEM-encoded PKIX public key and sets it as the current
// public key. PKIX format is algorithm-agnostic and widely supported.
//
// Parameters:
//   - pemPKIX: PEM-encoded PKIX public key string
//
// Returns:
//   - error: Error if decoding fails or key is not ECDSA
//
// Behavior:
//   - Expects "-----BEGIN PUBLIC KEY-----" header
//   - Automatically identifies ECDSA algorithm from key data
//   - Type assertion ensures the key is ECDSA
//
// Example:
//
//	pemData, _ := os.ReadFile("ecdsa_public.pem")
//	publicKey := &ecdsa.PublicKey{}
//	err := publicKey.SetPemPKIX(string(pemData))
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PublicKey) SetPemPKIX(pemPKIX string) error {
	block, _ := pem.Decode([]byte(pemPKIX))

	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return err
	} else {
		pk.publicKey = *key.(*ecdsa.PublicKey)
		return nil
	}
}

// GetSsh encodes the public key as an SSH authorized_keys format string.
//
// This method converts the ECDSA public key to the SSH authorized_keys format,
// which is used for SSH authentication and other SSH-based systems.
//
// Returns:
//   - string: SSH authorized_keys format string
//   - error: Error if encoding fails
//
// Behavior:
//   - Output format: "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTY..." (single line)
//   - Curve name included in format (nistp256, nistp384, nistp521)
//   - Compatible with OpenSSH and other SSH implementations
//
// Example:
//
//	sshKey, err := publicKey.GetSsh()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("id_ecdsa.pub", []byte(sshKey), 0644)
func (pk *PublicKey) GetSsh() (string, error) {
	if publicKey, err := ssh.NewPublicKey(&pk.publicKey); err != nil {
		return "", err
	} else {
		return string(ssh.MarshalAuthorizedKey(publicKey)), nil
	}
}

// SetSsh loads a public key from an SSH authorized_keys format string.
//
// This method parses an SSH authorized_keys format string and sets the ECDSA
// public key. The format is commonly found in ~/.ssh/authorized_keys files.
//
// Parameters:
//   - sshKey: SSH authorized_keys format string (e.g., "ecdsa-sha2-nistp256 AAAAE...")
//
// Returns:
//   - error: Error if parsing fails or key is not ECDSA
//
// Behavior:
//   - Accepts single-line SSH format
//   - Ignores comments and options if present
//   - Type assertion ensures the key is ECDSA
//
// Example:
//
//	sshData, _ := os.ReadFile("id_ecdsa.pub")
//	publicKey := &ecdsa.PublicKey{}
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
// This method creates an ssh.PublicKey instance from the ECDSA public key,
// which can be used with the golang.org/x/crypto/ssh package for SSH operations.
//
// Returns:
//   - ssh.PublicKey: The SSH public key instance
//   - error: Error if conversion fails
//
// Behavior:
//   - Returns a type that implements ssh.PublicKey interface
//   - Can be used with ssh package functions
//   - Preserves all key material including curve
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
// This method converts an ssh.PublicKey instance to an ECDSA public key.
// The ssh.PublicKey must contain an ECDSA key, otherwise an error occurs.
//
// Parameters:
//   - publicKey: The ssh.PublicKey to convert
//
// Returns:
//   - error: Error if conversion fails or key is not ECDSA
//
// Behavior:
//   - Type assertion ensures the key is ECDSA
//   - Marshals and unmarshals to ensure proper conversion
//   - Replaces the current key if successful
//
// Example:
//
//	sshPubKey, _ := ssh.ParsePublicKey(sshKeyBytes)
//	publicKey := &ecdsa.PublicKey{}
//	err := publicKey.SetSshPublicKey(sshPubKey)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (pk *PublicKey) SetSshPublicKey(publicKey ssh.PublicKey) error {
	if key, err := ssh.ParsePublicKey(publicKey.Marshal()); err != nil {
		return err
	} else {
		pk.publicKey = *key.(ssh.CryptoPublicKey).CryptoPublicKey().(*ecdsa.PublicKey)
		return nil
	}
}
