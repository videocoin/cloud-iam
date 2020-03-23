package datastore

import (
	"github.com/videocoin/common/dbutil/models"
)

type Method struct {
	models.Base
	ID         string `gorm:"primary_key"`
	Name       string
	Permission *Permission
}
