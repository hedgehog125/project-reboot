package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hedgehog125/project-reboot/core"
)

type guess struct {
	password         string
	decryptedContent string
}

func main() {
	password := flag.String("password", "", "the password to try to guess")
	flag.Parse()
	if *password == "" {
		log.Fatalf("missing required argument \"password\"")
	}

	fmt.Println("benchmarking...")
	threads := runtime.GOMAXPROCS(0)
	fmt.Printf("running on %v threads\n\n", threads)

	encrypted, err := core.Encrypt([]byte("Hello world"), *password)
	if err != nil {
		log.Fatalf("unable to encrypt test data. error:\n%v", err.Error())
	}

	startTime := time.Now().UTC()
	nextPasswordChan := make(chan string, threads)
	guessChan := make(chan guess)

	for range threads {
		go workerLoop(nextPasswordChan, guessChan, encrypted)
	}

	alphabet := []rune("abcdefghijklmnopqrstuvwxyz")
	currentPassword := make([]int32, len(*password))

	completedChecks := int64(-threads)
	go performanceLoop(&completedChecks, currentPassword)

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

	fmt.Printf("\nsuccessfully guessed password after ~%v attempts in %v seconds: \"%v\"\ndecrypted content:\n%v\n",
		completedChecks,
		math.Round(time.Now().UTC().Sub(startTime).Seconds()),
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
	nextPasswordChan chan string, guessChan chan guess,
	encrypted *core.EncryptedData,
) {
	for {
		select {
		case password := <-nextPasswordChan:
			decrypted, err := core.Decrypt(password, encrypted)
			if err == nil {
				guessChan <- (
				//exhaustruct:enforce
				guess{
					password:         password,
					decryptedContent: string(decrypted),
				})
			}
		case <-guessChan:
			break
		}
	}
}

func performanceLoop(completedChecks *int64, currentPassword []int32) {
	completedChecksWas := int64(0)
	for {
		time.Sleep(time.Minute)
		c := *completedChecks

		asStrings := make([]string, len(currentPassword))
		for i, charID := range currentPassword {
			asStrings[i] = strconv.Itoa(int(charID))
		}

		fmt.Printf("\nChecks per minute: %v\nCurrent guess: [%v]\n", c-completedChecksWas, strings.Join(asStrings, ", "))

		completedChecksWas = c
	}
}
