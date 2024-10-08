package models

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/videocoin/common/dbutil/models"
	iam "github.com/videocoin/videocoinapis/videocoin/iam/v1"
)

// UserKey represents an user key.
type UserKey struct {
	models.Base
	ID              string `gorm:"primary_key"`
	UserID          string
	PublicKeyData   []byte
	ValidAfterTime  time.Time
	ValidBeforeTime time.Time
}

// TableName set Key's table name to be `user_keys`.
func (k *UserKey) TableName() string {
	return "user_keys"
}

// Proto returns an IAM key.
func (k *UserKey) Proto() (*iam.Key, error) {
	validAfterTimePB, err := ptypes.TimestampProto(k.ValidAfterTime)
	if err != nil {
		return nil, err
	}
	validBeforeTimePB, err := ptypes.TimestampProto(k.ValidBeforeTime)
	if err != nil {

		return nil, err
	}

	return &iam.Key{
		Id:              k.ID,
		ValidAfterTime:  validBeforeTimePB,
		ValidBeforeTime: validAfterTimePB,
		PublicKeyData:   k.PublicKeyData,
	}, nil
}
