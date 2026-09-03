package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GoodsModel = (*customGoodsModel)(nil)

type (
	// GoodsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGoodsModel.
	GoodsModel interface {
		goodsModel
		withSession(session sqlx.Session) GoodsModel
	}

	customGoodsModel struct {
		*defaultGoodsModel
	}
)

// NewGoodsModel returns a model for the database table.
func NewGoodsModel(conn sqlx.SqlConn) GoodsModel {
	return &customGoodsModel{
		defaultGoodsModel: newGoodsModel(conn),
	}
}

func (m *customGoodsModel) withSession(session sqlx.Session) GoodsModel {
	return NewGoodsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultGoodsModel) PageBreakFind(ctx context.Context, categoryId int64, page int32, pageSize int32) ([]*Goods, error) {
	var list []*Goods
	sqlStr := fmt.Sprintf("select * from %s where status=1 ", m.table)
	if categoryId != 0 {
		sqlStr += fmt.Sprintf(" and category_id=%d ", categoryId)
	}
	sqlStr += fmt.Sprintf("order by created_at desc limit ?,?")
	err := m.conn.QueryRowsCtx(ctx, &list, sqlStr, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	return list, nil
}
func (m *defaultGoodsModel) FindByOwnId(ctx context.Context, userId string) ([]*Goods, error) {
	var list []*Goods
	sqlStr := fmt.Sprintf("select * from %s where own_userId=? ", m.table)
	err := m.conn.QueryRowsCtx(ctx, &list, sqlStr, userId)
	if err != nil {
		return nil, err
	}
	return list, err
}
func (m *defaultGoodsModel) UpdateFields(ctx context.Context, goodsId string, setMap map[string]any) (int64, error) {
	if len(setMap) == 0 {
		return 0, nil
	}
	var keys []string
	var values []any
	for k, v := range setMap {
		keys = append(keys, fmt.Sprintf("%s=?", k))
		values = append(values, v)
	}
	keyStr := strings.Join(keys, ",")
	values = append(values, goodsId)
	sqlStr := fmt.Sprintf("update %s set %s where goods_id=? ", m.table, keyStr)
	res, err := m.conn.ExecCtx(ctx, sqlStr, values...)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return num, nil
}
func (m *defaultGoodsModel) FindRowsByGoodsId(ctx context.Context, goodsIds []string) ([]*Goods, error) {
	placeHolder := strings.Repeat("?,", len(goodsIds))
	placeHolder = placeHolder[:len(placeHolder)-1]
	sqlStr := fmt.Sprintf("select goods_id,name,cover,price_cent,original_price_cent from %s where goods_id in (%s) ", m.table, placeHolder)
	var list []*Goods
	var anyIds []any
	for _, goodsId := range goodsIds {
		anyIds = append(anyIds, goodsId)
	}
	err := m.conn.QueryRowsPartialCtx(ctx, &list, sqlStr, anyIds...)
	return list, err
}
