// Package rsa provides rsa crypto related implementations.
package rsa

// KeyPair is struct that provides key pair related methods.
type KeyPair struct {
	privateKey PrivateKey
	publicKey  PublicKey
}

// Generate is to generate a key pair.
//
// ex) err := keyPair.Generate(4096)
func (this *KeyPair) Generate(bits int) error {
	if err := this.privateKey.SetBits(bits); err != nil {
		return err
	} else {
		this.publicKey = this.privateKey.GetPublicKey()
		return nil
	}
}

// EncryptPKCS1v15 is encrypt plaintext.
//
// ex) ciphertext, err := keyPair.EncryptPKCS1v15(plaintext)
func (this *KeyPair) EncryptPKCS1v15(plaintext string) ([]byte, error) {
	return this.privateKey.EncryptPKCS1v15(plaintext)
}

// DecryptPKCS1v15 is decrypt ciphertext.
//
// ex) plaintext, err := keyPair.DecryptPKCS1v15(ciphertext)
func (this *KeyPair) DecryptPKCS1v15(ciphertext []byte) (string, error) {
	return this.privateKey.DecryptPKCS1v15(ciphertext)
}

// EncryptOAEP is encrypt plaintext.
//
// ex) ciphertext, err := keyPair.EncryptOAEP(plaintext)
func (this *KeyPair) EncryptOAEP(plaintext string) ([]byte, error) {
	return this.privateKey.EncryptOAEP(plaintext)
}

// DecryptOAEP is decrypt ciphertext.
//
// ex) plaintext, err := keyPair.DecryptOAEP(ciphertext)
func (this *KeyPair) DecryptOAEP(ciphertext []byte) (string, error) {
	return this.privateKey.DecryptOAEP(ciphertext)
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
