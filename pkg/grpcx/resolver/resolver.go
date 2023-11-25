package resolver

import (
	"context"

	"google.golang.org/grpc/resolver"
)

type Discovery interface {
	Target() string

	resolver.Builder
}

type Register interface {
	Registry(ctx context.Context, addr string) error
	UnRegistry()
}
