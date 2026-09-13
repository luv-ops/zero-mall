package model

import (
	"context"
	"fmt"
	"zeromall/common/mq"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StockModel = (*customStockModel)(nil)

type (
	// StockModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStockModel.
	StockModel interface {
		stockModel
		withSession(session sqlx.Session) StockModel
	}

	customStockModel struct {
		*defaultStockModel
	}
)

// NewStockModel returns a model for the database table.
func NewStockModel(conn sqlx.SqlConn) StockModel {
	return &customStockModel{
		defaultStockModel: newStockModel(conn),
	}
}

func (m *customStockModel) withSession(session sqlx.Session) StockModel {
	return NewStockModel(sqlx.NewSqlConnFromSession(session))
}
func (m *defaultStockModel) StockBatchReturn(ctx context.Context, list []*mq.ReturnStockItem) (int64, error) {
	var sum int64
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, item := range list {
			num, err := m.BatchUpdate(ctx, session, item)
			if err != nil {
				return err
			}
			sum += num
		}
	})
	return sum, err
}

func (m *defaultStockModel) BatchUpdate(ctx context.Context, session sqlx.Session, item *mq.ReturnStockItem) (int64, error) {
	var stock Stock
	selectStr := fmt.Sprintf("select version from %s where goods_id=?", m.table)
	err := session.QueryRowPartialCtx(ctx, &stock, selectStr, item.GoodsId)
	if err != nil {
		return 0, err
	}
	//乐观锁
	sqlStr := fmt.Sprintf("update %s set frozen_stock=frozen_stock-? and available_stock=available_stock+? and version=version+1 where goodsId=? and version=?", m.table)
	vals := []interface{}{item.Num, item.Num, item.GoodsId, stock.Version}
	res, err := session.ExecCtx(ctx, sqlStr, vals...)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	return num, err
}
