// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"helloapi/internal/svc"
	"helloapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeptCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeptCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeptCountLogic {
	return &GetDeptCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeptCountLogic) GetDeptCount(req *types.DeptCountRequest) (resp *types.DeptCountResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
