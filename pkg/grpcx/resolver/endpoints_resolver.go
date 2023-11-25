package resolver

import (
	"context"
	"sync"

	"github.com/IguoChan/go-project/pkg/util"
	"google.golang.org/grpc/resolver"
)

const (
	defaultEndpointsScheme = "endpoints"
)

type EndpointsResolver struct {
	endpoints []string
	scheme    string

	cc         resolver.ClientConn
	serverList sync.Map
}

func NewEndpointsResolver(endpoints []string, opts ...Option) *EndpointsResolver {
	// set options
	defOpts := defaultOptions()
	for _, apply := range opts {
		apply(defOpts)
	}

	e := &EndpointsResolver{
		endpoints: endpoints,
		scheme:    util.SetIfEmpty(defOpts.scheme, defaultEndpointsScheme),
	}
	for _, endpoint := range endpoints {
		e.serverList.Store(endpoint, resolver.Address{Addr: endpoint})
	}

	return e
}

func (e *EndpointsResolver) Register() {
	resolver.Register(e)
}

func (e *EndpointsResolver) Target() string {
	return e.scheme + ":///" + "whatever"
}

func (e *EndpointsResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	e.cc = cc
	addrs := e.getServices()
	_ = e.cc.UpdateState(resolver.State{Addresses: addrs})
	return e, nil
}

func (e *EndpointsResolver) Scheme() string {
	return e.scheme
}

func (e *EndpointsResolver) ResolveNow(nowOptions resolver.ResolveNowOptions) {}

func (e *EndpointsResolver) Close() {}

func (e *EndpointsResolver) Registry(ctx context.Context, addr string) error {
	return nil
}

func (e *EndpointsResolver) UnRegistry() {}

func (e *EndpointsResolver) AddServerList(endpoint string) {
	e.serverList.Store(endpoint, resolver.Address{Addr: endpoint})
	_ = e.cc.UpdateState(resolver.State{Addresses: e.getServices()})
}

func (e *EndpointsResolver) DelServerList(endpoint string) {
	e.serverList.Delete(endpoint)
	_ = e.cc.UpdateState(resolver.State{Addresses: e.getServices()})
}

func (e *EndpointsResolver) getServices() []resolver.Address {
	addrs := make([]resolver.Address, 0, 10)
	e.serverList.Range(func(k, v any) bool {
		addrs = append(addrs, v.(resolver.Address))
		return true
	})
	return addrs
}
