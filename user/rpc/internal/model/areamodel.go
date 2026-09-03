package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AreaModel = (*customAreaModel)(nil)

type (
	// AreaModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAreaModel.
	AreaModel interface {
		areaModel
		withSession(session sqlx.Session) AreaModel
	}

	customAreaModel struct {
		*defaultAreaModel
	}
)

// NewAreaModel returns a model for the database table.
func NewAreaModel(conn sqlx.SqlConn) AreaModel {
	return &customAreaModel{
		defaultAreaModel: newAreaModel(conn),
	}
}

func (m *customAreaModel) withSession(session sqlx.Session) AreaModel {
	return NewAreaModel(sqlx.NewSqlConnFromSession(session))
}

// 自定义sql
func (m *defaultAreaModel) SelectFields(ctx context.Context, level int64, pId int64) ([]*Area, error) {

	sqlStr := fmt.Sprintf("select * from %s where pid=? and level=?", m.table)
	var list []*Area
	err := m.conn.QueryRowsCtx(ctx, &list, sqlStr, pId, level)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *defaultAreaModel) AreaLevel3IsExist(ctx context.Context, addressId int64) (bool, error) {
	sqlStr := fmt.Sprintf("select level from %s where id=? ", m.table)
	var level int64
	err := m.conn.QueryRowPartialCtx(ctx, &level, sqlStr, addressId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return level == 3, nil
}
