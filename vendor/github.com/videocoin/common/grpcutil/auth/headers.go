package auth

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
)

const (
	headerAuthority = ":authority"
	headerPath      = ":path"
	bearer          = "bearer"
)

// BearerFromMD is a helper function to extract the bearer token from the gRPC metadata of the request.
func BearerFromMD(ctx context.Context) (string, error) { return grpc_auth.AuthFromMD(ctx, bearer) }

// AuthorityFromMD is a helper function for extracting the :authority header from the gRPC metadata of the request.
func AuthorityFromMD(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(headerAuthority)
}

// PathFromMD is a helper function for extracting the :path header from the gRPC metadata of the request.
func PathFromMD(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(headerPath)
}
