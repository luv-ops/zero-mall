package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TransactionLogModel = (*customTransactionLogModel)(nil)

type (
	// TransactionLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTransactionLogModel.
	TransactionLogModel interface {
		transactionLogModel
		withSession(session sqlx.Session) TransactionLogModel
	}

	customTransactionLogModel struct {
		*defaultTransactionLogModel
	}
)

// NewTransactionLogModel returns a model for the database table.
func NewTransactionLogModel(conn sqlx.SqlConn) TransactionLogModel {
	return &customTransactionLogModel{
		defaultTransactionLogModel: newTransactionLogModel(conn),
	}
}

func (m *customTransactionLogModel) withSession(session sqlx.Session) TransactionLogModel {
	return NewTransactionLogModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultTransactionLogModel) FindStatusByOrderNo(orderNo int64) (*TransactionLog, error) {
	sqlStr := fmt.Sprintf("select status from %s where order_no=?", m.table)
	var tlog TransactionLog
	err := m.conn.QueryRowPartialCtx(context.TODO(), &tlog, sqlStr, orderNo)
	return &tlog, err
}
func (m *defaultTransactionLogModel) UpdateStatusByOrderNo(orderNo int64) error {
	sqlStr := fmt.Sprintf("update %s set status=1 where order_no=? and status=0", m.table)
	_, err := m.conn.ExecCtx(context.TODO(), sqlStr, orderNo)
	return err
}
