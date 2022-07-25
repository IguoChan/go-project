package lockx

import "errors"

var (
	ErrLockFailed = errors.New("get lock failed")
)

type Locker interface {
	Lock() error
	Unlock() error
	RLock() error
	RUnlock() error
	TryLock() bool
}
