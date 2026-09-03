package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ OrderLogisticsModel = (*customOrderLogisticsModel)(nil)

type (
	// OrderLogisticsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderLogisticsModel.
	OrderLogisticsModel interface {
		orderLogisticsModel
		withSession(session sqlx.Session) OrderLogisticsModel
	}

	customOrderLogisticsModel struct {
		*defaultOrderLogisticsModel
	}
)

// NewOrderLogisticsModel returns a model for the database table.
func NewOrderLogisticsModel(conn sqlx.SqlConn) OrderLogisticsModel {
	return &customOrderLogisticsModel{
		defaultOrderLogisticsModel: newOrderLogisticsModel(conn),
	}
}

func (m *customOrderLogisticsModel) withSession(session sqlx.Session) OrderLogisticsModel {
	return NewOrderLogisticsModel(sqlx.NewSqlConnFromSession(session))
}
