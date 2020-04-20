package auth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

var (
	ErrTokenExpiredOrNotActive = errors.New("token is either expired or not active yet")
	ErrMalformedToken          = errors.New("malformed token")
	ErrKidRequired             = errors.New("kid required")
)

// ServiceAccount handles authentication for service accounts.
func ServiceAccount(audience string, hmacSecret string, pubKeyFunc PubKeyFunc) AuthenticatorFunc {
	var (
		jwtCache        = NewJWTCache(JwtCacheSize)
		hmacSecretBytes = []byte(hmacSecret)
	)

	return func(ctx context.Context) (interface{}, error) {
		tokenStr, err := grpc_auth.AuthFromMD(ctx, "Bearer")
		if err != nil {
			return "", err
		}

		userInfo, found := jwtCache.Get(tokenStr)
		if found {
			return userInfo, nil
		}

		claims := new(jwt.StandardClaims)
		_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			url, err := url.Parse(claims.Audience)
			if err != nil {
				return nil, err
			}

			if _, err := uuid.Parse(claims.Subject); err != nil {
				return nil, fmt.Errorf("invalid subject: %s", claims.Subject)
			}

			if url.Hostname() != audience {
				return nil, fmt.Errorf("unexpected audience: %s", claims.Audience)
			}

			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"]
			if !ok {
				return nil, ErrKidRequired
			}

			if _, err := uuid.Parse(kid.(string)); err != nil {
				return nil, fmt.Errorf("invalid kid: %v", kid)
			}

			return pubKeyFunc(ctx, claims.Subject, kid.(string))
		})
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					return "", ErrMalformedToken
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					return "", ErrTokenExpiredOrNotActive
				}
			}
			return "", fmt.Errorf("couldn't handle this token: %v", err)
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{Subject: claims.Subject, IssuedAt: time.Now().Unix()})
		tokenStr, _ = token.SignedString(hmacSecretBytes)
		userInfo = &UserInfo{ID: claims.Subject, HMACToken: tokenStr}
		jwtCache.Add(tokenStr, userInfo, time.Unix(claims.ExpiresAt, 0))

		return userInfo, nil
	}
}

// HMACJWT handles authentication based on JWT with HMAC protection.
func HMACJWT(secret string) AuthenticatorFunc {
	jwtCache := NewJWTCache(JwtCacheSize)

	return func(ctx context.Context) (interface{}, error) {
		tokenStr, err := grpc_auth.AuthFromMD(ctx, "Bearer")
		if err != nil {
			return "", err
		}

		userInfo, found := jwtCache.Get(tokenStr)
		if found {
			return userInfo, nil
		}

		claims := new(jwt.StandardClaims)
		_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			if _, err := uuid.Parse(claims.Subject); err != nil {
				return nil, fmt.Errorf("invalid subject: %s", claims.Subject)
			}

			return []byte(secret), nil
		})
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					return "", ErrMalformedToken
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					return "", ErrTokenExpiredOrNotActive
				}
			}
			return "", fmt.Errorf("couldn't handle this token: %v", err)
		}

		userInfo = &UserInfo{ID: claims.Subject}
		jwtCache.Add(tokenStr, userInfo, time.Unix(claims.ExpiresAt, 0))

		return userInfo, nil
	}
}
