package mq

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/service"
)

type ConsumerManager struct {
	consumers []*Consumer
}

// NewConsumerManager 创建管理器
func NewConsumerManager() *ConsumerManager {
	return &ConsumerManager{
		consumers: make([]*Consumer, 0),
	}
}

// Add 注册一个消费者
func (m *ConsumerManager) Add(cr *Consumer) {
	m.consumers = append(m.consumers, cr)
}

// Start 实现 service.Service 接口，依次启动所有消费者
func (m *ConsumerManager) Start() {
	for _, cr := range m.consumers {
		if err := cr.Start(); err != nil {
			panic(fmt.Sprintf("mq: start consumer[%s] failed: %v", cr.Name(), err))
		}
	}
}

// Stop 实现 service.Service 接口，反向关闭所有消费者
func (m *ConsumerManager) Stop() {
	for i := len(m.consumers) - 1; i >= 0; i-- {
		cr := m.consumers[i]
		if err := cr.Stop(); err != nil {
			fmt.Printf("mq: shutdown consumer[%s] error: %v\n", cr.Name(), err)
		}
	}
}

var _ service.Service = (*ConsumerManager)(nil)
