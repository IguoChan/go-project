package logx

import (
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {

	log := NewLogger(SetWriteSyncer("/Users/chenyiguo/workspace/log/go-project/current.log"), SetLevel(zap.ErrorLevel))
	log.Debug("this is debug message")
	log.Info("this is info message")
	log.Info("this is info message with fileds",
		zap.Int("age", 24), zap.String("agender", "man"))
	log.Warn("this is warn message")
	log.Error("this is error message")
	log.DPanic("This is a DPANIC message")
	// log.Panic("this is panic message")
	// log.Fatal("This is a FATAL message")

	for i := 0; i < 1000000; i++ {
		log.Info("this is info message")
	}
}
