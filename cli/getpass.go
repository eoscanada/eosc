package cli

import (
	crypto_rand "crypto/rand"
	"fmt"
	"math/big"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func GetPassword(input string) (string, error) {
	fd := os.Stdin.Fd()
	fmt.Printf(input)
	pass, err := terminal.ReadPassword(int(fd))
	fmt.Println("")
	return string(pass), err
}

// GetConfirmation will prompt for a 4 random number from
// 1000-9999. The `prompt` should always contain a %d to display the
// confirmation security code.
func GetConfirmation(prompt string) (bool, error) {
	value, err := crypto_rand.Int(crypto_rand.Reader, big.NewInt(8999))
	if err != nil {
		return false, err
	}

	randVal := 1000 + value.Int64()

	pw, err := GetPassword(fmt.Sprintf(prompt, randVal))
	if err != nil {
		return false, err
	}

	confirmPW := fmt.Sprintf("%d", randVal)

	return pw == confirmPW, nil
}
