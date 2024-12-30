package util

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func RequireEnv(name string) string {
	value, specified := os.LookupEnv(name)
	if !specified {
		log.Fatalf("required environment variable \"%v\" hasn't been specified", name)
	}

	return value
}
func RequireStrArrEnv(name string) []string {
	rawValue := RequireEnv(name)

	return strings.Split(rawValue, ",")
}
func RequireIntEnv(name string) int {
	rawValue := RequireEnv(name)

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		log.Fatalf("couldn't parse environment variable \"%v\" into an integer", name)
	}

	return value
}
func RequireInt64Env(name string) int64 {
	rawValue := RequireEnv(name)

	value, err := strconv.ParseInt(rawValue, 10, 0)
	if err != nil {
		log.Fatalf("couldn't parse environment variable \"%v\" into an integer", name)
	}

	return value
}
