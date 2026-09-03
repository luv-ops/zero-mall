package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

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

type AddressDBRow struct {
	ReceiverName     string `db:"receiver_name"`
	ReceiverPhone    string `db:"receiver_phone"`
	Detail           string `db:"detail"`
	IsDefault        int64  `db:"is_default"`
	AddressMergeName string `db:"address_merge_name"` // 和sql as别名完全一致，下划线！
}

func (m *defaultUserReceiveAddressModel) FindAddressWithArea(ctx context.Context, userId string) ([]*AddressDBRow, error) {
	sqlStr := fmt.Sprintf(` select ua.receiver_name,ua.receiver_phone,ua.detail,ua.is_default,a.mergename as address_merge_name from %s ua 
 	left join area a on a.id = ua.address_id where ua.user_id = ?
 	`, m.table)
	var list []*AddressDBRow
	err := m.conn.QueryRowsPartialCtx(ctx, &list, sqlStr, userId)
	return list, err
}
func (m *defaultUserReceiveAddressModel) ExistReceiveAddr(ctx context.Context, userId string) (bool, error) {
	sqlStr := fmt.Sprintf(`select 1 from %s where user_id = ?`, m.table)
	var dummy int
	err := m.conn.QueryRowPartialCtx(ctx, &dummy, sqlStr, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return dummy == 1, nil
}
func (m *defaultUserReceiveAddressModel) TxRecAddrInsert(ctx context.Context, data *UserReceiveAddress) (int64, error) {
	if data.IsDefault == 1 {
		var sum int64
		//先更新所有数据的is_default为0
		err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
			//先将用户所有收货地值默认字段置0
			res, err := session.ExecCtx(ctx, fmt.Sprintf(`update %s set is_default = 0 where user_id = ?`, m.table), data.UserId)
			if err != nil {
				return err
			}
			num1, _ := res.RowsAffected()
			num2, err := RecAddrInsert(ctx, session, data)
			if err != nil {
				return err
			}
			sum = num1 + num2
			return nil
		})
		if err != nil {
			return 0, err
		}
		return sum, nil

	}
	//普通插入
	return RecAddrInsert(ctx, m.conn, data)

}
func RecAddrInsert(ctx context.Context, conn sqlx.Session, data *UserReceiveAddress) (int64, error) {
	sqlStr := fmt.Sprintf(`insert into %s (user_id,receiver_name,receiver_phone,address_id,detail,is_default) 
			values (?,?,?,?,?,?)`, "user_receive_address")
	res, err := conn.ExecCtx(ctx, sqlStr, data.UserId, data.ReceiverName, data.ReceiverPhone, data.AddressId, data.Detail, data.IsDefault)
	num, _ := res.RowsAffected()

	return num, err
}
func (m *defaultUserReceiveAddressModel) FindOneByUIdWithDefault(ctx context.Context, userId string) (*AddressDBRow, error) {
	sqlStr := fmt.Sprintf(`select ua.receiver_name,ua.receiver_phone,a.mergename as address_merge_name from %s ua 
	left join area a on a.id=ua.address_id where user_id=? and is_default=1`, m.table)
	var item AddressDBRow
	err := m.conn.QueryRowPartialCtx(ctx, &item, sqlStr, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
