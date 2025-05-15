package cli

import (
	"errors"
	"fmt"
)

func GetDecryptPassphrase() (string, error) {
	passphrase, err := GetPassword("Enter passphrase to decrypt your vault: ")
	if err != nil {
		return "", fmt.Errorf("reading password: %w", err)
	}

	return passphrase, nil
}

func GetEncryptPassphrase() (string, error) {
	passphrase, err := GetPassword("Enter passphrase to encrypt your vault: ")
	if err != nil {
		return "", fmt.Errorf("reading password: %w", err)
	}

	passphraseConfirm, err := GetPassword("Confirm passphrase: ")
	if err != nil {
		return "", fmt.Errorf("reading confirmation password: %w", err)
	}

	if passphrase != passphraseConfirm {
		fmt.Println()
		return "", errors.New("passphrase mismatch!")
	}
	return passphrase, nil
}
