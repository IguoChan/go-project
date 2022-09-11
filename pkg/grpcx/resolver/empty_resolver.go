package resolver

import (
	"context"

	"google.golang.org/grpc/resolver/manual"
)

type EmptyResolver struct {
	target string

	*manual.Resolver // do nothing
}

func NewEmptyResolver(target string) *EmptyResolver {
	return &EmptyResolver{
		target:   target,
		Resolver: manual.NewBuilderWithScheme("whatever"),
	}
}

// Register EmptyResolver没有注册到resolver，所以会走 defaultResolver 的解析，即当作{{ip:port}}
func (e *EmptyResolver) Register() {}

func (e *EmptyResolver) Target() string {
	return e.target
}

func (e *EmptyResolver) Registry(ctx context.Context, addr string) error {
	return nil
}

func (e *EmptyResolver) UnRegistry() {}
