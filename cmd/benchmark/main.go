package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/core"
)

type guess struct {
	password         string
	decryptedContent string
}

func main() {
	// TODO: revise these explanations, runtime.GC() made a massive difference
	password := flag.String("password", "", "the password to try to guess")
	hashTime := flag.Uint("hash-time", 0, "the time parameter for Argon2ID")
	hashMemory := flag.Uint("hash-memory", 0, "the memory parameter for Argon2ID in KiB")
	hashThreads := flag.Uint(
		"hash-threads",
		0,
		"the threads parameter for Argon2ID (note: changing this affects the hashes produced)",
	)
	benchmarkThreads := flag.Uint(
		"benchmark-threads",
		0,
		"the number of simultaneous decryptions to run, should be at most ceil(CPU threads / hash-threads). "+
			"But ensure you have sufficient RAM to avoid slowdown due to swap "+
			"(note: each benchmark thread often consumes twice of hash-memory)",
	)
	spacing := flag.Uint(
		"spacing",
		0,
		"the time in ms that each thread should wait before trying the next attempt."+
			" setting this allows the garbage collector to reduce the average RAM usage",
	)
	flag.Parse()
	if *password == "" {
		log.Fatalf("missing required argument \"password\"")
	}
	if *hashTime == 0 {
		log.Fatalf("missing required argument \"hash-time\"")
	}
	if *hashMemory == 0 {
		log.Fatalf("missing required argument \"hash-memory\"")
	}
	if *hashThreads == 0 {
		log.Fatalf("missing required argument \"hash-threads\"")
	}
	if *benchmarkThreads == 0 {
		log.Fatalf("missing required argument \"benchmark-threads\"")
	}
	hashSettings := &common.PasswordHashSettings{
		Time:    uint32(*hashTime),
		Memory:  uint32(*hashMemory),
		Threads: uint8(*hashThreads),
	}

	fmt.Fprintln(os.Stdout, "benchmarking...")

	salt := core.GenerateSalt()
	encryptionKey := core.HashPassword(*password, salt, hashSettings)
	encrypted, nonce, err := core.Encrypt([]byte("Hello world"), encryptionKey)
	if err != nil {
		log.Fatalf("unable to encrypt test data. error:\n%v", err.Error())
	}

	fmt.Fprintf(os.Stdout, "running on %v threads\n\n", *benchmarkThreads)

	startTime := time.Now()
	nextPasswordChan := make(chan string, *benchmarkThreads)
	guessChan := make(chan guess)

	for range *benchmarkThreads {
		go workerLoop(
			nextPasswordChan, guessChan,
			time.Duration(*spacing)*time.Millisecond,
			salt, hashSettings, encrypted, nonce,
		)
	}

	alphabet := []rune("abcdefghijklmnopqrstuvwxyz")
	currentPassword := make([]int32, len(*password))

	completedChecks := int64(-*benchmarkThreads)
	go performanceLoop(&completedChecks, currentPassword, *benchmarkThreads)

	var successfulGuess guess
MainLoop:
	for {
		asString := ""
		for _, charID := range currentPassword {
			asString += string(alphabet[charID])
		}

		select {
		case nextPasswordChan <- asString:
			completedChecks++
			hasOverflowed := addIntArray(currentPassword, 1,
				//#nosec - this is a constant that should always be in range
				int32(len(alphabet)),
			)
			if hasOverflowed {
				panic("couldn't find password after trying all combinations (with limitations)")
			}
		case successfulGuess = <-guessChan:
			break MainLoop
		}
	}

	fmt.Fprintf(os.Stdout,
		"\nsuccessfully guessed password after ~%v attempts in %v seconds: \"%v\"\ndecrypted content:\n%v\n",
		completedChecks,
		math.Round(time.Since(startTime).Seconds()),
		successfulGuess.password,
		successfulGuess.decryptedContent,
	)
}

func addIntArray(arr []int32, amount int32, maxValue int32) bool {
	hasOverflowed := false
	remainingPlaceValueAmount := amount
	for digitIndex := len(arr) - 1; digitIndex >= 0; digitIndex-- {
		arr[digitIndex] += remainingPlaceValueAmount

		remainingPlaceValueAmount = arr[digitIndex] / maxValue
		if remainingPlaceValueAmount == 0 {
			break
		}
		if digitIndex == 0 && arr[digitIndex] >= maxValue {
			hasOverflowed = true
		}
		arr[digitIndex] %= maxValue
	}

	return hasOverflowed
}

func workerLoop(
	nextPasswordChan chan string, guessChan chan guess, spacing time.Duration,
	salt []byte, passwordHashSettings *common.PasswordHashSettings, encrypted []byte, nonce []byte,
) {
	for {
		select {
		case password := <-nextPasswordChan:
			encryptionKey := core.HashPassword(password, salt, passwordHashSettings)
			decrypted, err := core.Decrypt(encrypted, encryptionKey, nonce)
			if err == nil {
				guessChan <- (
				//exhaustruct:enforce
				guess{
					password:         password,
					decryptedContent: string(decrypted),
				})
			}
		case <-guessChan:
			return
		}
		time.Sleep(spacing)
	}
}

func performanceLoop(completedChecksPointer *int64, currentPassword []int32, benchmarkThreads uint) {
	completedChecksWas := int64(0)
	for {
		time.Sleep(time.Minute)
		completedChecks := *completedChecksPointer

		asStrings := make([]string, len(currentPassword))
		for i, charID := range currentPassword {
			asStrings[i] = strconv.Itoa(int(charID))
		}

		completedChange := completedChecks - completedChecksWas
		fmt.Fprintf(
			os.Stdout,
			"\nTotal attempts per minute: %v\nLatency per guess: %vms "+
				"(not average time, which decreases with more benchmark threads)\nCurrent guess: [%v]\n",
			completedChange,
			math.Round(
				(60_000/
					float64(completedChange))*
					float64(benchmarkThreads),
			),
			strings.Join(asStrings, ", "),
		)

		completedChecksWas = completedChecks
	}
}
