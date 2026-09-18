package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BalanceModel = (*customBalanceModel)(nil)

type (
	// BalanceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBalanceModel.
	BalanceModel interface {
		balanceModel
		withSession(session sqlx.Session) BalanceModel
	}

	customBalanceModel struct {
		*defaultBalanceModel
	}
)

// NewBalanceModel returns a model for the database table.
func NewBalanceModel(conn sqlx.SqlConn) BalanceModel {
	return &customBalanceModel{
		defaultBalanceModel: newBalanceModel(conn),
	}
}

func (m *customBalanceModel) withSession(session sqlx.Session) BalanceModel {
	return NewBalanceModel(sqlx.NewSqlConnFromSession(session))
}
func (m *defaultBalanceModel) DeductBalance(ctx context.Context, payCent int64, version int64, userId string) (int64, error) {
	sqlStr := fmt.Sprintf(`update %s set available_cent=available_cent-?,version=version+1 where user_id=? and version=?`, m.table)
	args := []interface{}{payCent, userId, version}
	res, err := m.conn.ExecCtx(ctx, sqlStr, args...)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	return num, err
}
func (m *defaultBalanceModel) RefundBalance(ctx context.Context, payCent int64, userId string) (int64, error) {
	selectStr := fmt.Sprintf(`select version from %s where user_id=? `, m.table)
	var balance Balance
	err := m.conn.QueryRowPartialCtx(ctx, &balance, selectStr, userId)
	if err != nil {
		return 0, err
	}
	sqlStr := fmt.Sprintf(`update %s set available_cent=available_cent+?,version=version+1 where user_id=? and version=?`, m.table)
	args := []interface{}{payCent, userId, balance.Version}
	res, err := m.conn.ExecCtx(ctx, sqlStr, args...)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	return num, err
}
