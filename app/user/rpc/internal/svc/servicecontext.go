package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gomall/app/user/rpc/internal/config"
	"gomall/app/user/rpc/internal/model"
)

type ServiceContext struct {
	Config     config.Config
	UserModel  model.UserModel
	CacheRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	rdb := redis.MustNewRedis(c.CacheRedis[0].RedisConf)
	return &ServiceContext{
		Config:     c,
		UserModel:  model.NewUserModel(conn, c.CacheRedis),
		CacheRedis: rdb,
	}
}
