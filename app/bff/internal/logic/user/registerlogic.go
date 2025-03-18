package user

import (
	"context"
	regexp "github.com/dlclark/regexp2"
	"github.com/zeromicro/x/errors"
	"gomall/app/user/rpc/userservice"

	"gomall/app/bff/internal/svc"
	"gomall/app/bff/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	phoneExp    *regexp.Regexp
	PasswordExp *regexp.Regexp
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	const (
		emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*?&.])[A-Za-z\d@$!%*?&.]{8,72}$`
		phoneRegexPattern    = `^1[3-9]\d{9}$`
	)
	return &RegisterLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		phoneExp:    regexp.MustCompile(phoneRegexPattern, regexp.None),
		PasswordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// todo: add your logic here and delete this line
	ok, err := l.phoneExp.MatchString(req.Phone)
	if err != nil {
		return nil, errors.New(500, "系统错误")
	}
	if !ok {
		return nil, errors.New(400, "手机号匹配错误")
	}
	if req.Password != req.ConfirmPassword {
		return nil, errors.New(400, "两次密码不一致")
	}
	ok, err = l.PasswordExp.MatchString(req.Password)
	if err != nil {
		return nil, errors.New(500, "系统错误")
	}
	if !ok {
		return nil, errors.New(400, "密码格式有误")
	}
	_, err = l.svcCtx.UserRpc.Register(l.ctx, &userservice.RegisterRequest{
		User: &userservice.User{
			Phone:    req.Phone,
			Password: req.Password,
		},
	})
	if err != nil {
		// TODO 使用特定的err判定是否冲突
		return nil, errors.New(500, err.Error())
	}
	return
}
