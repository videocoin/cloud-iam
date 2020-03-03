package auth

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

const subjectKey = "sub"

var (
	// ErrMissingRequiredMetadata indicates that authentication metadata is missing.
	ErrMissingRequiredMetadata = errors.New("invalid metadata")
	// ErrMissingRequiredSubject indicates that the authentication subject is missing.
	ErrMissingRequiredSubject = errors.New("missing required subject")
)

// NewOutgoingContext creates a new context with authentication subject attached.
func NewOutgoingContext(ctx context.Context, subject string) context.Context {
	md := make(metadata.MD)
	md[subjectKey] = []string{subject}
	return metadata.NewOutgoingContext(ctx, md)
}

// FromIncomingContext returns the incoming authentication subject in ctx if it exists.
func FromIncomingContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMissingRequiredMetadata
	}

	values := md.Get(subjectKey)
	if len(values) != 1 {
		return "", ErrMissingRequiredSubject
	}

	return values[0], nil
}
