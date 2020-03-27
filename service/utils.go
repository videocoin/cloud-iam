package service

import (
	"context"
	"crypto/rsa"
	"errors"
	"io"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	guuid "github.com/google/uuid"
	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/cloud-iam/helpers"
)

const (
	keyValidityPeriodYears = 10
	bitsRSA                = 2048
)

// generateKey generates an internal user key.
func generateKey(rand io.Reader, userID string) ([]byte, *datastore.UserKey, error) {
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

	return helpers.PrivKeyToBytesPEM(rand, priv), &datastore.UserKey{
		ID:              guuid.New().String(),
		UserID:          userID,
		PublicKeyData:   pubBytes,
		ValidAfterTime:  validAfter,
		ValidBeforeTime: validBefore,
	}, nil
}

func subjectFromCtx(ctx context.Context) (string, error) {
	token, ok := ctx.Value("token").(*jwt.Token)
	if !ok {
		return "", errors.New("invalid token info")
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", errors.New("invalid token info")
	}

	return claims.Subject, nil
}
