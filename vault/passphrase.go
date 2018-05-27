package vault

import (
	crypto_rand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	saltLength         = 16
	nonceLength        = 24
	keyLength          = 32
	shamirSecretLength = 32
)

// SealWithPassphrases uses the passphrase-based encryption.
func (v *Vault) SealWithPassphrase(passphrase string) error {
	jstr, err := json.Marshal(v.KeyBag)
	if err != nil {
		return err
	}

	var nonce [nonceLength]byte
	if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
		return err
	}

	salt := make([]byte, saltLength)
	if _, err := crypto_rand.Read(salt); err != nil {
		return err
	}
	secretKey := deriveKey(passphrase, salt)
	prefix := append(salt, nonce[:]...)

	ciphertext := secretbox.Seal(prefix, jstr, &nonce, &secretKey)

	v.SecretBoxWrap = "passphrase"
	v.SecretBoxCiphertext = base64.RawStdEncoding.EncodeToString(ciphertext)

	return nil
}

// OpenWithPassphrase will load the KeyBag from the decrypted material, unlocked with the passphrase.
func (v *Vault) OpenWithPassphrase(passphrase string) error {
	buf, err := base64.RawStdEncoding.DecodeString(v.SecretBoxCiphertext)
	if err != nil {
		return err
	}

	salt := make([]byte, saltLength)
	copy(salt, buf[:saltLength])
	var nonce [nonceLength]byte
	copy(nonce[:], buf[saltLength:nonceLength+saltLength])

	secretKey := deriveKey(passphrase, salt)
	decrypted, ok := secretbox.Open(nil, buf[nonceLength+saltLength:], &nonce, &secretKey)
	if !ok {
		return fmt.Errorf("failed to decrypt")
	}
	if err := json.Unmarshal(decrypted, &v.KeyBag); err != nil {
		return err
	}
	return nil
}

// deriveKey derives a passphrase using argon2, strenthens against
// brute force attacks. Parameters chosen as per
// https://godoc.org/golang.org/x/crypto/argon2#IDKey
func deriveKey(passphrase string, salt []byte) [keyLength]byte {
	secretKeyBytes := argon2.IDKey([]byte(passphrase), salt, 4, 64*1024, 4, 32)
	var secretKey [keyLength]byte
	copy(secretKey[:], secretKeyBytes)
	return secretKey
}
