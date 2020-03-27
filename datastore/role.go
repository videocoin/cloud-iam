package datastore

import "github.com/videocoin/common/dbutil/models"

type Role struct {
	models.Base
	ID          string `gorm:"primary_key"`
	Name        string
	Title       string
	Description string
	Permissions []*Permission `gorm:"many2many:role_permissions;"`
}

// TableName sets Role's table name to be `roles`.
func (r *Role) TableName() string {
	return "roles"
}

type Permission struct {
	models.Base
	ID          string `gorm:"primary_key"`
	MethodID    string
	Name        string
	Title       string
	Description string
}

// TableName sets Role's table name to be `roles`.
func (p *Permission) TableName() string {
	return "permissions"
}
