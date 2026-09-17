package consumer

import (
	"zeromall/balance/rpc/internal/svc"
	"zeromall/common/mq"
)

func NewConsumerManager(svc *svc.ServiceContext) (*mq.ConsumerManager, error) {
	mgr := mq.NewConsumerManager()
	cg := mq.ConsumerConfig{
		Endpoint:          svc.Config.RocketMqConf.Endpoint,
		Topic:             svc.Config.RocketMqConf.Topics.TopicInsertBalance,
		ConsumerGroup:     svc.Config.RocketMqConf.Consumer.Group.BalanceInsertGroup,
		AwaitDuration:     svc.Config.RocketMqConf.Consumer.AwaitDuration,
		MaxMsgNum:         svc.Config.RocketMqConf.Consumer.MaxMsgNum,
		InvisibleDuration: svc.Config.RocketMqConf.Consumer.InvisibleDuration,
	}
	balanceInsert := NewBalanceInsertConsumer(svc)
	cr, err := mq.NewConsumer("balance-insert", cg, balanceInsert.Consume)
	if err != nil {
		return nil, err
	}
	mgr.Add(cr)
	return mgr, nil
}
