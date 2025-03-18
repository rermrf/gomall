package logic

import (
	"context"
	"database/sql"
	"github.com/zeromicro/x/errors"
	"golang.org/x/crypto/bcrypt"

	"gomall/app/user/rpc/internal/svc"
	"gomall/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 用户登录
func (l *LoginLogic) Login(in *pb.LoginRequest) (*pb.LoginResponse, error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, sql.NullString{
		String: in.Phone,
		Valid:  in.Phone != "",
	})
	if err != nil {
		return nil, errors.New(400, "手机号或密码错误")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		return nil, errors.New(400, "手机号或密码错误")
	}
	return &pb.LoginResponse{
		User: &pb.User{
			Id: user.Id,
		},
	}, nil
}
