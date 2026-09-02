package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		withSession(session sqlx.Session) UserModel
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *customUserModel) withSession(session sqlx.Session) UserModel {
	return NewUserModel(sqlx.NewSqlConnFromSession(session))
}

// 自定义sql
func (m *defaultUserModel) UpdateField(ctx context.Context, updateMap map[string]any, userId string) (int64, error) {
	if len(updateMap) == 0 {
		return 0, nil
	}
	var keys []string
	var values []any
	for k, v := range updateMap {
		keys = append(keys, fmt.Sprintf("%s=?", k))
		values = append(values, v)
	}
	setSql := strings.Join(keys, ",")
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE `user_id` = ?", m.table, setSql)
	values = append(values, userId)
	res, err := m.conn.ExecCtx(ctx, sqlStr, values...)
	affectedRow, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affectedRow, err
}

func (m *defaultUserModel) SelectOneByField(ctx context.Context, userId string, fields ...string) (*User, error) {
	selectFields := strings.Join(fields, ",")
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE user_id=?", selectFields, m.table)
	user := &User{}
	err := m.conn.QueryRowPartialCtx(ctx, user, sqlStr, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
