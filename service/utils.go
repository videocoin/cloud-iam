package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"io"
	"math/big"
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

func createSelfSignedCert(notBefore time.Time, notAfter time.Time, priv *rsa.PrivateKey) ([]byte, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"VideoCoin Development Association Ltd"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, fmt.Errorf("Failed to create certificate: %v", err)
	}

	return derBytes, nil
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
