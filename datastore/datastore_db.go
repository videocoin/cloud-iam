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
func (db *database) CreateServiceAccount(acc *ServiceAccount) (*ServiceAccount, error) {
	if err := db.Create(acc).Error; err != nil {
		return nil, err
	}
	return acc, nil
}

// GetServiceAccount gets a service account.
func (db *database) GetServiceAccountByEmail(userID string, email string) (*ServiceAccount, error) {
	return getServiceAccountByEmail(db.DB, email)
}

func getServiceAccountByEmail(DB *gorm.DB, email string) (*ServiceAccount, error) {
	acc := &ServiceAccount{}
	if err := DB.Find(acc, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return acc, nil
}

// ListServiceAccounts lists service accounts for a project.
func (db *database) ListServiceAccounts(projID string) ([]*ServiceAccount, error) {
	accs := []*ServiceAccount{}
	if err := db.Find(&accs, "project_id = ?", projID).Error; err != nil {
		return nil, err
	}
	return accs, nil
}

// DeleteServiceAccount deletes a service account.
func (db *database) DeleteServiceAccount(email string) error {
	return db.Delete(ServiceAccount{}, "email = ?", email).Error
}

// CreateServiceAccountKey creates a service account key.
func (db *database) CreateServiceAccountKey(accEmail string, passphrase string) (*ServiceAccountKey, error) {
	var key *ServiceAccountKey

	err := db.Transaction(func(tx *gorm.DB) error {
		acc, err := getServiceAccountByEmail(tx, accEmail)
		if err != nil {
			return err
		}
		key, err = generateKey(rand.Reader, passphrase, acc.ID)
		if err != nil {
			return err
		}
		if err := tx.Create(key).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GetServiceAccountKey gets a service account key.
func (db *database) GetServiceAccountKey(ID string) (*ServiceAccountKey, error) {
	key := &ServiceAccountKey{}
	if err := db.Find(key, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return key, nil
}

func listServiceAccountKeys(DB *gorm.DB, accID string) ([]*ServiceAccountKey, error) {
	keys := []*ServiceAccountKey{}
	if err := DB.Find(&keys, "account_id = ?", accID).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

// ListServiceAccountKeys lists the service account keys.
func (db *database) ListServiceAccountKeysByEmail(accEmail string) ([]*ServiceAccountKey, error) {
	var keys []*ServiceAccountKey

	err := db.Transaction(func(tx *gorm.DB) error {
		// replace with join?
		acc, err := getServiceAccountByEmail(tx, accEmail)
		if err != nil {
			return err
		}
		keys, err = listServiceAccountKeys(tx, acc.ID)
		if err != nil {
			return nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// DeleteServiceAccountKey deletes a service account key.
func (db *database) DeleteServiceAccountKey(id string) error {
	return db.Delete(ServiceAccountKey{ID: id}).Error
}
