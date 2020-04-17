package runtime

import (
	"net/http"
)

// AuthenticatorFunc turns a function into an authenticator
type AuthenticatorFunc func(r *http.Request) (interface{}, error)

func (f AuthenticatorFunc) Authenticate(r *http.Request) (interface{}, error) {
	return f(r)
}

// Authenticator represents an authentication strategy
type Authenticator interface {
	Authenticate(r *http.Request) (interface{}, error)
}

// AuthorizerFunc turns a function into an authorizer
type AuthorizerFunc func(r *http.Request, principal interface{}) error

// Authorize authorizes the processing of the request for the principal
func (f AuthorizerFunc) Authorize(r *http.Request, principal interface{}) error {
	return f(r, principal)
}

// Authorizer represents an authorization strategy
type Authorizer interface {
	Authorize(r *http.Request, principal interface{}) error
}
