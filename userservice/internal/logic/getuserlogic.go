// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"userservice/internal/svc"
	"userservice/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(in *user.UserRequest) (*user.UserResponse, error) {
	// 模拟数据库查询
	return &user.UserResponse{
		Id:   in.Id,
		Name: "ZhangSan",
		Age:  25,
	}, nil
}
