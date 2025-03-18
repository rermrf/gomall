package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	xhttp "github.com/zeromicro/x/http"
	"gomall/app/bff/internal/logic/user"
	"gomall/app/bff/internal/svc"
	"gomall/app/bff/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			err := svcCtx.JwtHandler.SetLoginToken(w, r, resp.Uid)
			if err != nil {
				xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			}
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
