package core

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/NicoClack/cryptic-stash/common"
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
func (adminCode AdminCode) Print() {
	fmt.Printf("\n==========\n\nadmin code:\n%v\n\n==========\n\n", adminCode.String())
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
