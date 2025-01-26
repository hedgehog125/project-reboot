package core

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/hedgehog125/project-reboot/intertypes"
	"github.com/hedgehog125/project-reboot/util"
)

// Doubled because the bytes are represented as base64
const ADMIN_CODE_BYTE_LENGTH = 128

func UpdateAdminCode(state *intertypes.State) {
	<-state.AdminCode

	adminCode := util.CryptoRandomBytes(ADMIN_CODE_BYTE_LENGTH)
	fmt.Printf("\n==========\n\nadmin code:\n%v\n\n==========\n\n", base64.StdEncoding.EncodeToString(adminCode))

	go func() { state.AdminCode <- adminCode }()
}

func CheckAdminCode(givenCode string, state *intertypes.State) bool {
	currentCode := <-state.AdminCode
	go func() { state.AdminCode <- currentCode }()
	if len(currentCode) != ADMIN_CODE_BYTE_LENGTH { // Failsafe in case this is somehow unset or only partly written
		return false
	}

	givenBytes, err := base64.StdEncoding.DecodeString(givenCode)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(givenBytes, currentCode) == 1
}
