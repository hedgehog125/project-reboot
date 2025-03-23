package core

import "github.com/hedgehog125/project-reboot/common"

// Doubled because the bytes are represented as base64
const AuthCodeByteLength = 128

func RandomAuthCode() []byte {
	return common.CryptoRandomBytes(AuthCodeByteLength)
}
