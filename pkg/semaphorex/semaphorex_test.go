package semaphorex

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/IguoChan/go-project/pkg/util"

	"github.com/IguoChan/go-project/pkg/cache/redisx"
)

func TestSemaphore(t *testing.T) {
	nameOpt := SetName("semaphore_test1")
	rc, err := redisx.NewClient(&redisx.Options{
		Addrs:       []string{"192.168.0.98:6379"},
		Password:    "123456",
		DialTimeout: 10 * time.Second,
		ReadTimeout: 600 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	rcOpt := SetRedisClient(rc)
	timeoutOpt := SetTimeout(2 * time.Minute)
	sema := NewSemaphore(SemaphoreRedis, 5, nameOpt, rcOpt, timeoutOpt)

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ok := sema.TryAcquire()
			if !ok {
				t.Logf("[%+v] [%d] [%s] get sema failed!", time.Now(), idx, identifyId())
				return
			}
			defer sema.Release()
			t.Logf("[%+v] [%d] [%s] get sema success!", time.Now(), idx, identifyId())
			defer func() {
				t.Logf("[%+v] [%d] [%s] release sema!", time.Now(), idx, identifyId())
			}()
			time.Sleep(3 * time.Second)
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second)
}

func identifyId() string {
	hostname, _ := os.Hostname()
	// 使用 hostname-pid-goroutineId 作为唯一标识
	return fmt.Sprintf("%v-%v-%v", hostname, os.Getpid(), util.GoroutineId())
}
