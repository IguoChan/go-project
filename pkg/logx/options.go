package logx

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

type setOptions struct {
	level  zapcore.Level
	writer io.Writer
}

type Option func(opts *setOptions)

type FileOption struct {
	MaxSize    int  // 在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups int  // 保留旧文件的最大个数
	MaxAges    int  // 保留旧文件的最大天数
	Compress   bool // 是否压缩/归档旧文件
}

func defaultOptions() *setOptions {
	return &setOptions{
		level:  zapcore.DebugLevel,
		writer: os.Stdout,
	}
}

func SetLevel(level zapcore.Level) Option {
	return func(opts *setOptions) {
		if level < zapcore.DebugLevel || level > zapcore.FatalLevel {
			return
		}
		opts.level = level
	}
}

func SetWriteSyncer(file string, fileOpt ...*FileOption) Option {
	return func(opts *setOptions) {
		if file == "" {
			return
		}

		maxSize := 5
		maxBackups := 5
		maxAges := 30
		compress := false
		if len(fileOpt) != 0 {
			maxSize = fileOpt[0].MaxSize
			maxBackups = fileOpt[0].MaxBackups
			maxAges = fileOpt[0].MaxAges
			compress = fileOpt[0].Compress
		}

		opts.writer = &lumberjack.Logger{
			Filename:   file,
			MaxSize:    maxSize,
			MaxAge:     maxAges,
			MaxBackups: maxBackups,
			Compress:   compress,
		}
	}
}
