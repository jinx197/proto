# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

这是一个 go-zero 微服务框架学习项目，包含学习手册和一个示例 API 服务。

## 环境要求

- Go 1.21+
- etcd 3.5.x（服务注册与发现）
- goctl 1.6.x+（代码生成工具）

```bash
# 安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 启动 etcd
etcd --listen-client-urls=http://localhost:2379 --advertise-client-urls=http://localhost:2379

# 验证 etcd
etcdctl endpoint health
```

## 常用命令

```bash
# 启动服务（指定配置文件）
go run helloapi.go -f etc/helloapi-api.yaml

# 重新生成代码（当 .api 文件变更后）
goctl api format --dir .
goctl api go --api helloapi.api --dir .

# 查看 etcd 注册的服务
etcdctl get --prefix /registry/services

# 测试接口
curl http://localhost:8080/hello/world
curl "http://localhost:8080/api/dept/count?dept=研发部"
```

## 项目结构

```
helloapi/
├── etc/
│   └── helloapi-api.yaml        # 服务配置（包含 etcd 注册信息）
├── internal/
│   ├── config/config.go      # 配置定义（嵌入 rest.RestConf）
│   ├── handler/              # 路由处理逻辑（生成，勿手动编辑）
│   ├── logic/                # 业务逻辑（在此写核心代码）
│   ├── svc/servicecontext.go # 服务上下文（依赖注入）
│   └── types/types.go        # 请求/响应结构体
├── helloapi.api              # API 定义文件（核心入口）
└── helloapi.go              # 程序入口
```

## go-zero API 开发流程

1. **编辑 .api 文件** - 定义 type 和 service
2. **生成代码** - `goctl api format --dir . && goctl api go --api helloapi.api --dir .`
3. **编写业务逻辑** - 编辑 `internal/logic/*logic.go`
4. **启动服务** - `go run helloapi.go -f etc/helloapi.yaml`

注意：`goctl api go` 不会覆盖已存在的 handler 和 logic 文件，但会更新 types 和路由注册。

## go-zero 架构要点

- **handler**: 路由处理入口，参数解析后调用 logic
- **logic**: 业务逻辑层，在此编写核心代码
- **svc**: 服务上下文，注入依赖（如数据库连接、RPC 客户端）
- **config**: 配置结构体，通过 yaml 文件加载
- **types**: 请求/响应结构体，由 goctl 根据 .api 自动生成

## 相关资源

- 学习手册：`go-zero-learning-guide.md`
- go-zero 官方文档：https://go-zero.dev/
- goctl 命令参考见学习手册第二章
