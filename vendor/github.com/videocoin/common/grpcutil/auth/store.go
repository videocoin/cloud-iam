package auth

import (
	"context"
	"fmt"

	cache "github.com/patrickmn/go-cache"
)

// PermissionStore ...
type PermissionStore interface {
	GetPermission(ctx context.Context, method string) (string, error)
}

// PermissionCache ...
type PermissionCache struct {
	PermissionStore
	*cache.Cache
}

// GetPermission ...
func (c *PermissionCache) GetPermission(ctx context.Context, fullMethod string) (string, error) {
	perm, ok := c.Get(fullMethod)
	if !ok {
		return "", fmt.Errorf("Invalid method %s", fullMethod)
	}
	return perm.(string), nil
}

// NewPermissionCache ...
func NewPermissionCache(store PermissionStore) *PermissionCache {
	c := cache.New(cache.NoExpiration, 0)
	c.Set("/videocoin.iam.v1.IAM/CreateKey", "iam.serviceAccountKeys.create", cache.NoExpiration)
	c.Set("/videocoin.iam.v1.IAM/ListKeys", "iam.serviceAccountKeys.list", cache.NoExpiration)
	c.Set("/videocoin.iam.v1.IAM/GetKey", "iam.serviceAccountKeys.get", cache.NoExpiration)
	c.Set("/videocoin.iam.v1.IAM/DeleteKey", "iam.serviceAccountKeys.delete", cache.NoExpiration)

	return &PermissionCache{
		Cache:           c,
		PermissionStore: store,
	}
}

// PermissionRepo ...
type PermissionRepo struct {
	// TODO
}

// GetPermission ...
func (c *PermissionRepo) GetPermission(ctx context.Context, method string) (string, error) {
	return "", nil
}

// NewPermissionRepo ...
func NewPermissionRepo() *PermissionRepo {
	return &PermissionRepo{}
}

// Role ...
type Role struct {
	IncludedPermissions []string
}

// RoleStore ...
type RoleStore interface {
	GetUserRole(ctx context.Context, principal string) (string, error)
	GetRole(ctx context.Context, name string) (*Role, error)
}

// RoleCache ...
type RoleCache struct {
	RoleStore
	*cache.Cache
}

// GetUserRole ...
func (c *RoleCache) GetUserRole(ctx context.Context, principal string) (string, error) {
	return "MINER", nil
}

// GetRole ...
func (c *RoleCache) GetRole(ctx context.Context, name string) (*Role, error) {
	// does not have delete for testing purposes
	return &Role{IncludedPermissions: []string{
		"iam.serviceAccountKeys.create",
		"iam.serviceAccountKeys.list",
		"iam.serviceAccountKeys.get",
	}}, nil
}

// NewRoleCache ...
func NewRoleCache(store RoleStore) *RoleCache {
	return &RoleCache{
		Cache:     cache.New(cache.NoExpiration, 0),
		RoleStore: store,
	}
}

// RoleRepo ...
type RoleRepo struct {
	// TODO
}

// GetUserRole ...
func (r *RoleRepo) GetUserRole(ctx context.Context, principal string) (string, error) {
	return "", nil
}

// GetRole ...
func (r *RoleRepo) GetRole(ctx context.Context, name string) (*Role, error) {
	return &Role{}, nil
}

// NewRoleRepo ...
func NewRoleRepo() *RoleRepo {
	return &RoleRepo{}
}
