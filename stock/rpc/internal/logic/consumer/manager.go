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
	mgr.Add(cr)
	mgr.Add(cr2)
	return mgr, nil
}
