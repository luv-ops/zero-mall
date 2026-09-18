package apiHandler

import (
	"context"
	"errors"
	"net/http"
	"zeromall/common/xerr"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ApiHandler() {
	httpx.SetOkHandler(func(ctx context.Context, a any) any {
		type Body struct {
			Code int         `json:"code"`
			Msg  string      `json:"msg"`
			Data interface{} `json:"data"`
		}
		return Body{Code: 200, Msg: "success", Data: a}
	})
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		var e *xerr.CodeError
		if errors.As(err, &e) {
			return http.StatusOK, map[string]interface{}{
				"code": e.GetCode(), "message": e.Error(), "data": nil,
			}
		}
		// 非 CodeError 的未知错误 → 通用提示
		return http.StatusOK, map[string]interface{}{
			"code": xerr.ServerErr, "message": "服务繁忙，请稍后重试",
		}
	})
}
