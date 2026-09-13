package mq

import (
	"context"
	"fmt"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
)

type Producer struct {
	p rmq_client.Producer
}

func NewProducer(c *ProducerConfig) (*Producer, error) {
	if c.Endpoint == "" {
		return nil, fmt.Errorf("mq: producer Endpoint is empty")
	}
	if len(c.Topics) == 0 {
		return nil, fmt.Errorf("mq: producer Topics is empty")
	}
	producer, err := rmq_client.NewProducer(
		&rmq_client.Config{
			Endpoint:    c.Endpoint,
			Credentials: nil,
		},
		rmq_client.WithTopics(c.Topics...),
	)
	if err != nil {
		return nil, fmt.Errorf("mq: new producer failed: %w", err)
	}

	if err = producer.Start(); err != nil {
		return nil, fmt.Errorf("mq: start producer failed: %w", err)
	}
	return &Producer{producer}, nil

}
func (pr *Producer) Send(ctx context.Context, topic string, body []byte) error {
	msg := &rmq_client.Message{
		Topic: topic,
		Body:  body,
	}

	resp, err := pr.p.Send(ctx, msg)
	if err != nil {
		return fmt.Errorf("mq: send topic=%s failed: %w", topic, err)
	}
	if len(resp) == 0 {
		return fmt.Errorf("mq: send topic=%s returned empty response", topic)
	}
	return nil
}
func (pr *Producer) SendDelay(ctx context.Context, topic string, body []byte, duration time.Duration) error {
	msg := &rmq_client.Message{
		Topic: topic,
		Body:  body,
	}
	msg.SetDelayTimestamp(time.Now().Add(duration))
	resp, err := pr.p.Send(ctx, msg)
	if err != nil {
		return fmt.Errorf("mq: send topic=%s failed: %w", topic, err)
	}
	if len(resp) == 0 {
		return fmt.Errorf("mq: send topic=%s returned empty response", topic)
	}
	return nil
}
func (pr *Producer) Stop() error {
	if pr.p == nil {
		return nil
	}
	return pr.p.GracefulStop()
}
