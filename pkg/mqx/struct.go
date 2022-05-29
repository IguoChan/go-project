package mqx

import "github.com/sirupsen/logrus"

type Type string

type Msg struct {
	Topic   string // 来源topic
	MsgId   string // 消息ID，标记消息的唯一id
	MsgBody []byte // 消息内容，byte
}

type Options struct {
	MQType Type

	Endpoints    []string
	Username     string
	Password     string
	Namespace    string
	InstanceName string
	GroupName    string

	// logger
	LogrusLogger *logrus.Logger
	LogMode      string

	// retry times
	ErrRetryTimes int
}
