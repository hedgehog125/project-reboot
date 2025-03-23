package core

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/hedgehog125/project-reboot/common"
)

// Doubled because the bytes are represented as base64
const AdminCodeByteLength = 128

func UpdateAdminCode(state *common.State) {
	<-state.AdminCode

	adminCode := common.CryptoRandomBytes(AdminCodeByteLength)
	fmt.Printf("\n==========\n\nadmin code:\n%v\n\n==========\n\n", base64.StdEncoding.EncodeToString(adminCode))

	go func() { state.AdminCode <- adminCode }()
}

func CheckAdminCode(givenCode string, state *common.State) bool {
	currentCode := <-state.AdminCode
	go func() { state.AdminCode <- currentCode }()
	if len(currentCode) != AdminCodeByteLength { // Failsafe in case this is somehow unset or only partly written
		return false
	}

	givenBytes, err := base64.StdEncoding.DecodeString(givenCode)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(givenBytes, currentCode) == 1
}
