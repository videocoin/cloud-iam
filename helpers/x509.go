package helpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
)

// PrivKeyToBytesPEM returns the PEM encoding of the given
// DER-encoded private key.
func PrivKeyToBytesPEM(random io.Reader, priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
}

// PubKeyToBytesPEM returns the PEM encoding the DER-encoded public key.
func PubKeyToBytesPEM(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1}), nil
}

func PubKeyFromBytesPEM(pemBytes []byte) (interface{}, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	return x509.ParsePKIXPublicKey(block.Bytes)
}
