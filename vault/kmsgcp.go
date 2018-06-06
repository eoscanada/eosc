package vault

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"sync"

	"encoding/json"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudkms/v1"
)

type KMSManager interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

func NewKMSGCPManager(keyRing string) (*KMSGCPManager, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	kmsService, err := cloudkms.New(client)
	if err != nil {
		return nil, err
	}

	manager := &KMSGCPManager{
		service: kmsService,
		keyRing: keyRing,
	}

	if err := manager.setupEncryption(); err != nil {
		return nil, err
	}

	return manager, nil
}

type KMSGCPManager struct {
	dekCache        map[string][32]byte
	dekCacheLock    sync.Mutex
	localDEK        [32]byte
	localWrappedDEK string
	service         *cloudkms.Service
	keyRing         string
}

func (k *KMSGCPManager) setupEncryption() error {
	_, err := rand.Read(k.localDEK[:])
	if err != nil {
		return err
	}

	req := &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString(k.localDEK[:]),
	}

	resp, err := k.service.Projects.Locations.KeyRings.CryptoKeys.Encrypt(k.keyRing, req).Do()
	if err != nil {
		return err
	}

	k.localWrappedDEK = resp.Ciphertext
	k.dekCache = map[string][32]byte{resp.Ciphertext: k.localDEK}

	return nil
}

func (k *KMSGCPManager) fetchPlainDEK(wrappedDEK string) (out [32]byte, err error) {
	k.dekCacheLock.Lock()
	defer k.dekCacheLock.Unlock()

	if cachedKey, found := k.dekCache[wrappedDEK]; found {
		return cachedKey, nil
	}

	req := &cloudkms.DecryptRequest{
		Ciphertext: wrappedDEK,
	}
	resp, err := k.service.Projects.Locations.KeyRings.CryptoKeys.Decrypt(k.keyRing, req).Do()
	if err != nil {
		return
	}

	plainKey, err := base64.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		return
	}

	copy(out[:], plainKey)

	k.dekCache[wrappedDEK] = out

	return
}

type BlobV1 struct {
	Version       int      `bson:"version"`
	WrappedDEK    string   `bson:"wrapped_dek"`
	Nonce         [24]byte `bson:"nonce"`
	EncryptedData []byte   `bson:"data"`
}

func (k *KMSGCPManager) Encrypt(in []byte) ([]byte, error) {

	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	var sealedMsg []byte
	sealedMsg = secretbox.Seal(sealedMsg, in, &nonce, &k.localDEK)

	blob := &BlobV1{
		Version:       1,
		WrappedDEK:    k.localWrappedDEK,
		Nonce:         nonce,
		EncryptedData: sealedMsg,
	}

	cereal, err := json.Marshal(blob)
	if err != nil {
		return nil, err
	}

	return cereal, nil
}

func (k *KMSGCPManager) Decrypt(in []byte) ([]byte, error) {
	var blob BlobV1
	err := json.Unmarshal(in, &blob)
	if err != nil {
		return nil, err
	}

	// No need to check `blob.Version` == 1, we did it already with
	// the `magicFound` comparison.

	plainDEK, err := k.fetchPlainDEK(blob.WrappedDEK)
	if err != nil {
		return nil, err
	}

	plainData, ok := secretbox.Open(nil, blob.EncryptedData, &blob.Nonce, &plainDEK)
	if !ok {
		return nil, fmt.Errorf("failed decrypting data, that's all we know")
	}

	return plainData, nil
}

// Passthrough encryption (no encryption, that is)

func NewPassthroughKeyManager() *PassthroughKeyManager {
	return &PassthroughKeyManager{}
}

type PassthroughKeyManager struct{}

func (k PassthroughKeyManager) Encrypt(in []byte) ([]byte, error) { return in, nil }
func (k PassthroughKeyManager) Decrypt(in []byte) ([]byte, error) { return in, nil }
