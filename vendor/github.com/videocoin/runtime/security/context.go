package security

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	userKey  struct{}
	tokenKey struct{}
)

func UserFromCtx(ctx context.Context) (string, error) {
	user, ok := ctx.Value(userKey{}).(string)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "Bad user string")
	}
	if user == "" {
		return "", status.Errorf(codes.Unauthenticated, "User unauthenticated with Bearer")
	}

	return user, nil
}
