package vault

import (
	"fmt"

	"github.com/eoscanada/eosc/cli"
	"github.com/pkg/errors"
)

type SecretBoxer interface {
	Seal(in []byte) (string, error)
	Open(in string) ([]byte, error)
	WrapType() string
}

func SecretBoxerForType(boxerType string, keypath string) (SecretBoxer, error) {
	switch boxerType {
	case "kms-gcp":
		if keypath == "" {
			return nil, errors.New("missing kms-gcp keypath")
		}
		return NewKMSGCPBoxer(keypath), nil
	case "passphrase":
		password, err := cli.GetDecryptPassphrase()
		if err != nil {
			return nil, err
		}
		return NewPassphraseBoxer(password), nil
	default:
		return nil, fmt.Errorf("unknown secret boxer: %s", boxerType)
	}
}
