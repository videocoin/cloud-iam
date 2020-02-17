package datastore

import (
	"github.com/videocoin/cloud-pkg/dbutil/models"
)

// Policy is a collection of bindings.
type Policy struct {
	models.Base
	ID       string `gorm:"primary_key"`
	version  int
	Bindings []Binding `gorm:"foreignkey:PolicyID"`
}

// TableName set Policy's table name to be `policies`.
func (p *Policy) TableName() string {
	return "policies"
}

// Binding binds one member to a single role.
type Binding struct {
	models.Base
	ID       string `gorm:"primary_key"`
	PolicyID string
	Role     string
	Member   string
}

// TableName set Binding's table name to be `bindings`.
func (b *Binding) TableName() string {
	return "bindings"
}
