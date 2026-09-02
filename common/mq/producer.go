package mq

import (
	"context"
	"log"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
)

type Producer struct {
	producer rmq_client.Producer
}

func NewProducer(endPoint string, topic string) (*Producer, error) {
	config := &rmq_client.Config{
		Endpoint:    endPoint,
		Credentials: nil,
	}
	//必须加withTopics防止死锁
	pro, err := rmq_client.NewProducer(config, rmq_client.WithTopics(topic))
	if err != nil {
		log.Printf("fail to create producer: %v", err)
		return nil, err
	}
	if err = pro.Start(); err != nil {
		log.Printf("fail to start producer: %v", err)
		_ = pro.GracefulStop()
		return nil, err
	}
	return &Producer{
		producer: pro,
	}, err
}

func (p *Producer) Send(ctx context.Context, topic string, message []byte) error {
	//组装message
	msg := rmq_client.Message{
		Topic: topic,
		Body:  message,
	}
	_, err := p.producer.Send(ctx, &msg)
	if err != nil {
		return err
	}
	return nil
}
func (p *Producer) SendDelay(ctx context.Context, topic string, message []byte, delay time.Duration) error {
	msg := rmq_client.Message{
		Topic: topic,
		Body:  message,
	}
	//设置消息延时
	msg.SetDelayTimestamp(time.Now().Add(delay))
	_, err := p.producer.Send(ctx, &msg)
	if err != nil {
		return err
	}
	return nil
}
func (p *Producer) Stop() error {
	return p.producer.GracefulStop()
}
