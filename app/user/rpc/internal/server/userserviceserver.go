// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"gomall/app/user/rpc/internal/logic"
	"gomall/app/user/rpc/internal/svc"
	"gomall/app/user/rpc/pb"
)

type UserServiceServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedUserServiceServer
}

func NewUserServiceServer(svcCtx *svc.ServiceContext) *UserServiceServer {
	return &UserServiceServer{
		svcCtx: svcCtx,
	}
}

// 用户注册
func (s *UserServiceServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

// 用户登录
func (s *UserServiceServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

func (s *UserServiceServer) Profile(ctx context.Context, in *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	l := logic.NewProfileLogic(ctx, s.svcCtx)
	return l.Profile(in)
}
