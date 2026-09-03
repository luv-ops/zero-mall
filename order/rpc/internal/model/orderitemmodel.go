package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

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
