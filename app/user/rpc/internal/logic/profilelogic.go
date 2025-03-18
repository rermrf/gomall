package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"gomall/app/user/rpc/internal/svc"
	"gomall/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProfileLogic) Profile(in *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		return nil, errors.New(400, "未找到")
	}
	return &pb.ProfileResponse{
		User: &pb.User{
			Id:       user.Id,
			Phone:    user.Phone.String,
			Email:    user.Email.String,
			Password: user.Password,
			NickName: user.NickName.String,
			Avatar:   user.Avatar.String,
			Ctime:    timestamppb.New(time.UnixMilli(user.Ctime)),
		},
	}, nil
}
