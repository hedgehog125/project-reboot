package core

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/hedgehog125/project-reboot/common"
)

// Adapted from: https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/
func Encrypt(data []byte, encryptionKey []byte) ([]byte, []byte, *common.Error) {
	passwordCipher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, nil, ErrWrapperEncrypt.Wrap(err)
	}
	gcm, err := cipher.NewGCM(passwordCipher)
	if err != nil {
		return nil, nil, ErrWrapperEncrypt.Wrap(err)
	}
	nonce := common.CryptoRandomBytes(gcm.NonceSize())

	encrypted := gcm.Seal(nil, nonce, data, nil)
	return encrypted, nonce, nil
}

func Decrypt(encrypted []byte, encryptionKey []byte, nonce []byte) ([]byte, *common.Error) {
	passwordCipher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, ErrWrapperDecrypt.Wrap(err)
	}

	gcm, err := cipher.NewGCM(passwordCipher)
	if err != nil {
		return nil, ErrWrapperDecrypt.Wrap(err)
	}

	decrypted, err := gcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, ErrWrapperDecrypt.Wrap(err)
	}
	return decrypted, nil
}
