package auth

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

const (
	AuthzCacheTimeoutSec = 120
	AuthzCacheSize       = 1024

	PermissionCacheTimeoutSec = 600
	PermissionCacheSize       = 250
)

type AuthzValue struct {
	Success            bool
	RequiredPermission string
	expiresAt          time.Time
}

// AuthzCache is a local cache to expedite the authorization process. The key of
// the cache is the hash of the concatenation of JWT auth token and request
// method.
type AuthzCache struct {
	cache *lru.Cache
}

func NewAuthzCache(size int) *AuthzCache {
	if size <= 0 {
		size = AuthzCacheSize
	}
	cache, _ := lru.New(size)
	return &AuthzCache{cache: cache}
}

func (c *AuthzCache) Add(key string, val *AuthzValue) {
	val.expiresAt = time.Now().Add(AuthzCacheTimeoutSec * time.Second)
	c.cache.Add(key, val)
}

func (c *AuthzCache) Get(key string) (*AuthzValue, bool) {
	val, found := c.cache.Get(key)
	if !found {
		return nil, false
	}

	authzVal := val.(*AuthzValue)
	if authzVal.expiresAt.Before(time.Now()) {
		c.cache.Remove(key)
		return nil, false
	}

	return authzVal, true
}

func (c *AuthzCache) ComposeKey(parts ...string) string {
	hasher := md5.New()
	for _, part := range parts {
		hasher.Write([]byte(part))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
