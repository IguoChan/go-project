package etcdx

import (
	"context"
	"errors"

	"github.com/IguoChan/go-project/pkg/util"
	v3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Register struct {
	*ClientV3
	key, val      string
	leaseID       v3.LeaseID
	keepAliveChan <-chan *v3.LeaseKeepAliveResponse
}

func NewRegister(opt *Options) (*Register, error) {
	c, err := NewClient(opt)
	if err != nil {
		return nil, err
	}
	c.TTL = util.SetIf0(opt.TTL, defaultTTL)
	c.Scheme = util.SetIfEmpty(opt.Scheme, defaultGrpcScheme)
	return &Register{
		ClientV3: c,
	}, nil
}

func (r *Register) Registry(ctx context.Context, serviceName, addr string) error {
	r.key = "/" + r.Scheme + "/" + serviceName + "/" + addr
	r.val = addr

	// 设置租约
	resp, err := r.Grant(ctx, r.TTL)
	if err != nil {
		return err
	}

	// 注册服务并绑定租约
	_, err = r.Put(ctx, r.key, r.val, v3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	// 注意此处的ctx不可以使用传入的ctx，因为其将作为后台任务一直运行
	keepAlive, err := r.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	if keepAlive == nil {
		return errors.New("keepAlive is nil")
	}

	r.keepAliveChan = keepAlive
	r.leaseID = resp.ID

	go r.listenLeaseRespChan()

	r.Logger.Info("register to etcd success", zap.String("key", r.key), zap.String("val", r.val), zap.Int64("leaseID", int64(r.leaseID)))

	return nil
}

func (r *Register) listenLeaseRespChan() {
	for _ = range r.keepAliveChan {
		//r.Logger.Info(resp.String())
	}
	r.Logger.Warn("lease closed, cancel or timeout")
}

func (r *Register) Revoke() error {
	_, err := r.ClientV3.Revoke(context.Background(), r.leaseID)
	if err != nil {
		return err
	}
	return nil
}
