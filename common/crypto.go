package common

import (
	"crypto/rand"
	"fmt"
)

func CryptoRandomBytes(length int) []byte {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		UnrecoverablePanic(fmt.Sprintf("CryptoRandomBytes: couldn't get random byte. error:\n%v", err))
	}
	return salt
}
