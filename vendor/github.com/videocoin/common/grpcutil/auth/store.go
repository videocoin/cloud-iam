package auth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
	role, found := c.Cache.Get(principal)
	if found {
		return role.(string), nil
	}

	role, err := c.RoleStore.GetUserRole(ctx, principal)
	if err != nil {
		return "", err
	}
	c.Cache.Set(principal, role, cache.NoExpiration)

	return role.(string), nil
}

// GetRole ...
func (c *RoleCache) GetRole(ctx context.Context, name string) (*Role, error) {
	if name == "USER_ROLE_MINER" {
		return &Role{IncludedPermissions: []string{
			"iam.serviceAccountKeys.create",
			"iam.serviceAccountKeys.list",
			"iam.serviceAccountKeys.get",
		}}, nil
	}

	return &Role{IncludedPermissions: []string{}}, nil
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
	url string // https://studio.dev.videocoin.network/api/v1/user
}

// GetUserRole ...
func (r *RoleRepo) GetUserRole(ctx context.Context, principal string) (string, error) {
	tokenStr, err := BearerFromMD(ctx)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", "https://studio.dev.videocoin.network/api/v1/user", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenStr))

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	/* fixme: Get https://studio.dev.videocoin.network/api/v1/user: x509: certificate signed by unknown authority
	res, err := cleanhttp.DefaultClient().Do(req)
	if err != nil {
		return "", err
	}
	*/

	props := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&props); err != nil {
		return "", fmt.Errorf("unable to decode JSON response: %v", err)
	}

	role, ok := props["role"]
	if !ok {
		return "", errors.New("role not available")
	}

	return role.(string), nil
}

// GetRole ...
func (r *RoleRepo) GetRole(ctx context.Context, name string) (*Role, error) {
	return &Role{}, nil
}

// NewRoleRepo ...
func NewRoleRepo() *RoleRepo {
	return &RoleRepo{}
}
