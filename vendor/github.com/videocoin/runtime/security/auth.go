package security

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/videocoin/runtime"
	"github.com/videocoin/runtime/middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PubKeyFunc func(ctx context.Context, subject string, keyID string) (interface{}, error)

// Authnz returns an authentication and authorization handler for JWT-based auth.
func Authnz(audience string, hmacSecret string, pubKeyFunc PubKeyFunc) auth.AuthFunc {
	var (
		svcacc runtime.Authenticator = ServiceAccount(audience, hmacSecret, pubKeyFunc)
		hmac   runtime.Authenticator = HMACJWT(hmacSecret)
		rbac   runtime.Authorizer    = RBAC()
	)

	return func(ctx context.Context, fullMethod string) (context.Context, error) {
		tokenStr, err := grpc_auth.AuthFromMD(ctx, "Bearer")
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		token, err := parseHeader(tokenStr)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		var authenticator runtime.Authenticator
		if _, ok := token.Header["kid"]; ok {
			authenticator = svcacc
		} else {
			authenticator = hmac
		}

		user, err := authenticator.Authenticate(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		if err := rbac.Authorize(ctx, user, fullMethod); err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		return context.WithValue(ctx, userKey{}, user), nil
	}
}
