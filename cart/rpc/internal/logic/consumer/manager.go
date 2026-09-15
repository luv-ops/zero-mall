package consumer

import (
	"zeromall/cart/rpc/internal/svc"
	"zeromall/common/mq"
)

func NewConsumerManager(svc *svc.ServiceContext) (*mq.ConsumerManager, error) {
	mgr := mq.NewConsumerManager()

	cartChangeLogic := NewCartChangeConsumer(svc)
	cg := mq.ConsumerConfig{
		Endpoint:          svc.Config.RocketMqConf.Endpoint,
		Topic:             svc.Config.RocketMqConf.Topics.TopicSyncFiling,
		ConsumerGroup:     svc.Config.RocketMqConf.Consumer.Group.GroupCartFiling,
		AwaitDuration:     svc.Config.RocketMqConf.Consumer.AwaitDuration,
		MaxMsgNum:         svc.Config.RocketMqConf.Consumer.MaxMsgNum,
		InvisibleDuration: svc.Config.RocketMqConf.Consumer.InvisibleDuration,
	}
	cartCr, err := mq.NewConsumer("cart-change", cg, cartChangeLogic.Consume)
	if err != nil {
		return nil, err
	}
	//其他消费者注册都在这里
	cg2 := mq.ConsumerConfig{
		Endpoint:          svc.Config.RocketMqConf.Endpoint,
		Topic:             svc.Config.RocketMqConf.Topics.TopicDelCart,
		ConsumerGroup:     svc.Config.RocketMqConf.Consumer.Group.GroupCartDel,
		AwaitDuration:     svc.Config.RocketMqConf.Consumer.AwaitDuration,
		MaxMsgNum:         svc.Config.RocketMqConf.Consumer.MaxMsgNum,
		InvisibleDuration: svc.Config.RocketMqConf.Consumer.InvisibleDuration,
	}
	cartDelLogic := NewCartDelConsumer(svc)
	cartCr2, err := mq.NewConsumer("cart-del", cg2, cartDelLogic.Consume)
	if err != nil {
		return nil, err
	}
	mgr.Add(cartCr)
	mgr.Add(cartCr2)
	return mgr, nil
}
