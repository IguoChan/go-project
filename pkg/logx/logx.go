package logx

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func (l *Logger) Printf(ctx context.Context, format string, v ...any) {
	l.Logger.Logf(l.Logger.Level, format, v...)
}
