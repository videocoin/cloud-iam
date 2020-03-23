package auth

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	jwt "github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/google/uuid"
)

// AuthenticatorFunc ...
type AuthenticatorFunc func(ctx context.Context) (string, error)

// Authenticate ...
func (f AuthenticatorFunc) Authenticate(ctx context.Context) (string, error) {
	return f(ctx)
}

// Authenticator authenticates requests.
type Authenticator interface {
	Authenticate(ctx context.Context) (string, error)
}

// ServiceAccountAuthN handles authentication for service accounts.
func ServiceAccountAuthN(audience string, pubKeyURLTemplate string) AuthenticatorFunc {
	return func(ctx context.Context) (string, error) {
		tokenStr, err := BearerFromMD(ctx)
		if err != nil {
			return "", err
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
				return nil, fmt.Errorf("Key ID is missing")
			}

			if _, err := uuid.Parse(kid.(string)); err != nil {
				return nil, fmt.Errorf("Invalid kid: %v", kid)
			}

			return ServiceAccountPublicKey(pubKeyURLTemplate, claims.Subject, kid.(string))
		})
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					return "", errors.New("")
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					return "", errors.New("Token is either expired or not active yet")
				}
			}
			return "", fmt.Errorf("Couldn't handle this token: %v", err)
		}

		return claims.Subject, nil
	}
}

// HMACAuthN handles authentication based on JWT with HMAC protection.
func HMACAuthN(secret string) AuthenticatorFunc {
	return func(ctx context.Context) (string, error) {
		tokenStr, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return "", err
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
					return "", errors.New("")
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					return "", errors.New("Token is either expired or not active yet")
				}
			}
			return "", fmt.Errorf("Couldn't handle this token: %v", err)
		}

		return claims.Subject, nil
	}
}
