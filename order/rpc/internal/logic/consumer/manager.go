package consumer

import (
	"zeromall/common/mq"
	"zeromall/order/rpc/internal/svc"
)

func NewConsumerManager(svc *svc.ServiceContext) (*mq.ConsumerManager, error) {
	mgr := mq.NewConsumerManager()

	orderOffLogic := NewOrderOffConsumer(svc)
	cg := mq.ConsumerConfig{
		Endpoint:          svc.Config.RocketMqConf.Endpoint,
		Topic:             svc.Config.RocketMqConf.Topics.TopicDelayOrderOff,
		ConsumerGroup:     svc.Config.RocketMqConf.Consumer.Group.TopicDelayOffGroup,
		AwaitDuration:     svc.Config.RocketMqConf.Consumer.AwaitDuration,
		MaxMsgNum:         svc.Config.RocketMqConf.Consumer.MaxMsgNum,
		InvisibleDuration: svc.Config.RocketMqConf.Consumer.InvisibleDuration,
	}
	orderCr, err := mq.NewConsumer("order-off", cg, orderOffLogic.Consume)
	if err != nil {
		return nil, err
	}
	mgr.Add(orderCr)

	return mgr, nil
}
