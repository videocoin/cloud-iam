package datastore

import (
	"io"

	"github.com/jinzhu/gorm"
)

// DataStore is a repository for persistently storing collections of data
// related to service accounts. Account reference is the email address or the
// unique id of the service account.
type DataStore interface {
	CreateServiceAccount(acc *ServiceAccount) (*ServiceAccount, error)
	GetServiceAccountByEmail(email string) (*ServiceAccount, error)
	ListServiceAccounts(projID string) ([]*ServiceAccount, error)
	DeleteServiceAccount(email string) error
	CreateServiceAccountKey(accEmail string, passphrase string) (*ServiceAccountKey, string, error)
	GetServiceAccountKey(id string) (*ServiceAccountKey, error)
	ListServiceAccountKeysByEmail(accEmail string) ([]*ServiceAccountKey, string, error)
	DeleteServiceAccountKey(id string) error
	io.Closer
}

// Open gets a handle for a database.
func Open(uri string) (DataStore, error) {
	db, err := gorm.Open("mysql", uri)
	if err != nil {
		return nil, err
	}
	return &database{DB: db}, nil
}
