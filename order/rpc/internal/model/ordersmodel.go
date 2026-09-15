package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrdersModel = (*customOrdersModel)(nil)

type (
	// OrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrdersModel.
	OrdersModel interface {
		ordersModel
		withSession(session sqlx.Session) OrdersModel
	}

	customOrdersModel struct {
		*defaultOrdersModel
	}
)

// NewOrdersModel returns a model for the database table.
func NewOrdersModel(conn sqlx.SqlConn) OrdersModel {
	return &customOrdersModel{
		defaultOrdersModel: newOrdersModel(conn),
	}
}

func (m *customOrdersModel) withSession(session sqlx.Session) OrdersModel {
	return NewOrdersModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultOrdersModel) TxInsert(ctx context.Context, orderData *Orders, orderItemData []*OrderItem, txLogData *TransactionLog) (int64, error) {
	var sum int64
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		num1, err := InsertOrder(ctx, session, orderData)
		if err != nil {
			return err
		}
		num2, err := InsertOrderItems(ctx, session, orderItemData)
		if err != nil {
			return err
		}
		num3, err := InsertTxLog(ctx, session, txLogData)
		sum = num1 + num2 + num3
		return nil
	})
	return sum, err
}
func InsertOrder(ctx context.Context, session sqlx.Session, orderData *Orders) (int64, error) {
	if orderData == nil {
		return 0, nil
	}
	var cols []string
	var vals []any
	cols = append(cols, "order_no", "user_id", "total_cent", "pay_cent", "expire_time")
	vals = append(vals, orderData.OrderNo, orderData.UserId, orderData.TotalCent, orderData.PayCent, orderData.ExpireTime)
	if orderData.Remark.Valid || orderData.Remark.String != "" {
		cols = append(cols, "remark")
		vals = append(vals, orderData.Remark.String)
	}
	placeHolder := strings.Repeat("?,", len(cols))
	placeHolder = placeHolder[:len(placeHolder)-1]
	sqlStr := fmt.Sprintf(`insert into %s (%s) values (%s)`, "orders", strings.Join(cols, ","), placeHolder)
	res, err := session.ExecCtx(ctx, sqlStr, vals...)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	return num, err
}
func InsertOrderItems(ctx context.Context, session sqlx.Session, orderData []*OrderItem) (int64, error) {
	if len(orderData) == 0 {
		return 0, nil
	}
	var cols []string
	var vals []any
	cols = append(cols, "order_no", "goods_id", "goods_name", "goods_cover", "price_cent", "total_cent", "num")
	var placeHolder []string
	for _, v := range orderData {
		placeHolder = append(placeHolder, "(?,?,?,?,?,?,?)")
		vals = append(vals, v.OrderNo, v.GoodsId, v.GoodsName, v.GoodsCover, v.PriceCent, v.TotalCent, v.Num)
	}
	sqlStr := fmt.Sprintf(`insert into %s (%s) values %s`, "order_item", strings.Join(cols, ","), strings.Join(placeHolder, ","))
	res, err := session.ExecCtx(ctx, sqlStr, vals...)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	return num, err
}
func InsertTxLog(ctx context.Context, session sqlx.Session, txLogData *TransactionLog) (int64, error) {
	sqlStr := fmt.Sprintf("insert into %s (tx_id,order_no) values (?,?)", "transaction_log")
	vals := []any{txLogData.TxId, txLogData.OrderNo}
	res, err := session.ExecCtx(ctx, sqlStr, vals...)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	return num, err
}

func (m *defaultOrdersModel) CloseConditional(ctx context.Context, orderNo int64) (int64, error) {
	sqlStr := fmt.Sprintf(`update %s set status=4,cancel_time= NOW(),updated_at= NOW()
where order_no=? and status=0 and expire_time< NOW()`, m.table)
	res, err := m.conn.ExecCtx(ctx, sqlStr, orderNo)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	return num, err
}
