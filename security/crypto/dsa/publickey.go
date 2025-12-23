// Package dsa provides dsa crypto related implementations.
//
// Deprecated: DSA is a legacy algorithm and should not be used for new applications.
// Use Ed25519 (crypto/ed25519) or other modern alternatives instead.
// DSA keys with 1024-bit moduli are cryptographically weak.
package dsa

import (
	//lint:ignore SA1019 DSA is deprecated but kept for compatibility
	"crypto/dsa"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

// PublicKey is struct that provides public key related methods.
type PublicKey struct {
	publicKey dsa.PublicKey
}

// Verify verifies a digital signature against the message.
//
// Deprecated: DSA is cryptographically weak and should not be used.
// Use Ed25519 or ECDSA for new applications.
//
// This method verifies that a signature was created by the corresponding private key.
// The message is hashed with SHA-256 before verification.
//
// Parameters:
//   - message: The original message that was signed
//   - signature: The signature to verify
//
// Returns:
//   - bool: true if signature is valid, false otherwise
//
// Example:
//
//	// For legacy compatibility only
//	valid := publicKey.Verify("legacy message", signature)
func (pub *PublicKey) Verify(message string, signature Signature) bool {
	hash := sha256.Sum256([]byte(message))

	return dsa.Verify(&pub.publicKey, hash[:], signature.R, signature.S)
}

// Get returns the underlying dsa.PublicKey.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// This method provides direct access to the Go standard library dsa.PublicKey.
//
// Returns:
//   - dsa.PublicKey: The underlying public key
func (pub *PublicKey) Get() dsa.PublicKey {
	return pub.publicKey
}

// Set assigns an existing dsa.PublicKey to this PublicKey instance.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// Parameters:
//   - publicKey: The dsa.PublicKey to assign
func (pub *PublicKey) Set(publicKey dsa.PublicKey) {
	pub.publicKey = publicKey
}

// GetPemAsn1 encodes the public key as a PEM-encoded ASN.1 string.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// This method converts the DSA public key to PEM format using ASN.1 encoding.
//
// Returns:
//   - string: PEM-encoded DSA public key
//   - error: Error if encoding fails
//
// Example:
//
//	// For legacy compatibility only
//	pemString, err := publicKey.GetPemAsn1()
func (pub *PublicKey) GetPemAsn1() (string, error) {
	if blockBytes, err := asn1.Marshal(pub.publicKey); err != nil {
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

// SetPemAsn1 loads a public key from a PEM-encoded ASN.1 string.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// This method decodes a PEM-encoded DSA public key and sets it as the current key.
//
// Parameters:
//   - pemAsn1: PEM-encoded ASN.1 DSA public key string
//
// Returns:
//   - error: Error if decoding fails
//
// Example:
//
//	// For legacy compatibility only
//	pemData, _ := os.ReadFile("dsa_public.pem")
//	publicKey := &dsa.PublicKey{}
//	err := publicKey.SetPemAsn1(string(pemData))
func (pub *PublicKey) SetPemAsn1(pemAsn1 string) error {
	block, _ := pem.Decode([]byte(pemAsn1))
	if _, err := asn1.Unmarshal(block.Bytes, &pub.publicKey); err != nil {
		return err
	} else {
		return nil
	}

}

// GetSsh encodes the public key as an SSH authorized_keys format string.
//
// Deprecated: DSA is cryptographically weak and should not be used.
// Modern SSH implementations are removing DSA support.
//
// This method converts the DSA public key to SSH format.
//
// Returns:
//   - string: SSH authorized_keys format string
//   - error: Error if encoding fails
//
// Example:
//
//	// For legacy compatibility only
//	sshKey, err := publicKey.GetSsh()
func (pub *PublicKey) GetSsh() (string, error) {
	if publicKey, err := ssh.NewPublicKey(&pub.publicKey); err != nil {
		return "", err
	} else {
		return string(ssh.MarshalAuthorizedKey(publicKey)), nil
	}
}

// SetSsh loads a public key from an SSH authorized_keys format string.
//
// Deprecated: DSA is cryptographically weak and should not be used.
// Modern SSH implementations are removing DSA support.
//
// Parameters:
//   - sshKey: SSH authorized_keys format string
//
// Returns:
//   - error: Error if parsing fails
//
// Example:
//
//	// For legacy compatibility only
//	sshData, _ := os.ReadFile("id_dsa.pub")
//	publicKey := &dsa.PublicKey{}
//	err := publicKey.SetSsh(string(sshData))
func (pub *PublicKey) SetSsh(sshKey string) error {
	if key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshKey)); err != nil {
		return err
	} else {
		return pub.SetSshPublicKey(key)
	}
}

// GetSshPublicKey converts the public key to an ssh.PublicKey.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// Returns:
//   - ssh.PublicKey: The SSH public key instance
//   - error: Error if conversion fails
//
// Example:
//
//	// For legacy compatibility only
//	sshPubKey, err := publicKey.GetSshPublicKey()
func (pub *PublicKey) GetSshPublicKey() (ssh.PublicKey, error) {
	return ssh.NewPublicKey(&pub.publicKey)
}

// SetSshPublicKey loads a public key from an ssh.PublicKey.
//
// Deprecated: DSA is cryptographically weak and should not be used.
//
// Parameters:
//   - publicKey: The ssh.PublicKey to convert
//
// Returns:
//   - error: Error if conversion fails or key is not DSA
//
// Example:
//
//	// For legacy compatibility only
//	sshPubKey, _ := ssh.ParsePublicKey(sshKeyBytes)
//	publicKey := &dsa.PublicKey{}
//	err := publicKey.SetSshPublicKey(sshPubKey)
func (pub *PublicKey) SetSshPublicKey(publicKey ssh.PublicKey) error {
	if key, err := ssh.ParsePublicKey(publicKey.Marshal()); err != nil {
		return err
	} else {
		pub.publicKey = *key.(ssh.CryptoPublicKey).CryptoPublicKey().(*dsa.PublicKey)
		return nil
	}
}
