package datastore

import (
	"io"

	"github.com/jinzhu/gorm"
	"github.com/videocoin/cloud-pkg/dbutil/models"
)

// Role is a collection of permissions.
type Role struct {
	models.Base
	ID          string `gorm:"primary_key"`
	Name        string
	Title       string
	Description string
	Permissions []Permission `gorm:"many2many:roles_permissions"`
}

// TableName set Role's table name to be `roles`.
func (r *Role) TableName() string {
	return "roles"
}

// RoleDatastore ...
type RoleDatastore interface {
	GetRole(name string) (*Role, error)
	ListRoles() ([]*Role, error)
	io.Closer
}

// RoleMySQL ...
type RoleMySQL struct {
	*gorm.DB
}

// GetRole gets a predefined role.
func (db *RoleMySQL) GetRole(name string) (*Role, error) {
	role := &Role{}
	if err := db.Find(role, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// ListRoles lists the predefined roles.
func (db *RoleMySQL) ListRoles() ([]*Role, error) {
	roles := []*Role{}
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
