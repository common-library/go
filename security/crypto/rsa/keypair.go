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
func (kp *KeyPair) Generate(bits int) error {
	if err := kp.privateKey.SetBits(bits); err != nil {
		return err
	} else {
		kp.publicKey = kp.privateKey.GetPublicKey()
		return nil
	}
}

// EncryptPKCS1v15 is encrypt plaintext.
//
// ex) ciphertext, err := keyPair.EncryptPKCS1v15(plaintext)
func (kp *KeyPair) EncryptPKCS1v15(plaintext string) ([]byte, error) {
	return kp.privateKey.EncryptPKCS1v15(plaintext)
}

// DecryptPKCS1v15 is decrypt ciphertext.
//
// ex) plaintext, err := keyPair.DecryptPKCS1v15(ciphertext)
func (kp *KeyPair) DecryptPKCS1v15(ciphertext []byte) (string, error) {
	return kp.privateKey.DecryptPKCS1v15(ciphertext)
}

// EncryptOAEP is encrypt plaintext.
//
// ex) ciphertext, err := keyPair.EncryptOAEP(plaintext)
func (kp *KeyPair) EncryptOAEP(plaintext string) ([]byte, error) {
	return kp.privateKey.EncryptOAEP(plaintext)
}

// DecryptOAEP is decrypt ciphertext.
//
// ex) plaintext, err := keyPair.DecryptOAEP(ciphertext)
func (kp *KeyPair) DecryptOAEP(ciphertext []byte) (string, error) {
	return kp.privateKey.DecryptOAEP(ciphertext)
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
