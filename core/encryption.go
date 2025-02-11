package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"errors"

	"github.com/hedgehog125/project-reboot/util"
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
	nonce := util.CryptoRandomBytes(gcm.NonceSize())

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

func CheckPassword(givenPassword string, passwordHash []byte, passwordSalt []byte, hashSettings *HashSettings) bool {
	givenPasswordHash := hashWithSalt(givenPassword, passwordSalt, hashSettings)
	return subtle.ConstantTimeCompare(givenPasswordHash, passwordHash) == 1
}

var ErrIncorrectPassword error = errors.New("incorrect password")

func Decrypt(password string, encryptedData *EncryptedData) ([]byte, error) {
	if !CheckPassword(password, encryptedData.PasswordHash, encryptedData.PasswordSalt, &encryptedData.HashSettings) {
		return nil, ErrIncorrectPassword
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
	salt = util.CryptoRandomBytes(SALT_LENGTH)
	hash = hashWithSalt(password, salt, settings)
	return
}
func hashWithSalt(password string, salt []byte, settings *HashSettings) []byte {
	return argon2.IDKey([]byte(password), salt, settings.Time, settings.Memory, HASH_THREADS, settings.KeyLen)
}
