package datastore

import (
	"io"

	"github.com/jinzhu/gorm"
	"github.com/videocoin/cloud-iam/datastore/models"
)

// DataStore is a repository for persistently storing collections of data
// related to identity and access management.
type DataStore interface {
	CreateUserKey(key *models.UserKey) error
	GetUserKey(userID string, keyID string) (*models.UserKey, error)
	ListUserKeys(userID string) ([]*models.UserKey, error)
	DeleteUserKey(userID string, keyID string) error

	io.Closer
}

// Open gets a handle for a database.
func Open(uri string) (DataStore, error) {
	db, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}

	return &database{DB: db}, nil
}
