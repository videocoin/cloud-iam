package security

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/videocoin/runtime"
	"github.com/videocoin/runtime/middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	TokenHeaderKeyIDKey = "kid"
	SubCtxTagKey        = "auth.sub"
	TokenCtxKey         = "token"
)

type PubKeyFunc func(ctx context.Context, subject string, keyID string) (interface{}, error)

// Authnz returns an authentication and authorization handler for JWT-based auth.
func Authnz(audience string, hmacSecret string, pubKeyFunc PubKeyFunc) auth.AuthFunc {
	var (
		svcacc    runtime.Authenticator = ServiceAccount(audience, hmacSecret, pubKeyFunc)
		hmac      runtime.Authenticator = HMACJWT(hmacSecret)
		rbac      runtime.Authorizer    = RBAC()
		parserJWT                       = new(jwt.Parser)
	)

	return func(ctx context.Context, fullMethod string) (context.Context, error) {
		tokenStr, err := grpc_auth.AuthFromMD(ctx, "Bearer")
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		token, _, err := parserJWT.ParseUnverified(tokenStr, &jwt.StandardClaims{})
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		var authenticator runtime.Authenticator
		if _, ok := token.Header[TokenHeaderKeyIDKey]; ok {
			authenticator = svcacc
		} else {
			authenticator = hmac
		}

		sub, err := authenticator.Authenticate(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		if err := rbac.Authorize(ctx, sub, fullMethod); err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		grpc_ctxtags.Extract(ctx).Set(SubCtxTagKey, sub)
		newCtx := context.WithValue(ctx, TokenCtxKey, token)
		return newCtx, nil
	}
}
