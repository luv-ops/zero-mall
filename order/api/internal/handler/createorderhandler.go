// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"
	"zeromall/common/Res"

	"zeromall/order/api/internal/logic"
	"zeromall/order/api/internal/svc"
	"zeromall/order/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func createOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreateOrder(&req)
		Res.Response(r, w, resp, err)
	}
}
