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
		ConsumerGroup:     svc.Config.RocketMqConf.Consumer.Group.DelayOffGroup,
		AwaitDuration:     svc.Config.RocketMqConf.Consumer.AwaitDuration,
		MaxMsgNum:         svc.Config.RocketMqConf.Consumer.MaxMsgNum,
		InvisibleDuration: svc.Config.RocketMqConf.Consumer.InvisibleDuration,
	}
	orderCr, err := mq.NewConsumer("order-off", cg, orderOffLogic.Consume)
	if err != nil {
		return nil, err
	}
	cg2 := mq.ConsumerConfig{
		Endpoint:          svc.Config.RocketMqConf.Endpoint,
		Topic:             svc.Config.RocketMqConf.Topics.TopicPaySuccess,
		ConsumerGroup:     svc.Config.RocketMqConf.Consumer.Group.PaySuccessGroup,
		AwaitDuration:     svc.Config.RocketMqConf.Consumer.AwaitDuration,
		MaxMsgNum:         svc.Config.RocketMqConf.Consumer.MaxMsgNum,
		InvisibleDuration: svc.Config.RocketMqConf.Consumer.InvisibleDuration,
	}
	paySuccessLogic := NewPaySuccessConsumer(svc)
	payCr, err := mq.NewConsumer("pay-success", cg2, paySuccessLogic.Consume)
	if err != nil {
		return nil, err
	}
	mgr.Add(orderCr)
	mgr.Add(payCr)
	return mgr, nil
}
