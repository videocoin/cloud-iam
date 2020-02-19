package datastore

import (
	// mysql driver

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/jinzhu/gorm"
)

// database implements the DataStore interface.
type database struct {
	*gorm.DB
}

// Close closes the database connection.
func (db *database) Close() error {
	return db.DB.Close()
}

// CreateKey creates an user key.
func (db *database) CreateUserKey(key *UserKey) error {
	return db.Create(key).Error
}

// GetKey gets an user key.
func (db *database) GetUserKey(userID string, keyID string) (*UserKey, error) {
	key := &UserKey{}
	if err := db.Find(key, "user_id = ? AND id = ?", userID, keyID).Error; err != nil {
		return nil, err
	}
	return key, nil
}

// ListKeys lists the user keys.
func (db *database) ListUserKeys(userID string) ([]*UserKey, error) {
	keys := []*UserKey{}
	if err := db.Find(&keys, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

// DeleteKey deletes an user key.
func (db *database) DeleteUserKey(userID string, keyID string) error {
	return db.Delete(&UserKey{}, "user_id = ? AND key_id = ?", userID, keyID).Error
}

// GetRole gets a role.
func (db *database) GetRole(name string) (*Role, error) {
	return nil, nil
}

// ListRoles lists all the available roles
func (db *database) ListRoles() ([]*Role, error) {
	return nil, nil
}

// ListUserRoles lists all the available roles
func (db *database) ListUserRoles(userID string) ([]*Role, error) {
	return nil, nil
}

// CreateRoleBinding binds a role to an user.
func (db *database) CreateRoleBinding(binding *RoleBinding) error {
	return nil
}

// DeleteRoleBinding removes the binding between a role and an user.
func (db *database) DeleteRoleBinding(roleName string, userID string) error {
	return nil
}

// ListRoleBindings lists all bindings.
func (db *database) ListRoleBindings() ([]*RoleBinding, error) {
	return nil, nil
}
