package txProducer

import (
	"context"
	"encoding/json"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/pay/rpc/internal/svc"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

func NewPayTransactionChecker(svcCtx *svc.ServiceContext) *rmq_client.TransactionChecker {
	return &rmq_client.TransactionChecker{
		Check: func(msg *rmq_client.MessageView) rmq_client.TransactionResolution {
			// 1. 从消息体里解析出业务主键
			var payMsg mq.PaySuccessMsg
			err := json.Unmarshal(msg.GetBody(), &payMsg)
			if err != nil {
				logx.Errorf(constant.UnmarshalErr, "order/rpc/internal/logic/txProducer", err)
			}
			if payMsg.OrderNo == 0 {
				// 消息格式无法识别，回滚避免脏数据
				return rmq_client.ROLLBACK
			}

			res, err := svcCtx.PayModel.FindOneByPayNo(context.TODO(), payMsg.PayNo)
			if err != nil {
				// 查不到或查询出错，返回 UNKNOWN，让 Broker 稍后再回查
				return rmq_client.UNKNOWN
			}
			if res != nil {
				return rmq_client.COMMIT
			}
			return rmq_client.ROLLBACK
		},
	}
}
