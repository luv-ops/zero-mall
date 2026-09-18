package xerr

import (
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

// 错误码
const (
	OK                  = 0
	ParamErr            = 10001 // 参数错误
	NotFound            = 10002 // 找不到
	PhoneIllegal        = 10003
	CaptchaMistake      = 10004
	UserHasExist        = 10005
	PasswordErr         = 10006
	AddGoodsErr         = 10007
	AddUserErr          = 10008
	GoodsNotSell        = 10009
	GoodsNotInCart      = 10010
	GoodsNotSelected    = 10011
	DeleteCartErr       = 10012
	StockNotEnough      = 10013
	CreateOrderErr      = 10014
	PermissionDenied    = 10015
	OrderStatusChange   = 10016
	OrderNotCancel      = 10017
	OrderOffErr         = 10018
	DeductBalanceErr    = 10019
	BalanceNotEnough    = 10020
	MsgSmsTooFrequently = 10021
	ServerErr           = 50000 // 系统错误

)

// 错误码 → 友好提示
var msgMap = map[int]string{
	ParamErr:            "参数错误",
	NotFound:            "资源不存在",
	ServerErr:           "服务器繁忙，请稍后再试",
	PhoneIllegal:        "手机号格式错误",
	CaptchaMistake:      "验证码错误",
	UserHasExist:        "用户已存在",
	PasswordErr:         "密码错误",
	AddGoodsErr:         "添加商品失败",
	AddUserErr:          "注册用户失败",
	GoodsNotSell:        "商品已下架",
	GoodsNotInCart:      "商品不再购物车中",
	GoodsNotSelected:    "商品未勾选",
	DeleteCartErr:       "从购物车删除商品失败",
	StockNotEnough:      "库存不足",
	CreateOrderErr:      "创建订单失败",
	PermissionDenied:    "访问越权",
	OrderStatusChange:   "订单状态已变化",
	OrderNotCancel:      "订单不可取消",
	OrderOffErr:         "订单取消失败",
	DeductBalanceErr:    "扣款失败",
	BalanceNotEnough:    "余额不足",
	MsgSmsTooFrequently: "短信发送过于频繁，请稍后再试",
}

// CodeError 自定义业务错误类型（这就是官方示例里的结构）
type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// CodeErr实现了Error方法，它的结构体就属于error类型
func (e *CodeError) Error() string { return e.Msg }
func (e *CodeError) GetCode() int  { return e.Code }

// NewCodeError 根据 code 自动查表生成错误
func NewCodeError(code int) error {
	return &CodeError{Code: code, Msg: msgMap[code]}
}
func NewCodeErrorWithMsg(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

// NewDefaultError 只指定消息，使用默认错误码
func NewDefaultError(msg string) error {
	return &CodeError{Code: ServerErr, Msg: msg}
}
func IsCodeError(code int) bool {
	_, ok := msgMap[code]
	return ok
}
func FromRpcError(err error) error {
	// 1. 解析 gRPC 状态，非 gRPC 错误直接当系统错误
	st, ok := status.FromError(err)
	if !ok {
		logx.Errorf("not a grpc error: %+v", err)
		return NewCodeError(ServerErr)
	}
	if IsCodeError(int(st.Code())) {
		return NewCodeError(int(st.Code()))
	}
	logx.Errorf("rpc system error, code: %d, msg: %s, err: %+v",
		st.Code(), st.Message(), err)
	return NewCodeError(ServerErr)
}
func Server() error {
	return NewCodeError(ServerErr)
}
