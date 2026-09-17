package logic

import (
	"context"
	"encoding/json"
	"time"
	"zeromall/balance/rpc/balancePb"
	"zeromall/common/constant"
	"zeromall/common/mq"
	"zeromall/pay/rpc/internal/model"

	"zeromall/pay/rpc/internal/svc"
	"zeromall/pay/rpc/payPb"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreatePayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePayLogic {
	return &CreatePayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePayLogic) CreatePay(in *payPb.CreatePayReq) (*payPb.CreatePayResp, error) {
	// todo: add your logic here and delete this line
	//发送半消息
	payNo := uuid.NewString()
	msg := mq.PaySuccessMsg{
		PayNo:   payNo,
		OrderNo: in.OrderNo,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		l.Logger.Errorf(constant.MarshalErr, "createOrder", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	tx := l.svcCtx.TxProducer.BeginTransaction()
	_, err = l.svcCtx.TxProducer.SendWithTransaction(l.ctx, l.svcCtx.Config.RocketMqConf.Topics.TopicPaySuccess, data, tx)
	if err != nil {
		l.Logger.Errorf(constant.WhereFailed, "createPay", err.Error())
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	//调用balanceRpc扣减余额
	res, err := l.svcCtx.BalanceRpc.DeductBalance(l.ctx, &balancePb.DeductBalanceReq{
		UserId:  in.UserId,
		PayCent: in.PayCent,
	})
	if err != nil {
		_ = tx.RollBack()
		return nil, status.Error(codes.Internal, err.Error())
	}
	if res.Ok == false {
		_ = tx.RollBack()
		return nil, status.Error(codes.Internal, "扣款失败")
	}
	//插入mysql
	pay := model.Pay{
		OrderNo: in.OrderNo,
		UserId:  in.UserId,
		PayNo:   payNo,
		PayCent: in.PayCent,
		PayType: 1,
		PayTime: time.Now(),
		Status:  1,
	}
	resp, err := l.svcCtx.PayModel.Insert(l.ctx, &pay)
	if err != nil {
		//插入pay表失败，执行扣款的补偿 退款
		result, err := l.svcCtx.BalanceRpc.RefundBalance(l.ctx, &balancePb.RefundBalanceReq{
			UserId:  in.UserId,
			PayCent: in.PayCent,
		})
		if err != nil {
			l.Logger.Errorf(constant.WhereFailed, "refundBalance", err.Error())
		}
		if result.Ok == false {
			l.Logger.Info("补偿任务失败")
		}
		_ = tx.RollBack()
		l.Logger.Errorf(constant.WhereFailed, "createPay", "插入支付记录失败")
		return nil, status.Error(codes.Internal, constant.MiddlewareError)
	}
	if e := tx.Commit(); e != nil {
		l.Logger.Errorf(constant.WhereFailed, "commit err", e.Error())
		return nil, status.Error(codes.Internal, e.Error())
	}
	num, _ := resp.RowsAffected()

	return &payPb.CreatePayResp{
		Ok: num > 0,
	}, nil
}
