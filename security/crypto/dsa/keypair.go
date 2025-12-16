// Package dsa provides dsa crypto related implementations.
//
// Deprecated: DSA is a legacy algorithm and should not be used for new applications.
// Use Ed25519 (crypto/ed25519) or other modern alternatives instead.
// DSA keys with 1024-bit moduli are cryptographically weak.
package dsa

//lint:ignore SA1019 DSA is deprecated but kept for compatibility
import "crypto/dsa"

// KeyPair is struct that provides key pair related methods.
type KeyPair struct {
	privateKey PrivateKey
	publicKey  PublicKey
}

// Generate is to generate a key pair.
//
// ex) err := keyPair.Generate(crypto_dsa.L1024N160)
func (kp *KeyPair) Generate(parameterSizes dsa.ParameterSizes) error {
	if err := kp.privateKey.SetSizes(parameterSizes); err != nil {
		return err
	} else {
		kp.publicKey = kp.privateKey.GetPublicKey()
		return nil
	}
}

// Sign is create a signature for message.
//
// ex) signature, err := keyPair.Sign(message)
func (kp *KeyPair) Sign(message string) (Signature, error) {
	return kp.privateKey.Sign(message)
}

// Verify is verifies the signature.
//
// ex) result := keyPair.Verify(message, signature)
func (kp *KeyPair) Verify(message string, signature Signature) bool {
	return kp.publicKey.Verify(message, signature)
}

// GetKeyPair is to get a key pair.
//
// ex) privateKey, publicKey := keyPair.GetKeyPair()
func (kp *KeyPair) GetKeyPair() (privateKey PrivateKey, publicKey PublicKey) {
	return kp.privateKey, kp.publicKey
}

// SetKeyPair is to set a key pair.
//
// ex) keyPair.SetKeyPair(privateKey, publicKey)
func (kp *KeyPair) SetKeyPair(privateKey PrivateKey, publicKey PublicKey) {
	kp.privateKey = privateKey
	kp.publicKey = publicKey
}
