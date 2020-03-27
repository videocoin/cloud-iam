package grpcutil

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

const (
	Bearer = "bearer"
)

// BearerFromMD is a helper function to extract the bearer token from the gRPC metadata of the request.
func BearerFromMD(ctx context.Context) (string, error) { return grpc_auth.AuthFromMD(ctx, Bearer) }
