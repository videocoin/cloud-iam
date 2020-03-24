package auth

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	SubCtxTag   = "auth.sub"
	TokenCtxKey = "token"
	TokenKid    = "kid"
)

// AuthNZFunc is the pluggable function that performs authentication.
//
// The passed in `Context` will contain the gRPC metadata.MD object (for header-based authentication) and
// the peer.Peer information that can contain transport-based credentials (e.g. `credentials.AuthInfo`).
//
// The returned context will be propagated to handlers, allowing user changes to `Context`. However,
// please make sure that the `Context` returned is a child `Context` of the one passed in.
//
// If error is returned, its `grpc.Code()` will be returned to the user as well as the verbatim message.
// Please make sure you use `codes.Unauthenticated` (lacking auth) and `codes.PermissionDenied`
// (authed, but lacking perms) appropriately.
type AuthNZFunc func(ctx context.Context, fullMethod string) (context.Context, error)

// JWTAuthNZ returns an authentication and authorization handler for JWT-based auth.
func JWTAuthNZ(audience string, accPublicKeyURLTemplate string, secretHMAC string) AuthNZFunc {
	var (
		// TODO(rgeraldes) include token cache
		accAuthN   Authenticator = ServiceAccountAuthN(audience, accPublicKeyURLTemplate)
		hmacAuthN  Authenticator = HMACAuthN(secretHMAC)
		authorizer Authorizer    = RBACAuthZ()
		parserJWT                = new(jwt.Parser)
	)

	return func(ctx context.Context, fullMethod string) (context.Context, error) {
		tokenStr, err := BearerFromMD(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		token, _, err := parserJWT.ParseUnverified(tokenStr, &jwt.StandardClaims{})
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		var authenticator Authenticator
		if _, ok := token.Header[TokenKid]; ok {
			authenticator = accAuthN
		} else {
			authenticator = hmacAuthN
		}

		sub, err := authenticator.Authenticate(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		if err := authorizer.Authorize(ctx, sub, fullMethod); err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		grpc_ctxtags.Extract(ctx).Set(SubCtxTag, sub)
		newCtx := context.WithValue(ctx, TokenCtxKey, token)
		return newCtx, nil
	}
}
