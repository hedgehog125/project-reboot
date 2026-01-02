package core

import (
	"github.com/NicoClack/cryptic-stash/backend/common"
	"golang.org/x/crypto/argon2"
)

const (
	EncryptionKeyLength = 32  // Required by AES-256
	PasswordSaltLength  = 128 // Overkill but there shouldn't really be any downsides
)

func GenerateSalt() []byte {
	return common.CryptoRandomBytes(PasswordSaltLength)
}

// Returns an encryption key
func HashPassword(password string, salt []byte, settings *common.PasswordHashSettings) []byte {
	return argon2.IDKey(
		[]byte(password), salt,
		settings.Time, settings.Memory,
		settings.Threads, EncryptionKeyLength,
	)
}
