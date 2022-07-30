package semaphorex

import (
	"sync"
	"testing"
	"time"

	"github.com/IguoChan/go-project/pkg/cache/redisx"
)

func TestSemaphore(t *testing.T) {
	nameOpt := SetName("semaphore_test4")
	rc, err := redisx.NewClient(&redisx.Options{
		Addrs:       []string{"192.168.0.102:6379"},
		Password:    "123456",
		DialTimeout: 10 * time.Second,
		ReadTimeout: 60 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	rcOpt := SetRedisClient(rc)
	timeoutOpt := SetTimeout(10 * time.Second)
	sema := NewSemaphore(SemaphoreRedis, 5, nameOpt, rcOpt, timeoutOpt)

	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ok := sema.TryAcquire()
			if !ok {
				//t.Logf("[%+v] [%d] get sema failed!", time.Now(), idx)
				return
			}
			defer sema.Release()
			t.Logf("[%+v] [%d] get sema success!", time.Now(), idx)
			defer func() {
				t.Logf("[%+v] [%d] release sema!", time.Now(), idx)
			}()
			time.Sleep(15 * time.Second)
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second)
}
