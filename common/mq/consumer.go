package mq

import (
	"context"
	"fmt"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumeHandler func(ctx context.Context, mvs []*rmq_client.MessageView) error
type Consumer struct {
	c       rmq_client.SimpleConsumer
	name    string
	topic   string
	group   string
	config  ConsumerConfig
	closeCh chan struct{}  //主动关闭信号
	handler ConsumeHandler //业务消费函数
}

func NewConsumer(name string, c ConsumerConfig, handler ConsumeHandler) (*Consumer, error) {
	if c.ConsumerGroup == "" {
		return nil, fmt.Errorf("mq: consumer[%s] ConsumerGroup is required", name)
	}
	if c.Topic == "" {
		return nil, fmt.Errorf("mq: consumer[%s] Topic is required", name)
	}
	if c.Endpoint == "" {
		return nil, fmt.Errorf("mq: consumer[%s] Endpoint is empty", name)
	}
	consumer, err := rmq_client.NewSimpleConsumer(
		&rmq_client.Config{
			Endpoint:      c.Endpoint,
			ConsumerGroup: c.ConsumerGroup,
			Credentials:   nil,
		},
		rmq_client.WithSimpleAwaitDuration(c.AwaitDuration),
		rmq_client.WithSimpleSubscriptionExpressions(map[string]*rmq_client.FilterExpression{
			c.Topic: rmq_client.SUB_ALL,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("mq: new consumer[%s] failed: %w", name, err)
	}

	if err = consumer.Start(); err != nil {
		return nil, fmt.Errorf("mq: start consumer[%s] failed: %w", name, err)
	}

	return &Consumer{
		name:    name,
		c:       consumer,
		topic:   c.Topic,
		group:   c.ConsumerGroup,
		config:  c,
		closeCh: make(chan struct{}),
		handler: handler,
	}, nil
}
func (cr *Consumer) Start() error {
	go func() {
		for {
			select {
			case <-cr.closeCh:
				logx.Info("消费者已经主动关闭，不再拉取消息")
				return
			default:
			}
			mvs, err := cr.c.Receive(context.Background(), cr.config.MaxMsgNum, cr.config.InvisibleDuration)
			if err != nil {
				//暂时没有消息也会进入此分支
				select {
				case <-cr.closeCh:
					logx.Info("消费者已经主动关闭，不再拉取消息")
					return
				default:
				}
				time.Sleep(500 * time.Millisecond)
				continue
			}
			err = cr.handler(context.Background(), mvs)
			if err == nil {
				//消费成功才ack
				for _, mv := range mvs {
					if ackErr := cr.c.Ack(context.Background(), mv); ackErr != nil {
						fmt.Printf("mq: consumer[%s] ack error: %v\n", cr.name, ackErr)
					}
				}
			}
		}
	}()
	return nil
}
func (cr *Consumer) Stop() error {
	if cr.c == nil {
		return nil
	}
	close(cr.closeCh)
	return cr.c.GracefulStop()
}
func (cr *Consumer) Name() string {
	return cr.name
}
