// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"helloapi/internal/model"

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

func (l *GetDeptCountLogic) GetDeptCount(req *types.DeptCountRequest) (*types.DeptCountResponse, error) {
	var count int64
	l.svcCtx.DB.Model(&model.Employee{}).
		Where("dept = ?", req.Dept).
		Count(&count)

	return &types.DeptCountResponse{
		Dept:  req.Dept,
		Count: int(count),
	}, nil
}
