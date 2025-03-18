package jwtx

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"net/http"
	"strings"
	"time"
)

type RefreshClaims struct {
	UserId    int64
	UserAgent string `json:"userAgent"`
	Ssid      string `json:"ssid"`
	jwt.RegisteredClaims
}

type UserClaims struct {
	UserId    int64  `json:"userId"`
	UserAgent string `json:"userAgent"`
	Ssid      string `json:"ssid"`
	jwt.RegisteredClaims
}

type RedisJWTHandler struct {
	Rds           *redis.Redis
	AccessSecret  string
	AccessExpire  int64
	RefreshSecret string
	RefreshExpire int64
}

func NewRedisJWTHandler(rds *redis.Redis, accessSecret string, accessExpire int64, refreshSecret string, refreshExpire int64) Handler {
	return &RedisJWTHandler{Rds: rds, AccessSecret: accessSecret, AccessExpire: accessExpire, RefreshSecret: refreshSecret, RefreshExpire: refreshExpire}
}

func (r2 *RedisJWTHandler) ExtractToken(r *http.Request) string {
	// 使用 JWT 进行登录校验
	tokenHeader := r.Header.Get("Authorization")
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (r2 *RedisJWTHandler) SetJWTToken(w http.ResponseWriter, r *http.Request, userId int64, ssid string) error {
	claims := &UserClaims{
		UserId:    userId,
		UserAgent: r.UserAgent(),
		Ssid:      ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(r2.AccessExpire) * time.Second)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(r2.AccessSecret))
	if err != nil {
		return err
	}
	w.Header().Set("x-jwt-token", tokenStr)
	return nil
}

func (r2 *RedisJWTHandler) ClearToken(w http.ResponseWriter, r *http.Request) error {
	// 退出登录前将 header 设置为非法值
	w.Header().Set("x-jwt-token", "")
	w.Header().Set("x-refresh-token", "")
	claims := r.Context().Value("claims").(*UserClaims)
	return r2.Rds.SetexCtx(r.Context(), fmt.Sprintf("users:ssid:%s", claims.Ssid), "", int(r2.AccessExpire))
}

func (r2 *RedisJWTHandler) SetLoginToken(w http.ResponseWriter, r *http.Request, userId int64) error {
	ssid := uuid.New().String()
	// 设置登录token
	err := r2.SetJWTToken(w, r, userId, ssid)
	if err != nil {
		return err
	}
	// 设置刷新token
	err = r2.SetRefreshToken(w, r, userId, ssid)
	return err
}

func (r2 *RedisJWTHandler) SetRefreshToken(w http.ResponseWriter, r *http.Request, userId int64, ssid string) error {
	claims := RefreshClaims{
		UserId: userId,
		Ssid:   ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(r2.RefreshExpire) * time.Second)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(r2.RefreshSecret))
	if err != nil {
		return err
	}
	w.Header().Set("x-refresh-token", tokenStr)
	return nil
}

func (r2 *RedisJWTHandler) CheckSession(r *http.Request, ssid string) error {
	val, err := r2.Rds.ExistsManyCtx(r.Context(), fmt.Sprintf("users:ssid:%s", ssid))
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if val > 0 {
			return errors.New("session 已经无效")
		}
	default:
		return err
	}
	return err
}
