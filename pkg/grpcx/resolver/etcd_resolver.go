package resolver

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/IguoChan/go-project/pkg/grpcx/balancer"

	"github.com/IguoChan/go-project/pkg/etcdx"
	"github.com/IguoChan/go-project/pkg/util"
	"go.etcd.io/etcd/api/v3/mvccpb"
	v3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
)

const (
	defaultGrpcScheme = "grpc"
)

type EtcdResolver struct {
	*etcdx.ClientV3

	// discovery
	cc         resolver.ClientConn
	serverList sync.Map
	cancel     context.CancelFunc

	// register
	prefix, key, val string
	leaseID          v3.LeaseID
	keepAliveChan    <-chan *v3.LeaseKeepAliveResponse

	scheme      string
	serviceName string
}

func NewEtcdResolver(serviceName string, etcdOpts *etcdx.Options, opts ...Option) (*EtcdResolver, error) {
	// set options
	defOpts := defaultOptions()
	for _, apply := range opts {
		apply(defOpts)
	}

	// new etcd client
	c, err := etcdx.NewClient(etcdOpts)
	if err != nil {
		return nil, err
	}

	return &EtcdResolver{
		ClientV3:    c,
		scheme:      util.SetIfEmpty(defOpts.scheme, defaultGrpcScheme),
		serviceName: serviceName,
	}, nil
}

func (e *EtcdResolver) Register() {
	resolver.Register(e)
}

func (e *EtcdResolver) Target() string {
	return fmt.Sprintf("%s:///%s", e.scheme, e.serviceName)
}

func (e *EtcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	e.cc = cc

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	e.cancel = cancel
	e.prefix = "/" + target.URL.Scheme + target.URL.Path + "/"
	resp, err := e.Get(ctx, e.prefix, v3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, ev := range resp.Kvs {
		e.setServerList(string(ev.Key), string(ev.Value))
	}
	_ = e.cc.UpdateState(resolver.State{Addresses: e.getServices()})

	go e.watch(e.prefix)

	return e, nil
}

func (e *EtcdResolver) Scheme() string {
	return e.scheme
}

func (e *EtcdResolver) ResolveNow(nowOptions resolver.ResolveNowOptions) {}

func (e *EtcdResolver) Close() {
	_ = e.ClientV3.Close()
	e.cancel()
}

func (e *EtcdResolver) watch(prefix string) {
	w := e.Watch(context.Background(), prefix, v3.WithPrefix())
	for resp := range w {
		for _, ev := range resp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				e.setServerList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				e.delServerList(string(ev.Kv.Key))
			}
		}
	}
}

func (e *EtcdResolver) setServerList(k, v string) {
	addr := resolver.Address{Addr: strings.TrimPrefix(k, e.prefix)}
	weight, err := strconv.Atoi(v)
	if err != nil {
		weight = balancer.MinWeight
	}
	addr = balancer.SetAddrInfo(addr, balancer.WeightAddrInfo{Weight: weight})
	e.serverList.Store(k, addr)
	_ = e.cc.UpdateState(resolver.State{Addresses: e.getServices()})
}

func (e *EtcdResolver) delServerList(k string) {
	e.serverList.Delete(k)
	_ = e.cc.UpdateState(resolver.State{Addresses: e.getServices()})
}

func (e *EtcdResolver) getServices() []resolver.Address {
	addrs := make([]resolver.Address, 0, 10)
	e.serverList.Range(func(k, v any) bool {
		addrs = append(addrs, v.(resolver.Address))
		return true
	})
	return addrs
}

func (e *EtcdResolver) Registry(ctx context.Context, addr string) error {
	e.key = "/" + e.scheme + "/" + e.serviceName + "/" + addr
	e.val = strconv.Itoa(14)

	// 设置租约
	resp, err := e.Grant(ctx, e.TTL)
	if err != nil {
		return err
	}

	// 注册服务并绑定租约
	_, err = e.Put(ctx, e.key, e.val, v3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	// 注意此处的ctx不可以使用传入的ctx，因为其将作为后台任务一直运行
	keepAlive, err := e.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	if keepAlive == nil {
		return errors.New("keepAlive is nil")
	}

	e.keepAliveChan = keepAlive
	e.leaseID = resp.ID

	go e.listenLeaseRespChan()

	e.Logger.Info("register to etcd success", zap.String("key", e.key), zap.String("val", e.val), zap.Int64("leaseID", int64(e.leaseID)))

	return nil
}

func (e *EtcdResolver) UnRegistry() {
	_, _ = e.Revoke(context.Background(), e.leaseID)
}

func (e *EtcdResolver) listenLeaseRespChan() {
	for _ = range e.keepAliveChan {
		//e.Logger.Info(resp.String())
	}
	e.Logger.Warn("lease closed, cancel or timeout")
}
