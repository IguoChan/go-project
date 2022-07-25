// 参考：https://github.com/rfyiamcool/go_redis_semaphore/blob/433ae39b137d5b15a0f74a5672ffe940b31e3b6f/go_redis_semaphore.go#L114

package semaphorex

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/IguoChan/go-project/pkg/cache/redisx"
	"github.com/IguoChan/go-project/pkg/util"

	"github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
)

type RedisSem struct {
	permit int64
	rc     *redisx.Client

	name     string
	incrKey  string // 用作保证信号量获取的绝对顺序
	ownerKey string // 有序列表
	timeKey  string // 有序列表，做超时判断
	waitKey  string // list，用作队列判断

	timeout int64 // 单位纳秒
}

func NewRedisSem(n int64, semName string, rc *redisx.Client) *RedisSem {
	return &RedisSem{
		permit:   n,
		rc:       rc,
		name:     semName,
		incrKey:  semName + "_incr",
		ownerKey: semName + "_owner",
		timeKey:  semName + "_time",
		waitKey:  semName + "_wait",
		timeout:  600000000000, // 暂定10min
	}
}

func (r *RedisSem) Acquire(ctx context.Context) error {
	ok := r.TryAcquire()
	if !ok {
		cmd := r.rc.BRPop(ctx, 0, r.waitKey)
		if cmd.Err() != nil {
			logrus.Errorf("RedisSem Acquire BRPop failed: %+v", cmd.Err())
			return ErrGetSem
		}
		return r.Acquire(ctx)
	}
	return nil
}

func (r *RedisSem) TryAcquire() bool {
	//// 首先清除超时信号量
	//// 1. 清除时间戳 zset 的超时数据
	//r.rc.ZRemRangeByScore(context.Background(), r.timeKey, "0", strconv.FormatInt(time.Now().UnixNano()-r.timeout, 10))
	//// 2. 取交集，取最小值存入ownKey（一般而言最小值肯定是cnt），使得ownKey也过滤掉超时信号量
	//r.rc.ZInterStore(context.Background(), r.ownerKey, &redis.ZStore{
	//	Keys:      []string{r.timeKey, r.ownerKey},
	//	Aggregate: "MIN",
	//})
	//
	//// 设置
	//intCmd := r.rc.ZAdd(context.Background(), r.timeKey, &redis.Z{
	//	Score:  float64(time.Now().UnixNano()),
	//	Member: r.identifyId(),
	//})
	//if intCmd.Err() != nil {
	//	return false
	//}
	//// 首先获取自增计数
	//intCmd = r.rc.Incr(context.Background(), r.incrKey)
	//if intCmd.Err() != nil {
	//	logrus.Errorf("RedisSem TryAcquire Incr failed: %+v", intCmd.Err())
	//	return false
	//}
	//cnt := intCmd.Val()
	//intCmd = r.rc.ZAdd(context.Background(), r.ownerKey, &redis.Z{
	//	Score:  float64(cnt),
	//	Member: r.identifyId(),
	//})
	//if intCmd.Err() != nil {
	//	return false
	//}
	//
	//// 判断在有序列表中的位置
	//intCmd = r.rc.ZRank(context.Background(), r.ownerKey, r.identifyId())
	//if intCmd.Err() != nil || intCmd.Val() >= r.permit {
	//	r.Release()
	//	return false
	//}
	//
	//return true

	keys := []string{r.timeKey, r.ownerKey, r.waitKey, r.incrKey}
	values := []interface{}{r.identifyId(), time.Now().UnixNano(), time.Now().UnixNano() - r.timeout, r.permit}
	res, err := acquireScript.Run(context.Background(), r.rc, keys, values...).Int()
	if err != nil || res == 0 {
		r.Release()
		return false
	}

	return true
}

func (r *RedisSem) Release() {
	keys := []string{r.timeKey, r.ownerKey, r.waitKey}
	values := []interface{}{r.identifyId()}
	_, err := releaseScript.Run(context.Background(), r.rc, keys, values...).Int()
	if err != nil {
		r.rc.ZRem(context.Background(), r.timeKey, r.identifyId())
		r.rc.ZRem(context.Background(), r.ownerKey, r.identifyId())
	}
}

func (r *RedisSem) identifyId() string {
	hostname, _ := os.Hostname()
	// 使用 hostname-pid-goroutineId 作为唯一标识
	return fmt.Sprintf("%v-%v-%v", hostname, os.Getpid(), util.GoroutineId())
}

var acquireScript = redis.NewScript(`
	redis.call("ZREMRANGEBYSCORE", KEYS[1], 0, ARGV[3])
	redis.call("ZINTERSTORE", KEYS[2], 2, KEYS[1], KEYS[2], "AGGREGATE", "MIN")
	local cnt = redis.call("INCR", KEYS[4])
	redis.call("ZADD", KEYS[1], ARGV[2], ARGV[1])
	redis.call("ZADD", KEYS[2], cnt, ARGV[1])
	print(ARGV[4])
	local res = redis.call("ZRANK", KEYS[2], ARGV[1])
	if res >= tonumber(ARGV[4]) then
		return 0
	else
		return 1
	end
`)

var releaseScript = redis.NewScript(`
	redis.call("ZREM", KEYS[1], ARGV[1])
	redis.call("ZREM", KEYS[2], ARGV[1])
	redis.call("RPUSH", KEYS[3], 1)
	redis.call("EXPIRE", KEYS[3], 5)
	return 1
`)
