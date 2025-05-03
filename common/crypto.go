package common

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func CryptoRandomBytes(length int) []byte {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		panic(fmt.Sprintf("CryptoRandomBytes: couldn't get random byte. error:\n%v", err))
	}
	return salt
}

func CryptoRandomInt(max int64) int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(fmt.Sprintf("CryptoRandomBytes: couldn't get random byte. error:\n%v", err))
	}
	return n.Int64()
}

func CryptoRandomAlphaNum(length int) string {
	characters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	randStr := ""
	for range length {
		randStr += string(characters[CryptoRandomInt(int64(len(characters)))])
	}
	return randStr
}
