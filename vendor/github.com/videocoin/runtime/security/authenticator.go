package security

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/videocoin/runtime"
)

var (
	ErrTokenExpiredOrNotActive = errors.New("token is either expired or not active yet")
	ErrMalformedToken          = errors.New("malformed token")
	ErrKidRequired             = errors.New("kid required")
)

var (
	headerAuthorize = "authorization"
)

// ServiceAccount handles authentication for service accounts.
func ServiceAccount(audience string, hmacSecret string, pubKeyFunc PubKeyFunc) runtime.AuthenticatorFunc {
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
			context.WithValue(ctx, tokenKey{}, userInfo.HMACToken)
			return userInfo.ID, nil
		}

		claims := new(jwt.StandardClaims)
		_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			url, err := url.Parse(claims.Audience)
			if err != nil {
				return nil, err
			}

			if _, err := uuid.Parse(claims.Subject); err != nil {
				return nil, fmt.Errorf("Invalid subject: %s", claims.Subject)
			}

			if url.Hostname() != audience {
				return nil, fmt.Errorf("Unexpected audience: %s", claims.Audience)
			}

			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"]
			if !ok {
				return nil, ErrKidRequired
			}

			if _, err := uuid.Parse(kid.(string)); err != nil {
				return nil, fmt.Errorf("Invalid kid: %v", kid)
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
			return "", fmt.Errorf("Couldn't handle this token: %v", err)
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{Subject: claims.Subject, IssuedAt: time.Now().Unix()})
		tokenStr, _ = token.SignedString(hmacSecretBytes)
		jwtCache.Add(tokenStr, &UserInfo{ID: claims.Subject, HMACToken: tokenStr}, time.Unix(claims.ExpiresAt, 0))
		context.WithValue(ctx, tokenKey{}, tokenStr)

		return claims.Subject, nil
	}
}

// HMACJWT handles authentication based on JWT with HMAC protection.
func HMACJWT(secret string) runtime.AuthenticatorFunc {
	jwtCache := NewJWTCache(JwtCacheSize)

	return func(ctx context.Context) (interface{}, error) {
		tokenStr, err := grpc_auth.AuthFromMD(ctx, "Bearer")
		if err != nil {
			return "", err
		}

		userInfo, found := jwtCache.Get(tokenStr)
		if found {
			context.WithValue(ctx, tokenKey{}, tokenStr)
			return userInfo.ID, nil
		}

		claims := new(jwt.StandardClaims)
		_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			if _, err := uuid.Parse(claims.Subject); err != nil {
				return nil, fmt.Errorf("Invalid subject: %s", claims.Subject)
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
			return "", fmt.Errorf("Couldn't handle this token: %v", err)
		}

		jwtCache.Add(tokenStr, &UserInfo{ID: claims.Subject}, time.Unix(claims.ExpiresAt, 0))
		context.WithValue(ctx, tokenKey{}, tokenStr)

		return claims.Subject, nil
	}
}
