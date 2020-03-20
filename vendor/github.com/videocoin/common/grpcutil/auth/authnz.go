package auth

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	healthCheckPath = "/grpc.health.v1.Health/Check"
)

// JWTAuthNZ returns an authentication and authorization handler.
func JWTAuthNZ(audience string, accPublicKeyURLTemplate string, secretHMAC string) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		// skip for health check
		if path := PathFromMD(ctx); path == healthCheckPath {
			return ctx, nil
		}

		tokenStr, err := BearerFromMD(ctx)
		if err != nil {
			return nil, err
		}

		token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &jwt.StandardClaims{})
		if err != nil {
			return nil, err
		}

		var authenticator Authenticator

		if _, ok := token.Header["kid"]; ok {
			authenticator = ServiceAccountAuthN(ctx, audience, accPublicKeyURLTemplate)
		} else {
			authenticator = HMACAuthN(ctx, secretHMAC)
		}

		sub, err := authenticator.Authenticate(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		grpc_ctxtags.Extract(ctx).Set("auth.sub", sub)

		return ctx, nil
	}
}
