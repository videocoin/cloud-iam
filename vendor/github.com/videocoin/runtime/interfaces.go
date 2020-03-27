package runtime

import (
	"context"
)

// AuthenticatorFunc turns a function into an authenticator
type AuthenticatorFunc func(ctx context.Context) (interface{}, error)

func (f AuthenticatorFunc) Authenticate(ctx context.Context) (interface{}, error) {
	return f(ctx)
}

// Authenticator represents an authentication strategy
type Authenticator interface {
	Authenticate(context.Context) (interface{}, error)
}

// AuthorizerFunc turns a function into an authorizer
type AuthorizerFunc func(ctx context.Context, principal interface{}, fullMethod string) error

// Authorize authorizes the processing of the request for the principal
func (f AuthorizerFunc) Authorize(ctx context.Context, principal interface{}, fullMethod string) error {
	return f(ctx, principal, fullMethod)
}

// Authorizer represents an authorization strategy
type Authorizer interface {
	Authorize(ctx context.Context, principal interface{}, fullMethod string) error
}
