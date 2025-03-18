package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	xerrors "github.com/zeromicro/x/errors"
	"golang.org/x/crypto/bcrypt"
	"gomall/app/user/rpc/internal/model"
	"time"

	"gomall/app/user/rpc/internal/svc"
	"gomall/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Register 用户注册
func (l *RegisterLogic) Register(in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	now := time.Now().UnixMilli()
	u := in.GetUser()
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hash)
	res, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Ctime:    now,
		Utime:    now,
	})
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			const uniqueConflictsErrNo = 1062
			if mysqlErr.Number == uniqueConflictsErrNo {
				// 邮箱冲突或者手机号冲突
				return nil, xerrors.New(400, "手机号已被注册")
			}
		}
		return nil, err
	}
	uid, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{
		Uid: uid,
	}, nil
}
