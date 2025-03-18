package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gomall/pkg/jwtx"
	"net/http"
)

type AuthMiddleware struct {
	jwtHandler   jwtx.Handler
	AccessSecret string
}

func NewAuthMiddleware(jwtHandler jwtx.Handler, accessSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtHandler: jwtHandler, AccessSecret: accessSecret}
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		tokenStr := m.jwtHandler.ExtractToken(r)
		if tokenStr == "" {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, Result{
				Code: http.StatusUnauthorized,
				Msg:  "未登录",
			})
			return
		}
		claims := &jwtx.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.AccessSecret), nil
		})
		if err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, Result{
				Code: http.StatusUnauthorized,
				Msg:  "验证有误",
			})
		}
		if token == nil || !token.Valid || claims.UserId == 0 {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, Result{
				Code: http.StatusUnauthorized,
				Msg:  "无效令牌",
			})
			return
		}

		// 验证 UserAgent
		if claims.UserAgent != r.UserAgent() {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, Result{
				Code: http.StatusUnauthorized,
				Msg:  "登录环境变更，请重新登录",
			})
			return
		}

		// 验证ssid
		if err = m.jwtHandler.CheckSession(r, claims.Ssid); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusUnauthorized, Result{
				Code: http.StatusUnauthorized,
				Msg:  "会话已消失",
			})
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)
		ctx = context.WithValue(ctx, "userId", claims.UserId)
		// Passthrough to next handler if need
		next(w, r.WithContext(ctx))
	}
}
