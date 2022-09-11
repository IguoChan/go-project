package redisx

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/IguoChan/go-project/pkg/logx"
	"github.com/go-redis/redis/v8"
)

type Client struct {
	redis.UniversalClient
}

type Options struct {
	Addrs      []string // len > 1, ClusterClient
	Username   string
	Password   string
	DB         int
	MasterName string // FailoverClient

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// logger
	LogrusLogger *zap.Logger
}

func NewClient(opt *Options) (*Client, error) {
	if opt == nil {
		return nil, errors.New("options is nil")
	}

	// logger
	if opt.LogrusLogger != nil {
		logger := &logx.Logger{Logger: opt.LogrusLogger}
		redis.SetLogger(logger)
	}

	// client
	c := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        opt.Addrs,
		DB:           opt.DB,
		Username:     opt.Username,
		Password:     opt.Password,
		MasterName:   opt.MasterName,
		DialTimeout:  opt.DialTimeout,
		ReadTimeout:  opt.ReadTimeout,
		WriteTimeout: opt.WriteTimeout,
	})

	// ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := c.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &Client{UniversalClient: c}, nil

}

func (c *Client) Close() {
	if c != nil {
		_ = c.UniversalClient.Close()
	}
}

func (c *Client) Set(ctx context.Context, key string, value string) error {
	return c.UniversalClient.Set(ctx, key, string(value), 0).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	stringCmd := c.UniversalClient.Get(ctx, key)
	if stringCmd.Err() != nil {
		return "", stringCmd.Err()
	}
	return stringCmd.String(), nil
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.UniversalClient.Del(ctx, key).Err()
}
