// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"
	"zeromall/common/Res"

	"zeromall/cart/api/internal/logic"
	"zeromall/cart/api/internal/svc"
)

func GetCartListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetCartListLogic(r.Context(), svcCtx)
		resp, err := l.GetCartList()
		Res.Response(r, w, resp, err)
	}
}
