package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserReceiveAddressModel = (*customUserReceiveAddressModel)(nil)

type (
	// UserReceiveAddressModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserReceiveAddressModel.
	UserReceiveAddressModel interface {
		userReceiveAddressModel
		withSession(session sqlx.Session) UserReceiveAddressModel
	}

	customUserReceiveAddressModel struct {
		*defaultUserReceiveAddressModel
	}
)

// NewUserReceiveAddressModel returns a model for the database table.
func NewUserReceiveAddressModel(conn sqlx.SqlConn) UserReceiveAddressModel {
	return &customUserReceiveAddressModel{
		defaultUserReceiveAddressModel: newUserReceiveAddressModel(conn),
	}
}

func (m *customUserReceiveAddressModel) withSession(session sqlx.Session) UserReceiveAddressModel {
	return NewUserReceiveAddressModel(sqlx.NewSqlConnFromSession(session))
}
