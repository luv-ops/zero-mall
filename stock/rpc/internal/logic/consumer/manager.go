package consumer

import (
	"zeromall/common/mq"
	"zeromall/stock/rpc/internal/svc"
)

func NewConsumerManager(svc *svc.ServiceContext) (*mq.ConsumerManager, error) {
	mgr := mq.NewConsumerManager()
	cg := mq.ConsumerConfig{
		Endpoint:          svc.Config.RocketMqConf.Endpoint,
		Topic:             svc.Config.RocketMqConf.Topics.TopicInsertStock,
		ConsumerGroup:     svc.Config.RocketMqConf.Consumer.Group.StockInsertGroup,
		AwaitDuration:     svc.Config.RocketMqConf.Consumer.AwaitDuration,
		MaxMsgNum:         svc.Config.RocketMqConf.Consumer.MaxMsgNum,
		InvisibleDuration: svc.Config.RocketMqConf.Consumer.InvisibleDuration,
	}
	stockInsert := NewStockInsertConsumer(svc)
	cr, err := mq.NewConsumer("stock-insert", cg, stockInsert.Consume)
	if err != nil {
		return nil, err
	}
	cg2 := mq.ConsumerConfig{
		Endpoint:          svc.Config.RocketMqConf.Endpoint,
		Topic:             svc.Config.RocketMqConf.Topics.TopicReturnStock,
		ConsumerGroup:     svc.Config.RocketMqConf.Consumer.Group.StockReturnGroup,
		AwaitDuration:     svc.Config.RocketMqConf.Consumer.AwaitDuration,
		MaxMsgNum:         svc.Config.RocketMqConf.Consumer.MaxMsgNum,
		InvisibleDuration: svc.Config.RocketMqConf.Consumer.InvisibleDuration,
	}
	stockReturn := NewStockReturnConsumer(svc)
	cr2, err := mq.NewConsumer("stock-return", cg2, stockReturn.Consume)
	if err != nil {
		return nil, err
	}
	cg3 := mq.ConsumerConfig{
		Endpoint:          svc.Config.RocketMqConf.Endpoint,
		Topic:             svc.Config.RocketMqConf.Topics.TopicTxFrozenStock,
		ConsumerGroup:     svc.Config.RocketMqConf.Consumer.Group.StockFrozenGroup,
		AwaitDuration:     svc.Config.RocketMqConf.Consumer.AwaitDuration,
		MaxMsgNum:         svc.Config.RocketMqConf.Consumer.MaxMsgNum,
		InvisibleDuration: svc.Config.RocketMqConf.Consumer.InvisibleDuration,
	}
	stockFrozen := NewStockFrozenConsumer(svc)
	cr3, err := mq.NewConsumer("stock-return", cg3, stockFrozen.Consume)
	if err != nil {
		return nil, err
	}
	cg4 := mq.ConsumerConfig{
		Endpoint:          svc.Config.RocketMqConf.Endpoint,
		Topic:             svc.Config.RocketMqConf.Topics.TopicDeductStock,
		ConsumerGroup:     svc.Config.RocketMqConf.Consumer.Group.StockDeductGroup,
		AwaitDuration:     svc.Config.RocketMqConf.Consumer.AwaitDuration,
		MaxMsgNum:         svc.Config.RocketMqConf.Consumer.MaxMsgNum,
		InvisibleDuration: svc.Config.RocketMqConf.Consumer.InvisibleDuration,
	}
	stockDeduct := NewStockDeductConsumer(svc)
	cr4, err := mq.NewConsumer("stock-deduct", cg4, stockDeduct.Consume)
	mgr.Add(cr)
	mgr.Add(cr2)
	mgr.Add(cr3)
	mgr.Add(cr4)
	return mgr, nil
}
