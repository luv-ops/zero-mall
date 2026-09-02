// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"
	"zeromall/common/Res"

	"zeromall/user/api/internal/logic"
	"zeromall/user/api/internal/svc"
	"zeromall/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getOtherInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OthersBaseInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetOtherInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetOtherInfo(&req)
		Res.Response(r, w, resp, err)
	}
}
