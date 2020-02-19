package datastore

import (
	"crypto/rsa"
	"io"
	"time"

	guuid "github.com/google/uuid"

	"github.com/videocoin/cloud-iam/helpers"
)

const (
	keyValidityPeriodYears = 10
	bitsRSA                = 2048
)

// generateKey generates an internal user key.
func generateKey(rand io.Reader, passphrase string, userID string) (*UserKey, error) {
	key, err := rsa.GenerateKey(rand, bitsRSA)
	if err != nil {
		return nil, err
	}

	keyBytes, err := helpers.PrivKeyToBytesWithPassphrasePEM(rand, key, passphrase)
	if err != nil {
		return nil, err
	}

	pubBytes, err := helpers.PubKeyToBytesPEM(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	return &UserKey{
		ID:              guuid.New().String(),
		UserID:          userID,
		PrivateKeyData:  keyBytes,
		PublicKeyData:   pubBytes,
		ValidAfterTime:  time.Now(),
		ValidBeforeTime: time.Now().AddDate(keyValidityPeriodYears, 0, 0),
	}, nil
}
