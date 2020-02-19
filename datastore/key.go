package datastore

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	iam "github.com/videocoin/cloud-api/iam/v1"

	"github.com/gogo/protobuf/types"
	"github.com/videocoin/cloud-pkg/dbutil/models"
)

// ErrPEMDataNotFound is returned when no PEM data is found.
var ErrPEMDataNotFound = errors.New("pem: data not found")

// UserKey represents an user key.
type UserKey struct {
	models.Base
	ID              string `gorm:"primary_key"`
	UserID          string
	PrivateKeyData  []byte
	PublicKeyData   []byte
	ValidAfterTime  time.Time
	ValidBeforeTime time.Time
}

// TableName set Key's table name to be `user_keys`.
func (k *UserKey) TableName() string {
	return "user_keys" // note: 'keys' is a reserved word in mysql
}

// CreationProto returns an IAM key for the key creation method.
func (k *UserKey) CreationProto(passphrase string) (*iam.Key, error) {
	block, _ := pem.Decode(k.PrivateKeyData)
	if block == nil {
		return nil, ErrPEMDataNotFound
	}

	decrypted, err := x509.DecryptPEMBlock(block, []byte(passphrase))
	if err != nil {
		return nil, err
	}
	block.Bytes = decrypted

	keyPB, err := k.Proto()
	if err != nil {
		return nil, err
	}

	keyPB.PrivateKeyData = pem.EncodeToMemory(block)

	return keyPB, nil
}

// Proto returns an IAM key.
func (k *UserKey) Proto() (*iam.Key, error) {
	validAfterTimePB, err := types.TimestampProto(k.ValidAfterTime)
	if err != nil {
		return nil, err
	}
	validBeforeTimePB, err := types.TimestampProto(k.ValidBeforeTime)
	if err != nil {

		return nil, err
	}

	return &iam.Key{
		Id:              k.ID,
		ValidAfterTime:  validBeforeTimePB,
		ValidBeforeTime: validAfterTimePB,
	}, nil
}
