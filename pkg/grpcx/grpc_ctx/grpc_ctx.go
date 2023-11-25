package grpc_ctx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	traceIdKey = "traceid" // grpc的metadata会将大写全部改为小写，即不区分大小写
)

type grpcCtx struct {
	// 内嵌context
	ctx context.Context

	// traceId
	traceIdKey, traceIdVal string
}

func GCtx(ctx context.Context) *grpcCtx {
	gCtx := &grpcCtx{ctx: ctx}

	// 取出traceId
	gCtx.traceIdVal = gCtx.GetTraceId()

	return gCtx
}

func (c *grpcCtx) GetTraceId(traceIdKeys ...string) string {
	if c.traceIdVal != "" {
		return c.traceIdVal
	}

	md, ok := metadata.FromIncomingContext(c.ctx)
	if !ok {
		return ""
	}

	for _, key := range traceIdKeys {
		if v, ok := md[key]; ok {
			c.traceIdKey, c.traceIdVal = key, v[0]
			return v[0]
		}
	}

	if v, ok := md[traceIdKey]; ok {
		c.traceIdKey, c.traceIdVal = traceIdKey, v[0]
		return v[0]
	}

	return ""
}

func (c *grpcCtx) SetTraceId(traceIdVal string, traceIdKeys ...string) {
	if len(traceIdKeys) > 0 {
		c.traceIdKey, c.traceIdVal = traceIdKeys[0], traceIdVal
		return
	}

	c.traceIdKey, c.traceIdVal = traceIdKey, traceIdVal
}
