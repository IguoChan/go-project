package mqx

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/IguoChan/go-project/pkg/util"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	consumer2 "github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	producer2 "github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

const CreateTopicMsg = "CTopic"

type rocketMQ struct {
	producer  rocketmq.Producer
	consumers map[string]rocketmq.PushConsumer
	mu        *sync.Mutex
	rlog.Logger
	retryTimes int
	admin      admin.Admin

	*Options
}

func NewRocketMQx(opt *Options) (*rocketMQ, error) {
	// logger
	//var logger rlog.Logger
	//if opt.LogrusLogger != nil {
	//	logger = &util.LogrusLogger{Logger: opt.LogrusLogger}
	//}
	//rlog.SetLogger(logger)
	//rlog.SetLogMode(opt.LogMode)

	// client
	producer, err := rocketmq.NewProducer(
		producer2.WithNameServer(opt.Endpoints),
		producer2.WithCredentials(primitive.Credentials{
			AccessKey: opt.Username,
			SecretKey: opt.Password,
		}),
		producer2.WithNamespace(opt.Namespace),
		producer2.WithInstanceName(opt.InstanceName),
		producer2.WithGroupName(opt.GroupName),
		producer2.WithRetry(opt.ErrRetryTimes),
	)
	if err != nil {
		return nil, err
	}

	// start producer
	err = producer.Start()
	if err != nil {
		return nil, err
	}

	//Admin
	admin, _ := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(opt.Endpoints)))

	return &rocketMQ{
		producer:   producer,
		consumers:  map[string]rocketmq.PushConsumer{},
		mu:         &sync.Mutex{},
		retryTimes: util.SetIf0(opt.ErrRetryTimes, 3),
		Options:    opt,
		admin:      admin,
	}, nil
}

func (r *rocketMQ) Subscribe(topic, tag, groupId string) (<-chan *Msg, error) {
	// get consumer
	consumer, err := r.getPushConsumer(groupId, topic)
	if err != nil {
		return nil, err
	}

	// subscribe
	selector := consumer2.MessageSelector{}
	if tag != "" {
		selector.Type = consumer2.TAG
		selector.Expression = tag
	}

	ch := make(chan *Msg)
	if err = consumer.Subscribe(topic, selector, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer2.ConsumeResult, error) {
		for _, msg := range ext {
			ch <- &Msg{
				Topic:   topic,
				MsgId:   msg.MsgId,
				MsgBody: msg.Body,
			}
		}
		return consumer2.ConsumeSuccess, nil
	}); err != nil {
		r.Logger.Error(fmt.Sprintf("[RocketMQ] subscribe err: %+v", err), map[string]any{
			rlog.LogKeyConsumerGroup: groupId,
			rlog.LogKeyTopic:         topic,
		})
		return nil, err
	}

	// start
	if err = consumer.Start(); err != nil {
		return nil, err
	}

	return ch, nil
}

func (r *rocketMQ) Unsubscribe(topic, tag, groupId string) error {
	// unsubscribe
	consumer, err := r.getPushConsumer(groupId, topic)
	if err != nil {
		return err
	}

	if err = consumer.Unsubscribe(topic); err != nil {
		r.Logger.Error(fmt.Sprintf("[RocketMQ] unsubscribe err: %+v", err), map[string]any{
			rlog.LogKeyConsumerGroup: groupId,
			rlog.LogKeyTopic:         topic,
		})
		return err
	}

	// shutdown
	_ = consumer.Shutdown()

	// delete consumer
	r.deleteConsumer(groupId)

	return nil
}

func (r *rocketMQ) Publish(topic, tag string, msg []byte) error {
	// new txMsg
	txMsg := primitive.NewMessage(topic, msg)
	if tag != "" {
		txMsg.WithTag(tag)
	}

	// retry, we have set producer2.WithRetry, we doubt is this necessary
	currentTimes := 0 // 当前重试次数
retry:
	result, err := r.producer.SendSync(context.Background(), txMsg)
	if err != nil {
		// 重试错误类型
		if strings.Contains(err.Error(), "[TIMEOUT_CLEAN_QUEUE]broker busy") {
			if currentTimes < r.retryTimes {
				currentTimes++
				r.Logger.Warning(fmt.Sprintf("[RocketMQ] send err:%+v, retry send %d times", err, currentTimes), map[string]any{
					"retry": "failed",
					"err":   err.Error(),
				})
				goto retry
			}
		}

		return err
	}

	r.Logger.Debug(fmt.Sprintf("[RocketMQ] SendMessageSync %s result %s", string(msg), result.String()), map[string]any{
		"send":           "success",
		rlog.LogKeyTopic: topic,
	})

	return nil
}

func (r *rocketMQ) CreateTopic(topic, tag string) error {
	return r.Publish(topic, tag, []byte(CreateTopicMsg))
}

func (r *rocketMQ) DeleteTopic(topic, tag string) error {
	topic = topic + tag
	return r.admin.DeleteTopic(
		context.Background(),
		admin.WithTopicDelete(topic),
	)
}

func (r *rocketMQ) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// shutdown producer
	_ = r.producer.Shutdown()

	// clear consumers
	for groupId, consumer := range r.consumers {
		_ = consumer.Shutdown()
		delete(r.consumers, groupId)
	}

	return nil
}

func (r *rocketMQ) getPushConsumer(groupId, topic string) (rocketmq.PushConsumer, error) {
	// lock for singleton pattern
	r.mu.Lock()
	defer r.mu.Unlock()

	var err error
	consumer, ok := r.consumers[groupId]
	if !ok {
		consumer, err = rocketmq.NewPushConsumer(
			consumer2.WithNameServer(r.Endpoints),
			consumer2.WithCredentials(primitive.Credentials{
				AccessKey: r.Username,
				SecretKey: r.Password,
			}),
			consumer2.WithNamespace(r.Namespace),
			consumer2.WithGroupName(groupId),
			consumer2.WithInstance(r.InstanceName+"consumer"), // 为了应对 https://github.com/apache/rocketmq-client-go/issues/699
			consumer2.WithConsumerModel(consumer2.Clustering),
		)
		if err != nil {
			r.Logger.Error(fmt.Sprintf("[RocketMQ] new pushConsumer err: %+v", err), map[string]any{
				rlog.LogKeyConsumerGroup: groupId,
				rlog.LogKeyTopic:         topic,
			})
			return nil, err
		}

		r.consumers[groupId] = consumer
	}

	return consumer, nil
}

func (r *rocketMQ) deleteConsumer(groupId string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.consumers[groupId]
	if ok {
		delete(r.consumers, groupId)
	}
}
