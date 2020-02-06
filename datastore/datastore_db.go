package datastore

import (
	// mysql driver
	"crypto/rand"

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

// CreateServiceAccount creates a service account.
func (db *database) CreateServiceAccount(sa *ServiceAccount) (*ServiceAccount, error) {
	if err := db.Create(sa).Error; err != nil {
		return nil, err
	}
	return sa, nil
}

// GetServiceAccount gets a service account.
func (db *database) GetServiceAccount(email string) (*ServiceAccount, error) {
	return getServiceAccount(db.DB, email)
}

func getServiceAccount(DB *gorm.DB, email string) (*ServiceAccount, error) {
	var sa ServiceAccount
	if err := DB.Where("email = ?", email).First(&sa).Error; err != nil {
		return nil, err
	}
	return &sa, nil
}

// ListServiceAccounts lists service accounts for a project.
func (db *database) ListServiceAccounts(projID string) ([]*ServiceAccount, error) {
	var accounts []*ServiceAccount
	if err := db.Where(&ServiceAccount{ProjectID: projID}).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// DeleteServiceAccount deletes a service account.
func (db *database) DeleteServiceAccount(email string) error {
	return db.Delete(ServiceAccount{}, "email = ?", email).Error
}

// CreateServiceAccountKey creates a service account key.
func (db *database) CreateServiceAccountKey(accEmail string, passphrase string) (*ServiceAccountKey, string, error) {
	var (
		keyDB  ServiceAccountKey
		projID string
	)

	if err := db.Transaction(func(tx *gorm.DB) error {
		acc, err := getServiceAccount(tx, accEmail)
		if err != nil {
			return err
		}
		key, err := generateKey(rand.Reader, passphrase, acc.ID)
		if err != nil {
			return err
		}
		if err := tx.Create(&key).Error; err != nil {
			return err
		}

		projID = acc.ProjectID

		return nil
	}); err != nil {
		return nil, "", err
	}

	return &keyDB, projID, nil
}

// GetServiceAccountKey gets a service account key.
func (db *database) GetServiceAccountKey(id string) (*ServiceAccountKey, error) {
	var key ServiceAccountKey
	if err := db.Where("id = ?", id).First(&key).Error; err != nil {
		return nil, err
	}
	return &key, nil
}

func listServiceAccountKeys(DB *gorm.DB, accID string) ([]*ServiceAccountKey, error) {
	var keys []*ServiceAccountKey
	if err := DB.Where("account_id = ?", accID).Find(&keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

// ListServiceAccountKeys lists the service account keys.
func (db *database) ListServiceAccountKeys(email string) ([]*ServiceAccountKey, string, error) {
	var (
		keys   []*ServiceAccountKey
		projID string
	)

	if err := db.Transaction(func(tx *gorm.DB) error {
		// replace with join?
		acc, err := getServiceAccount(tx, email)
		if err != nil {
			return err
		}
		keys, err = listServiceAccountKeys(tx, acc.ID)
		if err != nil {
			return nil
		}

		projID = acc.ProjectID

		return nil
	}); err != nil {
		return nil, "", err
	}

	return keys, projID, nil
}

// DeleteServiceAccountKey deletes a service account key.
func (db *database) DeleteServiceAccountKey(id string) error {
	return db.Delete(ServiceAccountKey{ID: id}).Error
}
