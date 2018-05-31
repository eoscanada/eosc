package vault

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func GetPassword(input string) (string, error) {
	fd := os.Stdin.Fd()
	fmt.Printf(input)
	pass, err := terminal.ReadPassword(int(fd))
	return string(pass), err
}
