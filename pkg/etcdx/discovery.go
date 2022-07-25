package etcdx

import (
	"context"
	"sync"
	"time"

	"github.com/IguoChan/go-project/pkg/util"
	"go.etcd.io/etcd/api/v3/mvccpb"
	v3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

type Discovery struct {
	*ClientV3
	conn       resolver.ClientConn
	serverList sync.Map
}

func NewDiscovery(opt *Options) (*Discovery, error) {
	c, err := NewClient(opt)
	if err != nil {
		return nil, err
	}
	c.Scheme = util.SetIfEmpty(opt.Scheme, defaultGrpcScheme)

	return &Discovery{
		ClientV3: c,
	}, nil
}

func (d *Discovery) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	d.conn = cc

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	prefix := "/" + target.URL.Scheme + "/" + target.URL.Host + target.URL.Path
	resp, err := d.Get(ctx, prefix, v3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, e := range resp.Kvs {
		d.setServerList(string(e.Key), string(e.Value))
	}
	_ = d.conn.UpdateState(resolver.State{Addresses: d.getServices()})

	go d.watch(prefix)

	return d, nil
}

func (d *Discovery) Scheme() string {
	return d.ClientV3.Scheme
}

// ResolveNow 监视目标更新
func (d *Discovery) ResolveNow(rn resolver.ResolveNowOptions) {}

//Scheme return schema
func (d *Discovery) Close() {
	_ = d.ClientV3.Close()
}

func (d *Discovery) watch(prefix string) {
	w := d.Watch(context.Background(), prefix, v3.WithPrefix())
	for resp := range w {
		for _, e := range resp.Events {
			switch e.Type {
			case mvccpb.PUT:
				d.setServerList(string(e.Kv.Key), string(e.Kv.Value))
			case mvccpb.DELETE:
				d.delServerList(string(e.Kv.Key))
			}
		}
	}
}

func (d *Discovery) setServerList(k, v string) {
	d.serverList.Store(k, resolver.Address{Addr: v})
	_ = d.conn.UpdateState(resolver.State{Addresses: d.getServices()})
}

func (d *Discovery) delServerList(k string) {
	d.serverList.Delete(k)
	_ = d.conn.UpdateState(resolver.State{Addresses: d.getServices()})
}

func (d *Discovery) getServices() []resolver.Address {
	addrs := make([]resolver.Address, 0, 10)
	d.serverList.Range(func(k, v any) bool {
		addrs = append(addrs, v.(resolver.Address))
		return true
	})
	return addrs
}
