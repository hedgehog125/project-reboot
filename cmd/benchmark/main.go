package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/hedgehog125/project-reboot/core"
)

func main() {
	encrypted, _ := core.Encrypt([]byte("Hello world"), "pass")

	nextPasswordChan := make(chan string)
	threads := runtime.NumCPU() * 2

	for i := 0; i < threads; i++ {
		go workerLoop(nextPasswordChan, encrypted)
	}

	completedChecks := int64(-threads)
	go performanceLoop(&completedChecks)

	alphabet := []rune("abcdefghijklmnopqrstuvwxyz")
	currentPassword := make([]int32, 4)
	for {
		asString := ""
		for _, charId := range currentPassword {
			asString += string(alphabet[charId])
		}
		nextPasswordChan <- asString
		completedChecks++
		addIntArray(&currentPassword, 1, int32(len(alphabet)))
	}
}

func addIntArray(arr *[]int32, amount int32, maxValue int32) {
	remainingPlaceValueAmount := amount
	for i := len(*arr) - 1; i >= 0; i-- {
		(*arr)[i] += remainingPlaceValueAmount

		remainingPlaceValueAmount = (*arr)[i] / maxValue
		if remainingPlaceValueAmount == 0 {
			break
		}
		(*arr)[i] %= maxValue
	}
}

func workerLoop(nextPasswordChan chan string, encrypted *core.EncryptedData) {
	for {
		password := <-nextPasswordChan

		decrypted, err := core.Decrypt(password, encrypted)
		if err == nil {
			fmt.Printf("password: %v\ndecrypted:\n%v\n", password, string(decrypted))
		}
	}
}

func performanceLoop(completedChecks *int64) {
	completedChecksWas := int64(0)
	for {
		time.Sleep(time.Minute)
		c := *completedChecks

		fmt.Printf("========== Checks per minute: %v ==========\n", c-completedChecksWas)

		completedChecksWas = c
	}
}
