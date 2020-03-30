package datastore

import (
	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/jinzhu/gorm"
	"github.com/videocoin/cloud-iam/datastore/models"
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
func (db *database) CreateUserKey(key *models.UserKey) error {
	return db.Create(key).Error
}

// GetKey gets an user key.
func (db *database) GetUserKey(userID string, keyID string) (*models.UserKey, error) {
	key := &models.UserKey{}
	if err := db.Find(key, "user_id = ? AND id = ?", userID, keyID).Error; err != nil {
		return nil, err
	}
	return key, nil
}

// ListKeys lists the user keys.
func (db *database) ListUserKeys(userID string) ([]*models.UserKey, error) {
	keys := []*models.UserKey{}
	if err := db.Find(&keys, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

// DeleteKey deletes an user key.
func (db *database) DeleteUserKey(userID string, keyID string) error {
	return db.Delete(&models.UserKey{}, "user_id = ? AND key_id = ?", userID, keyID).Error
}
