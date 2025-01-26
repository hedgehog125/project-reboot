package util

import (
	"crypto/rand"
	"log"
)

func CryptoRandomBytes(length int) []byte {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatalf("randomBytes: couldn't get random byte. error:\n%v", err)
	}
	return salt
}
