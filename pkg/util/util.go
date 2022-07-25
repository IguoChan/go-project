package util

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func SetIfNil[T int | int32 | int64 | string | float32 | float64](value *T, set T) T {
	if value != nil {
		return *value
	}
	return set
}

func SetIf0[T int | int32 | int64 | float32 | float64 | time.Duration](value T, set T) T {
	if value != T(0) {
		return value
	}
	return set
}

func SetIfEmpty(value string, set string) string {
	if value != "" {
		return value
	}
	return set
}

// GoroutineId 获取当前协程id
func GoroutineId() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return int64(id)
}
