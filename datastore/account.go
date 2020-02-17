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

// ServiceAccount is an account that belongs to your project instead
// of to an individual end user. It is used to authenticate calls
// to a VideoCoin API.
type ServiceAccount struct {
	models.Base
	ID     string `gorm:"primary_key"`
	UserID string
	Email  string
	Keys   []ServiceAccountKey `gorm:"foreignkey:AccountID"`
}

// TableName set ServiceAccount's table name to be `accounts`.
func (sa *ServiceAccount) TableName() string {
	return "accounts"
}

// Proto ...
func (sa *ServiceAccount) Proto() *iam.ServiceAccount {
	return &iam.ServiceAccount{
		UniqueId: sa.ID,
		Email:    sa.Email,
	}
}

// ServiceAccountKey represents a service account key.
type ServiceAccountKey struct {
	models.Base
	ID              string `gorm:"primary_key"`
	AccountID       string
	PrivateKeyData  []byte
	PublicKeyData   []byte
	ValidAfterTime  time.Time
	ValidBeforeTime time.Time
}

// TableName set ServiceAccountKey's table name to be `account_keys`.
func (k *ServiceAccountKey) TableName() string {
	return "account_keys" // note: 'keys' is a reserved word in mysql
}

// CreationProto ...
func (k *ServiceAccountKey) CreationProto(passphrase string) (*iam.ServiceAccountKey, error) {
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

// Proto ...
func (k *ServiceAccountKey) Proto() (*iam.ServiceAccountKey, error) {
	validAfterTimePB, err := types.TimestampProto(k.ValidAfterTime)
	if err != nil {
		return nil, err
	}
	validBeforeTimePB, err := types.TimestampProto(k.ValidBeforeTime)
	if err != nil {

		return nil, err
	}

	return &iam.ServiceAccountKey{
		Id:              k.ID,
		ValidAfterTime:  validBeforeTimePB,
		ValidBeforeTime: validAfterTimePB,
	}, nil
}
