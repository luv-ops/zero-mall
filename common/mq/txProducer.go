package mq

import (
	"context"
	"fmt"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
)

type TxProducer struct {
	p rmq_client.Producer
}

func NewTxProducer(c *ProducerConfig, checker *rmq_client.TransactionChecker) (*TxProducer, error) {
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
		rmq_client.WithTransactionChecker(checker), //比普通生产者多一个这个
		rmq_client.WithTopics(c.Topics...),
	)
	if err != nil {
		return nil, fmt.Errorf("mq: new producer failed: %w", err)
	}

	if err = producer.Start(); err != nil {
		return nil, fmt.Errorf("mq: start producer failed: %w", err)
	}
	return &TxProducer{producer}, nil

}

func (t *TxProducer) SendWithTransaction(ctx context.Context, topic string, body []byte, transaction rmq_client.Transaction) (*rmq_client.SendReceipt, error) {
	msg := &rmq_client.Message{
		Topic: topic,
		Body:  body,
	}
	resp, err := t.p.SendWithTransaction(ctx, msg, transaction)
	if err != nil {
		return nil, fmt.Errorf("mq: send tx topic=%s failed: %w", topic, err)
	}
	return resp[0], nil
}
func (t *TxProducer) BeginTransaction() rmq_client.Transaction {
	return t.p.BeginTransaction()
}

// GracefulStop 优雅关闭
func (t *TxProducer) Stop() error {
	if t.p != nil {
		return t.p.GracefulStop()
	}
	return nil
}
