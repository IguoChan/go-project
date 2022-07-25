package bigcachex

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

type Options struct {
	TTL time.Duration
}

type Client struct {
	*bigcache.BigCache
}

func NewClient(opt *Options) (*Client, error) {
	bc, _ := bigcache.NewBigCache(bigcache.DefaultConfig(opt.TTL))
	return &Client{
		BigCache: bc,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, value string) error {
	return c.BigCache.Set(key, []byte(value))
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	v, err := c.BigCache.Get(key)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.BigCache.Delete(key)
}
