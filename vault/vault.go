package vault

import (
	"encoding/json"
	"fmt"
	"os"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)

// Vault represents an `eosc` wallet.  It contains the encrypted
// material to load a KeyBag, which is the signing provider for
// signing transactions using the `eos-go` library (which includes the
// embedded keosd-compatible wallet).
type Vault struct {
	Kind    string `json:"kind"`
	Version int    `json:"version"`
	Comment string `json:"comment"`

	Shamir *ShamirWrapping `json:"shamir,omitempty"`

	SecretBoxWrap       string `json:"secretbox_wrap"`
	SecretBoxCiphertext string `json:"secretbox_ciphertext"`

	KeyBag *eos.KeyBag `json:"-"`
}

type ShamirWrapping struct {
	Shares    int      `json:"shares"`
	Threshold int      `json:"threshold"`
	Parts     [][]byte `json:"parts"`
}

func NewVaultFromWalletFile(filename string) (*Vault, error) {
	v := NewVault()
	fl, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	err = json.NewDecoder(fl).Decode(&v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func NewVaultFromKeysFile(keysFile string) (*Vault, error) {
	v := NewVault()
	return v, nil
}

func NewVaultFromSingleKey(privKey string) (*Vault, error) {
	v := NewVault()
	return v, nil
}

// NewVault returns an empty vault, unsaved and with no keys.
func NewVault() *Vault {
	return &Vault{
		Kind:    "eosc-vault-wallet",
		Version: 1,
		KeyBag:  eos.NewKeyBag(),
	}
}

// NewKeyPair creates a new EOS keypair, saves the private key in the
// local wallet and returns the public key. It does NOT save the
// wallet, you better do that soon after.
func (v *Vault) NewKeyPair() (pub ecc.PublicKey, err error) {
	privKey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return
	}

	v.KeyBag.Keys = append(v.KeyBag.Keys, privKey)

	pub = privKey.PublicKey()
	return
}

func (v *Vault) AddPrivateKey(privateKey *ecc.PrivateKey) (pub ecc.PublicKey) {
	v.KeyBag.Keys = append(v.KeyBag.Keys, privateKey)
	pub = privateKey.PublicKey()
	return
}

func (v *Vault) PrintPublicKeys() {
	fmt.Printf("Public keys contained within (%d in total):\n", len(v.KeyBag.Keys))
	for _, key := range v.KeyBag.Keys {
		fmt.Println("-", key.PublicKey().String())
	}
}

// WriteToFile writes the Vault to disk. You need to encrypt before
// writing to file, otherwise you might lose much :)
func (v *Vault) WriteToFile(filename string) error {
	cnt, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fl, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	_, err = fl.Write(cnt)
	if err != nil {
		fl.Close()
		return err
	}

	return fl.Close()
}

func (v *Vault) Open(boxer SecretBoxer) error {
	data, err := boxer.Open(v.SecretBoxCiphertext)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v.KeyBag)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) Seal(boxer SecretBoxer) error {
	payload, err := json.Marshal(v.KeyBag)
	if err != nil {
		return err
	}

	v.SecretBoxWrap = boxer.WrapType()
	cipherText, err := boxer.Seal(payload)
	if err != nil {
		return err
	}

	v.SecretBoxCiphertext = cipherText
	return nil
}
