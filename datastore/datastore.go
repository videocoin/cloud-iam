package datastore

import (
	"io"

	"github.com/jinzhu/gorm"
)

// DataStore is a repository for persistently storing collections of data
// related to identity and access management.
type DataStore interface {
	CreateUserKey(key *UserKey) error
	GetUserKey(userID string, keyID string) (*UserKey, error)
	ListUserKeys(userID string) ([]*UserKey, error)
	DeleteUserKey(userID string, keyID string) error

	GetRole(name string) (*Role, error)
	ListRoles() ([]*Role, error)
	ListUserRoles(userID string) ([]*Role, error)
	CreateRoleBinding(binding *RoleBinding) error
	DeleteRoleBinding(roleName string, userID string) error
	ListRoleBindings() ([]*RoleBinding, error)

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
