package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	// OrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderModel.
	OrderModel interface {
		orderModel
		withSession(session sqlx.Session) OrderModel
	}

	customOrderModel struct {
		*defaultOrderModel
	}
)

// NewOrderModel returns a model for the database table.
func NewOrderModel(conn sqlx.SqlConn) OrderModel {
	return &customOrderModel{
		defaultOrderModel: newOrderModel(conn),
	}
}

func (m *customOrderModel) withSession(session sqlx.Session) OrderModel {
	return NewOrderModel(sqlx.NewSqlConnFromSession(session))
}
func (m *defaultOrderModel) TxInsert(ctx context.Context, orderData *Order, orderItemData []*OrderItem) (int64, error) {
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
		sum = num1 + num2
		return nil
	})
	return sum, err
}
func InsertOrder(ctx context.Context, session sqlx.Session, orderData *Order) (int64, error) {
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
	placeHolder := strings.Repeat("?.", len(cols))
	placeHolder = placeHolder[:len(placeHolder)-1]
	sqlStr := fmt.Sprintf(`insert into %s (%s) values (%s)`, "order", strings.Join(cols, ","), placeHolder)
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
	onePlaceHolder := "(?,?,?,?,?,?,?)"
	var placeHolder []string
	for _, v := range orderData {
		placeHolder = append(placeHolder, onePlaceHolder)
		vals = append(vals, v.OrderNo, v.GoodsId, v.GoodsName, v.GoodsCover, v.PriceCent, v.TotalCent, v.Num)
	}
	sqlStr := fmt.Sprintf(`insert into %s (%s) values (%s)`, "order_item", strings.Join(cols, ","), strings.Join(placeHolder, ","))
	res, err := session.ExecCtx(ctx, sqlStr, vals...)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	return num, err
}

func (m *defaultOrderModel) CloseConditional(ctx context.Context, orderNo string) (int64, error) {
	sqlStr := fmt.Sprintf(`update %s set status=4,cancel_time= NOW(),updated_at= NOW()
where order_no=? and status=0 and expire_time< NOW()`, m.table)
	res, err := m.conn.ExecCtx(ctx, sqlStr, orderNo)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	return num, err
}

// 为了解决分布式事务，必须取出conn
func (m *defaultOrderModel) GetConn() sqlx.SqlConn {
	return m.conn
}
