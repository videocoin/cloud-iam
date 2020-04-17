package security

import (
	"context"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/videocoin/runtime"
)

var defaultAuthOptions = authOptions{}

// AuthOption configures how we set up the authentication and authorization layers.
type AuthOption interface {
	apply(*authOptions)
}

type authOptions struct {
	defaultAuthenticator runtime.Authenticator
	authStrategies       []*authStrategy
}

type funcAuthOption struct {
	f func(*authOptions)
}

func (fdo *funcAuthOption) apply(do *authOptions) {
	fdo.f(do)
}

func newFuncAuthOption(f func(*authOptions)) *funcAuthOption {
	return &funcAuthOption{
		f: f,
	}
}

type authStrategy struct {
	auth runtime.Authenticator
	rule AuthRuleFunc
}

type AuthRuleFunc func(token *jwt.Token) bool

func WithAuthentication(auth runtime.Authenticator, rules ...AuthRuleFunc) AuthOption {
	return newFuncAuthOption(func(o *authOptions) {
		if len(rules) == 0 {
			o.defaultAuthenticator = auth
			return
		}

		o.authStrategies = append(o.authStrategies, &authStrategy{
			auth: auth,
			rule: rules[0],
		})
	})
}

type authnz struct {
	opts authOptions
	h    http.Handler
}

func newAuthnzHandler(h http.Handler, opt ...AuthOption) http.Handler {
	opts := defaultAuthOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	return &authnz{
		opts: opts,
		h:    h,
	}
}

// Satisfies the http.Handler interface for basicAuth.
func (a authnz) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := a.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	/*
		if err := a.rbac.Authorize(ctx, user, fullMethod); err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
	*/

	a.h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey{}, user)))
}

func (a authnz) authenticate(r *http.Request) (interface{}, error) {
	tokenStr, err := authFromReq(r, "Bearer")
	if err != nil {
		return nil, err
	}
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &jwt.StandardClaims{})
	if err != nil {
		return nil, err
	}

	authenticator := a.opts.defaultAuthenticator
	for _, strategy := range a.opts.authStrategies {
		if strategy.rule(token) {
			authenticator = strategy.auth
			break
		}
	}

	return authenticator.Authenticate(r)
}

func Auth(opt ...AuthOption) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return newAuthnzHandler(h, opt...)
	}
}

/*
return &authnz{
	svcacc: ServiceAccount(opts.audience, opts.hmacSecret, opts.pubKeyFunc),
	hmac:   HMACJWT(opts.hmacSecret),
	rbac:   RBAC(),
}

token, err := parseHeader(tokenStr)
if err != nil {
	return nil, err
}

var authenticator runtime.Authenticator
if _, ok := token.Header["kid"]; ok {
	authenticator = a.svcacc
} else {
	authenticator = a.hmac
}

return authenticator.Authenticate(r)
*/
