package service

import (
	"crypto/rsa"
	"io"
	"time"

	guuid "github.com/google/uuid"

	"github.com/videocoin/cloud-iam/datastore/models"
	"github.com/videocoin/cloud-iam/helpers"
)

const (
	keyValidityPeriodYears = 10
	bitsRSA                = 2048
)

// generateKey generates an internal user key.
func generateKey(rand io.Reader, userID string) ([]byte, *models.UserKey, error) {
	priv, err := rsa.GenerateKey(rand, bitsRSA)
	if err != nil {
		return nil, nil, err
	}

	pubBytes, err := helpers.PubKeyToBytesPEM(&priv.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	validAfter := time.Now()
	validBefore := time.Now().AddDate(keyValidityPeriodYears, 0, 0)

	return helpers.PrivKeyToBytesPEM(priv), &models.UserKey{
		ID:              guuid.New().String(),
		UserID:          userID,
		PublicKeyData:   pubBytes,
		ValidAfterTime:  validAfter,
		ValidBeforeTime: validBefore,
	}, nil
}
