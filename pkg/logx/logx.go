package logx

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(opts ...Option) *Logger {
	// option
	defaultOpts := defaultOptions()
	for _, apply := range opts {
		apply(defaultOpts)
	}

	// new Encoder
	enc := encoder()

	// new core
	core := zapcore.NewCore(enc, zapcore.AddSync(defaultOpts.writer), defaultOpts.level)

	// new Logger
	logger := zap.New(core, zap.AddCaller()) // zap.AddCaller()  添加将调用函数信息记录到日志中的功能。

	return &Logger{
		Logger: logger,
	}
}

func encoder() zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder   // 修改时间编码器
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder // 日志中使用大写字母
	return zapcore.NewConsoleEncoder(cfg)
}

func (l *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	l.Logger.Sugar().Infof(format, v)
}
