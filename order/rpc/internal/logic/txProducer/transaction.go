package txProducer

import (
	"encoding/json"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/order/rpc/internal/svc"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

// NewOrderTransactionChecker 订单服务的事务回查器
// 注意：依赖的是底层资源（db），不是 svcCtx 本身，避免循环依赖
func NewOrderTransactionChecker(svcCtx *svc.ServiceContext) *rmq_client.TransactionChecker {
	return &rmq_client.TransactionChecker{
		Check: func(msg *rmq_client.MessageView) rmq_client.TransactionResolution {
			// 1. 从消息体里解析出业务主键
			var orderMsg mq.OrderTransactionMsg
			err := json.Unmarshal(msg.GetBody(), &orderMsg)
			if err != nil {
				logx.Errorf(constant.UnmarshalErr, "order/rpc/internal/logic/txProducer", err)
			}
			if orderMsg.OrderNo == 0 {
				// 消息格式无法识别，回滚避免脏数据
				return rmq_client.ROLLBACK
			}

			// 2. 查本地事务表，判断这条订单的本地事务是否已提交
			tlog, err := svcCtx.TransactionLog.FindStatusByOrderNo(orderMsg.OrderNo)
			if err != nil {
				// 查不到或查询出错，返回 UNKNOWN，让 Broker 稍后再回查
				return rmq_client.UNKNOWN
			}
			if tlog == nil {
				return rmq_client.UNKNOWN
			}
			// 3. 根据本地事务状态返回结果
			switch tlog.Status {
			case int64(1):
				return rmq_client.COMMIT
			case int64(2):
				return rmq_client.ROLLBACK
			case int64(0):
				//这是第一次插入消息事务表的初始state 0，检测到0就应该更新为1，防止无限回查
				err = svcCtx.TransactionLog.UpdateStatusByOrderNo(orderMsg.OrderNo)
				if err != nil {
					return rmq_client.UNKNOWN
				}
				return rmq_client.COMMIT

			default:
				return rmq_client.UNKNOWN
			}
		},
	}
}
