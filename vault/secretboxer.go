package vault

import (
	crypto_rand "crypto/rand"
	"io"

	"encoding/base64"

	"fmt"

	"github.com/eoscanada/eosc/cli"
	"github.com/spf13/viper"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

type SecretBoxer interface {
	Seal(in []byte) (string, error)
	Open(in string) ([]byte, error)
	WrapType() string
}

type PassphraseBoxer struct {
	passphrase string
}

func NewPassphraseBoxer(password string) *PassphraseBoxer {
	return &PassphraseBoxer{
		passphrase: password,
	}
}

func (b *PassphraseBoxer) WrapType() string {
	return "passphrase"
}

func (b *PassphraseBoxer) Seal(in []byte) (string, error) {

	var nonce [nonceLength]byte
	if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
		return "", err
	}

	salt := make([]byte, saltLength)
	if _, err := crypto_rand.Read(salt); err != nil {
		return "", err
	}
	secretKey := deriveKey(b.passphrase, salt)
	prefix := append(salt, nonce[:]...)

	cipherText := secretbox.Seal(prefix, in, &nonce, &secretKey)

	return base64.RawStdEncoding.EncodeToString(cipherText), nil
}

func (b *PassphraseBoxer) Open(in string) ([]byte, error) {
	buf, err := base64.RawStdEncoding.DecodeString(in)
	if err != nil {
		return []byte{}, err
	}

	salt := make([]byte, saltLength)
	copy(salt, buf[:saltLength])
	var nonce [nonceLength]byte
	copy(nonce[:], buf[saltLength:nonceLength+saltLength])

	secretKey := deriveKey(b.passphrase, salt)
	decrypted, ok := secretbox.Open(nil, buf[nonceLength+saltLength:], &nonce, &secretKey)
	if !ok {
		return []byte{}, fmt.Errorf("failed to decrypt")
	}
	return decrypted, nil
}

type KMSGCPBoxer struct {
	keyRing string
}

func NewKMSCGPBoxer(keyRing string) *KMSGCPBoxer {
	return &KMSGCPBoxer{
		keyRing: keyRing,
	}
}

func (b *KMSGCPBoxer) Seal(in []byte) (string, error) {
	panic("implement me")
	//mgr, err := NewKMSGCPManager("/key/ring/here")
	//if err != nil {
	//	fmt.Errorf("new kms gcp manager, %s", err)
	//}
}

func (b *KMSGCPBoxer) Open(in string) ([]byte, error) {
	panic("implement me")
}

func (b *KMSGCPBoxer) WrapType() string {
	return "kms-gcp"
}

const (
	saltLength         = 16
	nonceLength        = 24
	keyLength          = 32
	shamirSecretLength = 32
)

func deriveKey(passphrase string, salt []byte) [keyLength]byte {
	secretKeyBytes := argon2.IDKey([]byte(passphrase), salt, 4, 64*1024, 4, 32)
	var secretKey [keyLength]byte
	copy(secretKey[:], secretKeyBytes)
	return secretKey
}

func SecretBoxerForType(boxerType string) (SecretBoxer, error) {

	switch boxerType {
	case "kms-gcp":
		return NewKMSCGPBoxer(viper.GetString("kms-keyring")), nil
	case "passphrase":
		password, err := cli.GetPassphrase()
		if err != nil {
			return nil, err
		}
		return NewPassphraseBoxer(password), nil
	case "passphrase-create":
		password, err := cli.CreatePassphrase()
		if err != nil {
			return nil, err
		}
		return NewPassphraseBoxer(password), nil
	default:
		return nil, fmt.Errorf("unknown secret boxer: %s", boxerType)
	}
}
