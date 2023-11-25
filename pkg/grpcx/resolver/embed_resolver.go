package resolver

import (
	"context"
	"fmt"

	"google.golang.org/grpc/resolver"
)

type (
	Scheme string
)

func (s Scheme) String() string {
	return string(s)
}

const (
	SchemePassThrough Scheme = "passthrough"
	SchemeDNS         Scheme = "dns"
)

// EmbedResolver 用于grpc模块自带的passthrough和dns resolver的封装
type EmbedResolver struct {
	builder resolver.Builder
	scheme  Scheme
	target  string
}

func NewEmbedResolver(scheme Scheme, target string) *EmbedResolver {
	return &EmbedResolver{
		builder: resolver.Get(scheme.String()),
		scheme:  scheme,
		target:  target,
	}
}

func (d *EmbedResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return d.builder.Build(target, cc, opts)
}

func (d *EmbedResolver) Scheme() string {
	return d.builder.Scheme()
}

func (d *EmbedResolver) Register() {}

func (d *EmbedResolver) Target() string {
	return fmt.Sprintf("%s:///%s", d.builder.Scheme(), d.target)
}

func (d *EmbedResolver) Registry(ctx context.Context, addr string) error {
	return nil
}

func (d *EmbedResolver) UnRegistry() {}
