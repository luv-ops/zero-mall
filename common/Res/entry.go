package Res

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Response(r *http.Request, w http.ResponseWriter, resp any, err error) {
	if err != nil {
		body := Body{
			Code: -1,
			Msg:  err.Error(),
			Data: nil,
		}
		httpx.WriteJson(w, http.StatusOK, body)
		return
	}
	body := Body{
		Code: 0,
		Msg:  "success",
		Data: resp,
	}
	httpx.WriteJson(w, http.StatusOK, body)
	return
}
