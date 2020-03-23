// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package auth

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc"
)

// ServiceAuthFuncOverride allows a given gRPC service implementation to override the global `AuthFunc`.
//
// If a service implements the AuthFuncOverride method, it takes precedence over the `AuthFunc` method,
// and will be called instead of AuthFunc for all method invocations within that service.
type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryServerInterceptor(authNZFunc AuthNZFunc, opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		if o.shouldAuth(info.FullMethod) {
			var err error
			if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
				newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
			} else {
				newCtx, err = authNZFunc(ctx, info.FullMethod)
			}
			if err != nil {
				return nil, err
			}
		} else {
			newCtx = ctx
		}

		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new unary server interceptors that performs per-request auth.
func StreamServerInterceptor(authNZFunc AuthNZFunc, opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context

		if o.shouldAuth(info.FullMethod) {
			var err error
			if overrideSrv, ok := srv.(ServiceAuthFuncOverride); ok {
				newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
			} else {
				newCtx, err = authNZFunc(stream.Context(), info.FullMethod)
			}
			if err != nil {
				return err
			}
		} else {
			newCtx = stream.Context()
		}

		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
