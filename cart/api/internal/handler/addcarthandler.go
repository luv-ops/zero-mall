// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"
	"zeromall/common/Res"

	"zeromall/cart/api/internal/logic"
	"zeromall/cart/api/internal/svc"
	"zeromall/cart/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddCartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddCartReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAddCartLogic(r.Context(), svcCtx)
		err := l.AddCart(&req)
		Res.Response(r, w, nil, err)
	}
}
