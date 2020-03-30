package security

import "context"

type (
	userKey  struct{}
	tokenKey struct{}
)

func UserFromCtx(ctx context.Context) string {
	return ctx.Value(userKey{}).(string)
}
