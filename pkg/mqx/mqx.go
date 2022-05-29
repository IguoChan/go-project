package mqx

type MQer interface {
	MQInner
	MQOuter
	CreateTopic(topic, tag string) error
	DeleteTopic(topic, tag string) error
	Close() error
}

type MQInner interface {
	Subscribe(topic, tag, groupId string) (<-chan *Msg, error)
	Unsubscribe(topic, tag, groupId string) error
}

type MQOuter interface {
	Publish(topic, tag string, msg []byte) error
}

func NewClient(opt *Options) (MQer, error) {
	return NewRocketMQx(opt)
}
