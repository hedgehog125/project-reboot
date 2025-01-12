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
	HashSettings HashSettings
}

const HASH_THREADS = 2
const SALT_LENGTH = 128

type HashSettings struct {
	Time   uint32
	Memory uint32
	KeyLen uint32
}

// Adapted from: https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/
func Encrypt(data []byte, password string) (*EncryptedData, error) {
	hashSettings := HashSettings{
		Time:   5,
		Memory: 128 * 1024,
		KeyLen: 32,
	}
	passwordHash, passwordSalt := hash(password, &hashSettings)
	encryptionKey, encryptionKeySalt := hash(password, &hashSettings)

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
		HashSettings: hashSettings,
	}, nil
}

var ErrIncorrectPassword error = errors.New("incorrect password")

func Decrypt(password string, encryptedData *EncryptedData) ([]byte, error) {
	{
		passwordHash := hashWithSalt(password, encryptedData.PasswordSalt, &encryptedData.HashSettings)
		if subtle.ConstantTimeCompare(passwordHash, encryptedData.PasswordHash) == 0 {
			return nil, ErrIncorrectPassword
		}
	}
	encryptionKey := hashWithSalt(password, encryptedData.KeySalt, &encryptedData.HashSettings)

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
func hash(password string, settings *HashSettings) (hash, salt []byte) {
	salt = randomBytes(SALT_LENGTH)
	hash = hashWithSalt(password, salt, settings)
	return
}
func hashWithSalt(password string, salt []byte, settings *HashSettings) []byte {
	return argon2.IDKey([]byte(password), salt, settings.Time, settings.Memory, HASH_THREADS, settings.KeyLen)
}

func randomBytes(length int) []byte {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatalf("randomBytes: couldn't get random byte. error:\n%v", err)
	}
	return salt
}
