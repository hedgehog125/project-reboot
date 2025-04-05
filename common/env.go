package common

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

func RequireBoolEnv(name string) bool {
	rawValue := RequireEnv(name)
	lower := strings.ToLower(rawValue)

	if lower == "true" {
		return true
	} else if lower == "false" {
		return false
	}

	log.Fatalf("couldn't parse environment variable \"%v\". \"%v\" is not true or false", name, rawValue)
	return false
}

// Mainly used for the CLI commands eg. BENCHMARK_THREAD_COUNT
func OptionalEnv(name string, defaultValue string) string {
	value, specified := os.LookupEnv(name)
	if specified {
		return value
	} else {
		return defaultValue
	}
}

// Mainly used for the CLI commands eg. BENCHMARK_THREAD_COUNT
func OptionalIntEnv(name string, defaultValue int) int {
	_, specified := os.LookupEnv(name)
	if specified {
		return RequireIntEnv(name)
	} else {
		return defaultValue
	}
}
