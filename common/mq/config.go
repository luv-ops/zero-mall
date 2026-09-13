package mq

import "time"

type ProducerConfig struct {
	Endpoint string   `json:"endpoint"`
	Topics   []string `json:"topics"`
}

type ConsumerConfig struct {
	Endpoint          string `json:"endpoint"`
	Topic             string `json:"topic"`
	ConsumerGroup     string
	AwaitDuration     time.Duration
	MaxMsgNum         int32
	InvisibleDuration time.Duration
}
