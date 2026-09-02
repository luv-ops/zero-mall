// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"
	"zeromall/common/Res"

	"zeromall/goods/api/internal/logic"
	"zeromall/goods/api/internal/svc"
)

func GetAdminGoodsListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetAdminGoodsListLogic(r.Context(), svcCtx)
		resp, err := l.GetAdminGoodsList()
		Res.Response(r, w, resp, err)
	}
}
