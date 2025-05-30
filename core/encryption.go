package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"

	"github.com/hedgehog125/project-reboot/common"
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

// Hash run settings
const (
	HashThreads = 2
	SaltLength  = 128
)

var defaultHashSettings = HashSettings{
	Time:   5,
	Memory: 128 * 1024,
	KeyLen: 32,

	// Minimum
	// Time:   1,
	// Memory: 1 * 1024,
	// KeyLen: 16,
}

type HashSettings struct {
	Time   uint32
	Memory uint32
	KeyLen uint32
}

// Adapted from: https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/
func Encrypt(data []byte, password string) (*EncryptedData, *common.Error) {
	passwordHash, passwordSalt := hash(password, defaultHashSettings)
	encryptionKey, encryptionKeySalt := hash(password, defaultHashSettings)

	passwordCipher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, ErrWrapperEncrypt.Wrap(err)
	}
	gcm, err := cipher.NewGCM(passwordCipher)
	if err != nil {
		return nil, ErrWrapperEncrypt.Wrap(err)
	}
	nonce := common.CryptoRandomBytes(gcm.NonceSize())

	encrypted := gcm.Seal(nil, nonce, data, nil)
	return &EncryptedData{
		Data:         encrypted,
		Nonce:        nonce,
		KeySalt:      encryptionKeySalt,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
		HashSettings: defaultHashSettings,
	}, nil
}

func CheckPassword(givenPassword string, passwordHash []byte, passwordSalt []byte, hashSettings HashSettings) bool {
	givenPasswordHash := hashWithSalt(givenPassword, passwordSalt, hashSettings)
	return subtle.ConstantTimeCompare(givenPasswordHash, passwordHash) == 1
}

func Decrypt(password string, encryptedData *EncryptedData) ([]byte, *common.Error) {
	if !CheckPassword(password, encryptedData.PasswordHash, encryptedData.PasswordSalt, encryptedData.HashSettings) {
		return nil, ErrIncorrectPassword.RemoveHighestCategory().AddCategory(ErrTypeDecrypt)
	}

	encryptionKey := hashWithSalt(password, encryptedData.KeySalt, encryptedData.HashSettings)
	passwordCipher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, ErrWrapperDecrypt.Wrap(err)
	}

	gcm, err := cipher.NewGCM(passwordCipher)
	if err != nil {
		return nil, ErrWrapperDecrypt.Wrap(err)
	}

	decrypted, err := gcm.Open(nil, encryptedData.Nonce, encryptedData.Data, nil)
	if err != nil {
		return nil, ErrWrapperDecrypt.Wrap(err)
	}
	return decrypted, nil
}

func hash(password string, settings HashSettings) (hash, salt []byte) {
	salt = common.CryptoRandomBytes(SaltLength)
	hash = hashWithSalt(password, salt, settings)
	return
}

func hashWithSalt(password string, salt []byte, settings HashSettings) []byte {
	return argon2.IDKey([]byte(password), salt, settings.Time, settings.Memory, HashThreads, settings.KeyLen)
}
