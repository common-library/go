// Package dsa provides dsa crypto related implementations.
package dsa

import "crypto/dsa"

// KeyPair is struct that provides key pair related methods.
type KeyPair struct {
	privateKey PrivateKey
	publicKey  PublicKey
}

// Generate is to generate a key pair.
//
// ex) err := keyPair.Generate(crypto_dsa.L1024N160)
func (this *KeyPair) Generate(parameterSizes dsa.ParameterSizes) error {
	if err := this.privateKey.SetSizes(parameterSizes); err != nil {
		return err
	} else {
		this.publicKey = this.privateKey.GetPublicKey()
		return nil
	}
}

// Sign is create a signature for message.
//
// ex) signature, err := keyPair.Sign(message)
func (this *KeyPair) Sign(message string) (Signature, error) {
	return this.privateKey.Sign(message)
}

// Verify is verifies the signature.
//
// ex) result := keyPair.Verify(message, signature)
func (this *KeyPair) Verify(message string, signature Signature) bool {
	return this.publicKey.Verify(message, signature)
}

// GetKeyPair is to get a key pair.
//
// ex) privateKey, publicKey := keyPair.GetKeyPair()
func (this *KeyPair) GetKeyPair() (privateKey PrivateKey, publicKey PublicKey) {
	return this.privateKey, this.publicKey
}

// SetKeyPair is to set a key pair.
//
// ex) keyPair.SetKeyPair(privateKey, publicKey)
func (this *KeyPair) SetKeyPair(privateKey PrivateKey, publicKey PublicKey) {
	this.privateKey = privateKey
	this.publicKey = publicKey
}
