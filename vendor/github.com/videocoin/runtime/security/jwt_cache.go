package security

import (
	"time"

	lru "github.com/hashicorp/golang-lru"
)

const (
	JwtCacheTimeoutSec = 300
	JwtCacheSize       = 1024
)

type UserInfo struct {
	ID        string
	HMACToken string
}

type JWTValue struct {
	UserInfo  *UserInfo
	expiresAt time.Time
}

type JWTCache struct {
	cache *lru.Cache
}

func NewJWTCache(size int) *JWTCache {
	if size <= 0 {
		size = JwtCacheSize
	}
	cache, _ := lru.New(size)
	return &JWTCache{cache: cache}
}

func (c *JWTCache) Add(tokenStr string, userInfo *UserInfo, expiresAt time.Time) {
	timeoutAt := time.Now().Add(JwtCacheTimeoutSec * time.Second)
	var t time.Time
	if timeoutAt.Before(expiresAt) {
		t = timeoutAt
	} else {
		t = expiresAt
	}

	c.cache.Add(tokenStr, &JWTValue{
		UserInfo:  userInfo,
		expiresAt: t,
	})
}

func (c *JWTCache) Get(tokenStr string) (*UserInfo, bool) {
	val, found := c.cache.Get(tokenStr)
	if !found {
		return nil, false
	}

	jwtVal := val.(*JWTValue)
	if jwtVal.expiresAt.Before(time.Now()) {
		c.cache.Remove(tokenStr)
		return nil, false
	}

	return jwtVal.UserInfo, true
}
