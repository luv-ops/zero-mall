package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ PayModel = (*customPayModel)(nil)

type (
	// PayModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPayModel.
	PayModel interface {
		payModel
		withSession(session sqlx.Session) PayModel
	}

	customPayModel struct {
		*defaultPayModel
	}
)

// NewPayModel returns a model for the database table.
func NewPayModel(conn sqlx.SqlConn) PayModel {
	return &customPayModel{
		defaultPayModel: newPayModel(conn),
	}
}

func (m *customPayModel) withSession(session sqlx.Session) PayModel {
	return NewPayModel(sqlx.NewSqlConnFromSession(session))
}
