package cache

import (
	"context"
	"testing"
	"time"
)

func TestBigCache(t *testing.T) {
	bc, _ := NewCacher(&Options{
		h

	})
	bc.Set(context.Background(), "hello", []byte("world"))
	time.Sleep(10 * time.Second)
	v, _ := bc.Get(context.Background(), "hello")
	t.Log(v, string(v))
}
