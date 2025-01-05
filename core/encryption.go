package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"log"

	"golang.org/x/crypto/argon2"
)

type EncryptedData struct {
	Data         []byte
	Nonce        []byte
	KeySalt      []byte
	PasswordHash []byte
	PasswordSalt []byte
}

// Adapted from: https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/
func Encrypt(data []byte, password string) (*EncryptedData, error) {
	passwordHash, passwordSalt := hash(password)
	encryptionKey, encryptionKeySalt := hash(password)

	passwordCipher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(passwordCipher)
	if err != nil {
		return nil, err
	}
	nonce := randomBytes(gcm.NonceSize()) // TODO: how bad are collisions?

	encrypted := gcm.Seal(nil, nonce, data, nil)
	return &EncryptedData{
		Data:         encrypted,
		Nonce:        nonce,
		KeySalt:      encryptionKeySalt,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
	}, nil
}

var ErrIncorrectPassword error = errors.New("incorrect password")

func Decrypt(password string, encryptedData *EncryptedData) ([]byte, error) {
	{
		passwordHash := hashWithSalt(password, encryptedData.PasswordSalt)
		if subtle.ConstantTimeCompare(passwordHash, encryptedData.PasswordHash) == 0 {
			return nil, ErrIncorrectPassword
		}
	}
	encryptionKey := hashWithSalt(password, encryptedData.KeySalt)

	passwordCipher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(passwordCipher)
	if err != nil {
		return nil, err
	}

	decrypted, err := gcm.Open(nil, encryptedData.Nonce, encryptedData.Data, nil)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
func hash(password string) (hash, salt []byte) {
	salt = randomBytes(128)
	hash = hashWithSalt(password, salt)
	return
}
func hashWithSalt(password string, salt []byte) []byte {
	// TODO: save the constants in the db for each hash
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 2, 32)
}

func randomBytes(length int) []byte {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatalf("randomBytes: couldn't get random byte. error:\n%v", err)
	}
	return salt
}
