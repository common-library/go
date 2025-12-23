// Package rsa provides RSA public-key cryptography.
//
// This package implements RSA encryption and digital signatures with support
// for multiple key sizes, PKCS1v15 and OAEP padding schemes, and PEM format
// key storage.
//
// # Features
//
//   - RSA key pair generation (2048, 3072, 4096 bits)
//   - PKCS1v15 and OAEP encryption/decryption
//   - Digital signatures with PSS padding
//   - PEM/PKCS1/PKCS8 format support
//   - SSH public key format conversion
//
// # Basic Example
//
//	keyPair := &rsa.KeyPair{}
//	err := keyPair.Generate(2048)
//	ciphertext, _ := keyPair.EncryptOAEP("secret")
//	plaintext, _ := keyPair.DecryptOAEP(ciphertext)
package rsa

// KeyPair provides RSA key pair operations.
type KeyPair struct {
	privateKey PrivateKey
	publicKey  PublicKey
}

// Generate creates a new RSA key pair.
//
// # Parameters
//
//   - bits: Key size in bits (2048, 3072, or 4096 recommended)
//
// # Returns
//
//   - error: Error if generation fails, nil on success
//
// # Examples
//
//	keyPair := &rsa.KeyPair{}
//	err := keyPair.Generate(2048)
func (kp *KeyPair) Generate(bits int) error {
	if err := kp.privateKey.SetBits(bits); err != nil {
		return err
	} else {
		kp.publicKey = kp.privateKey.GetPublicKey()
		return nil
	}
}

// EncryptPKCS1v15 encrypts plaintext using PKCS#1 v1.5 padding.
//
// # Parameters
//
//   - plaintext: Text to encrypt
//
// # Returns
//
//   - []byte: Encrypted ciphertext
//   - error: Error if encryption fails, nil on success
//
// # Examples
//
//	ciphertext, err := keyPair.EncryptPKCS1v15("secret message")
func (kp *KeyPair) EncryptPKCS1v15(plaintext string) ([]byte, error) {
	return kp.privateKey.EncryptPKCS1v15(plaintext)
}

// DecryptPKCS1v15 decrypts ciphertext using PKCS#1 v1.5 padding.
//
// # Parameters
//
//   - ciphertext: Encrypted data
//
// # Returns
//
//   - string: Decrypted plaintext
//   - error: Error if decryption fails, nil on success
//
// # Examples
//
//	plaintext, err := keyPair.DecryptPKCS1v15(ciphertext)
func (kp *KeyPair) DecryptPKCS1v15(ciphertext []byte) (string, error) {
	return kp.privateKey.DecryptPKCS1v15(ciphertext)
}

// EncryptOAEP encrypts plaintext using OAEP padding (recommended).
//
// # Parameters
//
//   - plaintext: Text to encrypt
//
// # Returns
//
//   - []byte: Encrypted ciphertext
//   - error: Error if encryption fails, nil on success
//
// # Examples
//
//	ciphertext, err := keyPair.EncryptOAEP("secret message")
func (kp *KeyPair) EncryptOAEP(plaintext string) ([]byte, error) {
	return kp.privateKey.EncryptOAEP(plaintext)
}

// DecryptOAEP decrypts ciphertext using OAEP padding (recommended).
//
// # Parameters
//
//   - ciphertext: Encrypted data
//
// # Returns
//
//   - string: Decrypted plaintext
//   - error: Error if decryption fails, nil on success
//
// # Examples
//
//	plaintext, err := keyPair.DecryptOAEP(ciphertext)
func (kp *KeyPair) DecryptOAEP(ciphertext []byte) (string, error) {
	return kp.privateKey.DecryptOAEP(ciphertext)
}

// GetKeyPair retrieves the private and public keys.
//
// # Returns
//
//   - privateKey: RSA private key
//   - publicKey: RSA public key
//
// # Examples
//
//	privateKey, publicKey := keyPair.GetKeyPair()
func (kp *KeyPair) GetKeyPair() (privateKey PrivateKey, publicKey PublicKey) {
	return kp.privateKey, kp.publicKey
}

// SetKeyPair sets the private and public keys.
//
// # Parameters
//
//   - privateKey: RSA private key
//   - publicKey: RSA public key
//
// # Examples
//
//	keyPair.SetKeyPair(privateKey, publicKey)
func (kp *KeyPair) SetKeyPair(privateKey PrivateKey, publicKey PublicKey) {
	kp.privateKey = privateKey
	kp.publicKey = publicKey
}
