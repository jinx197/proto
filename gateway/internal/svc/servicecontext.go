// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"gateway/internal/config"

	"github.com/zeromicro/go-zero/zrpc"
	"userservice/userclient"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userclient.User // RPC 客户端
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)), // 连接 userservice
	}
}
