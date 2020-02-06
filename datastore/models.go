package datastore

import (
	"time"

	"github.com/videocoin/common/models"
)

// ServiceAccount is an account that belongs to your project instead
// of to an individual end user. It is used to authenticate calls
// to a VideoCoin API.
type ServiceAccount struct {
	models.Base
	ID        string `gorm:"primary_key"`
	ProjectID string
	Email     string
	Keys      []ServiceAccountKey `gorm:"foreignkey:AccountID"`
}

// TableName set ServiceAccount's table name to be `accounts`.
func (sa *ServiceAccount) TableName() string {
	return "accounts"
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
func (key *ServiceAccountKey) TableName() string {
	return "account_keys" // note: 'keys' is a reserved word in mysql
}
