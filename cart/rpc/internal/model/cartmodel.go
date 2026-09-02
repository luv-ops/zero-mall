package model

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CartModel = (*customCartModel)(nil)

type (
	// CartModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCartModel.
	CartModel interface {
		cartModel
		withSession(session sqlx.Session) CartModel
	}

	customCartModel struct {
		*defaultCartModel
	}
)

// NewCartModel returns a model for the database table.
func NewCartModel(conn sqlx.SqlConn) CartModel {
	return &customCartModel{
		defaultCartModel: newCartModel(conn),
	}
}

func (m *customCartModel) withSession(session sqlx.Session) CartModel {
	return NewCartModel(sqlx.NewSqlConnFromSession(session))
}
func (m *defaultCartModel) FindCartsByUserId(ctx context.Context, userId string) ([]*Cart, error) {
	sqlStr := fmt.Sprintf(`select * from %s where user_id = ?`, m.table)
	var carts []*Cart
	err := m.conn.QueryRowsPartialCtx(ctx, &carts, sqlStr, userId)
	if err != nil {
		return nil, err
	}
	return carts, nil
}
func (m *defaultCartModel) BatchDeleteTx(ctx context.Context, session sqlx.Session, userId string, goodsIds []string) error {
	if len(goodsIds) == 0 {
		return nil
	}
	placeholders := strings.Repeat("?,", len(goodsIds))
	//去调占位的最后一个 ,
	placeholders = placeholders[:len(placeholders)-1]
	sqlStr := fmt.Sprintf(`delete from %s where user_id=? and goods_id in (%s)`, m.table, placeholders)
	//[]string转[]any
	var values []any
	values = append(values, userId)
	for _, v := range goodsIds {
		values = append(values, v)
	}
	res, err := session.ExecCtx(ctx, sqlStr, values...)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		log.Printf("no affected rows in %s:%s", m.table, "batchDel")
	}
	return nil
}
func (m *defaultCartModel) BatchUpdateTx(ctx context.Context, userId string, session sqlx.Session, toUpdate []*Cart) error {
	if len(toUpdate) == 0 {
		return nil
	}
	var caseNumStr strings.Builder
	var caseSelectedStr strings.Builder
	var args []any
	for _, v := range toUpdate {
		caseNumStr.WriteString(" WHEN ? THEN ? ")
		args = append(args, v.GoodsId, v.Num)
		caseSelectedStr.WriteString(" WHEN ? THEN ? ")
		args = append(args, v.GoodsId, v.Selected)
	}
	placeholders := strings.Repeat("?,", len(toUpdate))
	placeholders = placeholders[:len(placeholders)-1]
	sqlStr := fmt.Sprintf(`update %s set num = case goods_id %s end, 
              selected = case goods_id %s end where user_id = ? and goods_id in (%s)`, m.table, caseNumStr.String(), caseSelectedStr.String(), placeholders)
	args = append(args, userId)
	for _, v := range toUpdate {
		args = append(args, v.GoodsId)
	}
	res, err := session.ExecCtx(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		log.Printf("no affected rows in %s:%s", m.table, "batchUpdate")
	}
	return nil
}
func (m *defaultCartModel) BatchInsertTx(ctx context.Context, session sqlx.Session, toInsert []*Cart) error {
	if len(toInsert) == 0 {
		return nil
	}
	var args []any
	var placeholders []string
	for _, v := range toInsert {
		placeholders = append(placeholders, "(?,?,?,?,?,?,?)")
		args = append(args, v.UserId, v.GoodsId, v.Num, v.Selected, v.Name, v.Cover, v.PriceCent)
	}
	sqlStr := fmt.Sprintf("insert into %s (user_id, goods_id,num,selected,name,cover,price_cent) values %s", m.table, strings.Join(placeholders, ","))
	res, err := session.ExecCtx(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		log.Printf("no affected rows in %s:%s", m.table, "batchInsert")
	}
	return nil
}
func (m *defaultCartModel) TransactCtx(ctx context.Context, userId string, toUpdate []*Cart, toInsert []*Cart, toDelete []string) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if len(toUpdate) > 0 {
			err := m.BatchUpdateTx(ctx, userId, session, toUpdate)
			if err != nil {
				return err
			}
		}
		if len(toInsert) > 0 {
			err := m.BatchInsertTx(ctx, session, toInsert)
			if err != nil {
				return err
			}
		}
		if len(toDelete) > 0 {
			err := m.BatchDeleteTx(ctx, session, userId, toDelete)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
