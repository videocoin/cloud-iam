package datastore

import (
	"io"

	iam "github.com/videocoin/cloud-api/iam/v1"

	"github.com/jinzhu/gorm"
	"github.com/videocoin/cloud-pkg/dbutil/models"
)

// Role is a collection of permissions.
type Role struct {
	models.Base
	Name        string `gorm:"primary_key"`
	Title       string
	Description string
	Permissions []string `gorm:"many2many:roles_permissions"`
}

// TableName set Role's table name to be `roles`.
func (r *Role) TableName() string {
	return "roles"
}

// Proto ...
func (r *Role) Proto() *iam.Role {
	// TODO
	return &iam.Role{}
}

// RoleBinding ...
type RoleBinding struct {
	models.Base
	UserID   string
	RoleName string
}

// RoleDatastore ...
type RoleDatastore interface {
	GetRole(name string) (*Role, error)
	ListRoles() ([]*Role, error)
	CreateRoleBinding(roleName string, userID string) error
	DeleteRoleBinding(roleName string, userID string) error
	ListRoleBindings() []*RoleBinding
	ListRoleBindingsByUser(userID string) []*RoleBinding
	io.Closer
}

// RoleCache ...
type RoleCache struct {
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

// Permission provides access to specific resources.
type Permission struct {
	models.Base
	Name        string `gorm:"primary_key"`
	Description string
	Roles       []Role `gorm:"many2many:roles_permissions"`
}

// TableName set Permission's table name to be `permissions`.
func (p *Permission) TableName() string {
	return "permissions"
}
