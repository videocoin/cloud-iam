package auth

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	jwt "github.com/videocoin/jwt-go"
)

type AuthOption interface {
	apply(*authnz)
}

type funcAuthOption struct {
	f func(*authnz)
}

func (fdo *funcAuthOption) apply(do *authnz) {
	fdo.f(do)
}

func newFuncAuthOption(f func(*authnz)) *funcAuthOption {
	return &funcAuthOption{
		f: f,
	}
}

type authStrategy struct {
	auth       Authenticator
	matchingFn MatchingFunc
}

type MatchingFunc func(token *jwt.Token) bool

func WithAuthentication(auth Authenticator, rules ...MatchingFunc) AuthOption {
	return newFuncAuthOption(func(a *authnz) {
		if len(rules) == 0 {
			a.defaultAuthenticator = auth
			return
		}

		a.authStrategies = append(a.authStrategies, &authStrategy{
			auth:       auth,
			matchingFn: rules[0],
		})
	})
}

func WithAuthorization(auth Authorizer) AuthOption {
	return newFuncAuthOption(func(a *authnz) {
		a.authorizer = auth
	})
}

type authnz struct {
	defaultAuthenticator Authenticator
	authStrategies       []*authStrategy
	authorizer           Authorizer
}

type PubKeyFunc func(ctx context.Context, subject string, keyID string) (interface{}, error)

func NewAuthnzHandler(opt ...AuthOption) *authnz {
	handler := &authnz{}
	for _, o := range opt {
		o.apply(handler)
	}
	return handler
}

func (h *authnz) HandleAuthnz(ctx context.Context, fullMethod string) (context.Context, error) {
	if h.defaultAuthenticator != nil {
		principal, err := h.authenticate(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		ctx = context.WithValue(ctx, userKey{}, principal)

		if h.authorizer != nil {
			if err := h.authorizer.Authorize(ctx, principal, fullMethod); err != nil {
				return nil, status.Error(codes.PermissionDenied, err.Error())
			}
		}
	}

	return ctx, nil
}

func (h authnz) authenticate(ctx context.Context) (interface{}, error) {
	authenticator := h.defaultAuthenticator

	if len(h.authStrategies) != 0 {
		tokenStr, err := grpc_auth.AuthFromMD(ctx, "Bearer")
		if err != nil {
			return nil, err
		}
		token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &jwt.StandardClaims{})
		if err != nil {
			return nil, err
		}

		for _, strategy := range h.authStrategies {
			if strategy.matchingFn(token) {
				authenticator = strategy.auth
				break
			}
		}
	}

	return authenticator.Authenticate(ctx)
}
