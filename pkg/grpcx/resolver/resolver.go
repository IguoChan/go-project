package resolver

import (
	"context"

	"google.golang.org/grpc/resolver"
)

type Discovery interface {
	Register()
	Target() string

	resolver.Builder
	resolver.Resolver
}

type Register interface {
	Registry(ctx context.Context, addr string) error
	UnRegistry()
}
