package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userKey struct{}

func UserFromCtx(ctx context.Context) (string, error) {
	userInfo, ok := ctx.Value(userKey{}).(*UserInfo)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "Bad user information")
	}
	if userInfo.ID == "" {
		return "", status.Errorf(codes.Unauthenticated, "User unauthenticated with Bearer")
	}

	return userInfo.ID, nil
}
