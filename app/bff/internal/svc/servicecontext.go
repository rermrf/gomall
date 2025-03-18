package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"gomall/app/bff/internal/config"
	"gomall/app/bff/internal/middleware"
	"gomall/app/user/rpc/pb"
	"gomall/app/user/rpc/userservice"
	"gomall/pkg/jwtx"
)

type ServiceContext struct {
	Config         config.Config
	AuthMiddleware rest.Middleware
	BizRedis       *redis.Redis
	JwtHandler     jwtx.Handler
	UserRpc        pb.UserServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	bizRedis := redis.MustNewRedis(c.BizRedis)
	jwtHandler := jwtx.NewRedisJWTHandler(bizRedis, c.Auth.AccessSecret, c.Auth.AccessExpire, c.Auth.RefreshSecret, c.Auth.AccessExpire)
	userRpc := zrpc.MustNewClient(c.UserRpc)
	return &ServiceContext{
		Config:         c,
		AuthMiddleware: middleware.NewAuthMiddleware(jwtHandler, c.Auth.AccessSecret).Handle,
		BizRedis:       bizRedis,
		JwtHandler:     jwtHandler,
		UserRpc:        userservice.NewUserService(userRpc),
	}
}
