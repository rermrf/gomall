package user

import (
	"context"
	"github.com/zeromicro/x/errors"
	"gomall/app/user/rpc/userservice"

	"gomall/app/bff/internal/svc"
	"gomall/app/bff/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileLogic) Profile() (resp *types.ProfileResponse, err error) {
	// todo: add your logic here and delete this line
	uid, ok := l.ctx.Value("userId").(int64)
	if !ok {
		return nil, errors.New(500, "系统错误")
	}
	res, err := l.svcCtx.UserRpc.Profile(l.ctx, &userservice.ProfileRequest{
		Uid: uid,
	})
	if err != nil {
		return nil, errors.New(500, err.Error())
	}
	resp = &types.ProfileResponse{
		Phone:    res.GetUser().Phone,
		Email:    res.GetUser().Email,
		NickName: res.GetUser().NickName,
		Avatar:   res.GetUser().Avatar,
		Ctime:    res.GetUser().Ctime.AsTime().UnixMilli(),
	}
	return
}
