// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"
	"zeromall/common/Res"

	"zeromall/goods/api/internal/logic"
	"zeromall/goods/api/internal/svc"
	"zeromall/goods/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetGoodsDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GoodsDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetGoodsDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetGoodsDetail(&req)
		Res.Response(r, w, resp, err)
	}
}
