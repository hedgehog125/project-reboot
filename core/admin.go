package core

import (
	"crypto/subtle"
	"encoding/base64"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/pquerna/otp/totp"
)

// Doubled because the bytes are represented as base64
const AdminCodeByteLength = 128

type AdminCode []byte

func NewAdminCode() AdminCode {
	return common.CryptoRandomBytes(AdminCodeByteLength)
}
func (adminCode AdminCode) String() string {
	return base64.StdEncoding.EncodeToString(adminCode)
}

func CheckAdminCode(givenCode string, expected AdminCode, logger common.Logger) bool {
	if len(expected) != AdminCodeByteLength { // Failsafe in case this is somehow unset or only partly written
		logger.Error(
			"current admin code is the wrong length, this should not happen!",
			"length", len(expected),
			"expectedLength", AdminCodeByteLength,
		)
		return false
	}
	givenBytes, stdErr := base64.StdEncoding.DecodeString(givenCode)
	if stdErr != nil {
		return false
	}
	return subtle.ConstantTimeCompare(givenBytes, expected) == 1
}

func CheckAdminCredentials(
	password string,
	totpCode string,
	expectedHash []byte,
	salt []byte,
	settings *common.PasswordHashSettings,
	totpSecret string,
) bool {
	encryptionKey := HashPassword(password, salt, settings)
	isHashValid := subtle.ConstantTimeCompare(encryptionKey, expectedHash)
	isTotpValidBool := totp.Validate(totpCode, totpSecret)
	isTotpValid := 0
	if isTotpValidBool {
		isTotpValid = 1
	}

	return isHashValid&isTotpValid == 1
}
