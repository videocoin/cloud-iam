package helpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// PrivKeyToBytesPEM returns the PEM encoding of a RSA private key.
func PrivKeyToBytesPEM(priv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
}

// PubKeyToBytesPEM returns the PEM encoding of a RSA public key.
func PubKeyToBytesPEM(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1}), nil
}

// PubKeyFromBytesPEM derives the RSA public key from its PEM encoding.
func PubKeyFromBytesPEM(pemBytes []byte) (interface{}, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	return x509.ParsePKIXPublicKey(block.Bytes)
}
