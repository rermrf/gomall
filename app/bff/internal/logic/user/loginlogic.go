package user

import (
	"context"
	regexp "github.com/dlclark/regexp2"
	"github.com/zeromicro/x/errors"
	"gomall/app/user/rpc/pb"

	"gomall/app/bff/internal/svc"
	"gomall/app/bff/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	phoneExp    *regexp.Regexp
	PasswordExp *regexp.Regexp
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	const (
		emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*?&.])[A-Za-z\d@$!%*?&.]{8,72}$`
		phoneRegexPattern    = `^1[3-9]\d{9}$`
	)
	return &LoginLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		phoneExp:    regexp.MustCompile(phoneRegexPattern, regexp.None),
		PasswordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line
	ok, err := l.phoneExp.MatchString(req.Phone)
	if err != nil {
		return nil, errors.New(500, "系统错误")
	}
	if !ok {
		return nil, errors.New(400, "手机格式不正确")
	}
	ok, err = l.PasswordExp.MatchString(req.Password)
	if err != nil {
		return nil, errors.New(500, "系统错误")
	}
	if !ok {
		return nil, errors.New(400, "密码格式不正确")
	}
	res, err := l.svcCtx.UserRpc.Login(l.ctx, &pb.LoginRequest{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, errors.New(500, err.Error())
	}
	resp = &types.LoginResponse{
		Uid: res.GetUser().GetId(),
	}
	return
}
