package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderItemModel = (*customOrderItemModel)(nil)

type (
	// OrderItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderItemModel.
	OrderItemModel interface {
		orderItemModel
		withSession(session sqlx.Session) OrderItemModel
	}

	customOrderItemModel struct {
		*defaultOrderItemModel
	}
)

// NewOrderItemModel returns a model for the database table.
func NewOrderItemModel(conn sqlx.SqlConn) OrderItemModel {
	return &customOrderItemModel{
		defaultOrderItemModel: newOrderItemModel(conn),
	}
}

func (m *customOrderItemModel) withSession(session sqlx.Session) OrderItemModel {
	return NewOrderItemModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultOrderItemModel) FindGIdsByOrderNo(ctx context.Context, orderNo int64) ([]*OrderItem, error) {
	sqlStr := fmt.Sprintf(`select goods_id,num from %s where order_no=?`, m.table)
	var list []*OrderItem
	err := m.conn.QueryRowsPartialCtx(ctx, &list, sqlStr, orderNo)
	return list, err
}
