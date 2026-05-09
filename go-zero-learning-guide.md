# go-zero 微服务学习手册

> 目标：入门 → 中级 | 形式：文档 + 手动实践 | 环境：Go + etcd 已就绪

***

## 一、环境检查清单

在开始之前，确认你的环境：

```bash
# 1. Go 版本
go version
# 预期：go1.21 或更高

# 2. etcd 是否运行
etcd --version
# 预期：etcd version 3.5.x

# 3. 启动 etcd（如果没运行）
etcd &
# 或指定端口
etcd --listen-client-urls=http://localhost:2379 --advertise-client-urls=http://localhost:2379

# 4. 验证 etcd
etcdctl endpoint health
# 预期：{"health":true}
```

***

## 二、goctl 命令速查表

goctl 是 go-zero 的代码生成工具，贯穿整个开发流程。

### 2.1 核心命令分类

| 命令 | 用途 | 示例 |
|------|------|------|
| **API 开发** | | |
| `goctl api new` | 创建新 API 服务 | `goctl api new userservice` |
| `goctl api go` | 根据 .api 生成 Go 代码 | `goctl api go --api helloapi.api --dir .` |
| `goctl api format` | 格式化 .api 文件 | `goctl api format --dir .` |
| `goctl api validate` | 校验 .api 文件语法 | `goctl api validate --api helloapi.api` |
| `goctl api doc` | 生成 Markdown 文档 | `goctl api doc --dir . --o ./docs` |
| `goctl api swagger` | 生成 Swagger 文档（v1.8.2+） | `goctl api swagger --api helloapi.api --dir ./swagger` |
| `goctl api ts` | 生成 TypeScript 代码 | `goctl api ts --api helloapi.api --dir ./ts` |
| **gRPC 开发** | | |
| `goctl rpc new` | 创建 RPC 演示服务 | `goctl rpc new greet` |
| `goctl rpc protoc` | 根据 proto 生成 gRPC 代码 | `goctl rpc protoc greet.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.` |
| `goctl rpc template` | 生成 proto 模板文件 | `goctl rpc -o greet.proto` |
| **Model / Database** | | |
| `goctl model mysql ddl` | 从 SQL 文件生成 model | `goctl model mysql ddl -src="./sql/*.sql" -dir=./model` |
| `goctl model mysql datasource` | 从数据库连接生成 model | `goctl model mysql datasource -url="user:pwd@tcp(127.0.0.1:3306)/db" -table="user" -dir=./model` |
| `goctl model pg datasource` | 从 PostgreSQL 生成 model | `goctl model pg datasource --url="postgres://..." --table="users" --dir=./model` |
| `goctl model mongo` | MongoDB 转 Go 模型 | `goctl model mongo --type User --dir ./model` |
| **运维相关** | | |
| `goctl docker` | 生成 Dockerfile | `goctl docker --go helloapi.go --port 8080 --version 1.21` |
| `goctl kube` | 生成 Kubernetes 部署文件 | `goctl kube deploy -name myapp -namespace default -image myapp:latest -o deploy.yaml -port 8080` |
| **模板管理** | | |
| `goctl template init` | 初始化所有模板（强制更新） | `goctl template init` |
| `goctl template update` | 更新指定类别模板 | `goctl template update -c api` |
| `goctl template revert` | 还原指定模板到原始版本 | `goctl template revert -c api -n handler.tpl` |
| `goctl template clean` | 清理缓存模板 | `goctl template clean` |
| **环境与工具** | | |
| `goctl env` | 查看环境变量 | `goctl env` |
| `goctl env check` | 检测依赖工具并安装 | `goctl env check --install --force` |
| `goctl upgrade` | 升级 goctl 到最新版本 | `goctl upgrade` |

### 2.2 常用参数

```bash
# --verbose 或 -v  输出详细日志（rpc protoc 支持）
goctl rpc protoc user.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. -v

# --home 和 --remote  指定模板源（不能同时使用，--remote 优先级更高）
goctl api new userservice --home /path/to/templates --remote https://github.com/zeromicro/go-zero-template.git

# --dir 或 -d  指定输出目录
goctl model mysql datasource -url="user:pwd@tcp(127.0.0.1:3306)/db" -table="user" -dir=./model
```

### 2.3 快速验证 goctl 安装

```bash
goctl --version
# 预期：goctl version 1.6.x 或更高
```

如果未安装：

```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

***

## 三、阶段一：单服务 Hello World

### 3.1 创建 API 服务

```bash
# 创建项目目录
mkdir go-zero-learning && cd go-zero-learning

# 生成 API 服务模板（注意：服务名不能有连字符，使用下划线或驼峰）
goctl api new helloapi
cd helloapi
```

### 3.2 生成的代码结构

```
helloapi/
├── etc/
│   └── helloapi-api.yaml        # 服务配置（包含 etcd 注册信息）
├── internal/
│   ├── config/
│   │   └── config.go          # 配置定义
│   ├── handler/
│   │   └── hellohandler.go    # 路由处理逻辑
│   ├── logic/
│   │   └── hellologic.go       # 业务逻辑
│   ├── svc/
│   │   └── servicecontext.go   # 服务上下文（依赖注入）
│   └── handler.go             # 路由注册
├── helloapi.api              # API 定义文件
├── go.mod
├── helloapi.go              # 入口文件（注意：不是 main.go）
└── go.mod
```

### 3.3 查看 etc 配置

```bash
cat etc/helloapi-api.yaml
```

你会看到类似：

```yaml
Name: helloapi-api
Host: 0.0.0.0
Port: 8080

# etcd 配置 - 重点关注这部分
Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: helloapi-api
```

### 3.4 启动服务（注册到 etcd）

```bash
go run helloapi.go -f etc/helloapi-api.yaml
```

### 3.5 验证服务注册到 etcd

```bash
# 查看 etcd 中的服务注册
etcdctl get --prefix ""
```

应该能看到 `helloapi-api` 的注册信息。

### 3.6 测试接口

```bash
curl http://localhost:8080/hello/world
# 预期响应：{"message":"Hello, go-zero!"}
```

### 3.7 预期结果

- [ ] 服务启动无报错
- [ ] curl 返回正确响应

***

## 三（补充）：API 业务开发指南

### 3.8 业务开发流程

一个 API 服务的开发流程：

```
定义 api 文件 → 生成代码 → 写业务逻辑 → 启动服务
```

**步骤一：定义 API 接口** — 编辑 `{服务名}.api`

```api
type (
    // 请求结构体
    DeptCountRequest {
        Dept string `json:"dept"`
    }

    // 响应结构体
    DeptCountResponse {
        Dept  string `json:"dept"`
        Count int    `json:"count"`
    }
)

service helloapi {
    @handler getDeptCountHandler
    get /api/dept/count (DeptCountRequest) returns (DeptCountResponse)
}
```

**步骤二：生成代码**

```bash
goctl api format --dir .
goctl api go --api helloapi.api --dir .
```

**步骤三：写业务逻辑** — 编辑 `internal/logic/getdeptcountlogic.go`

```go
package logic

import (
    "context"
    "helloapi/internal/svc"
    "helloapi/internal/types"
)

type GetDeptCountLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewGetDeptCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeptCountLogic {
    return &GetDeptCountLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *GetDeptCountLogic) GetDeptCount(req *types.DeptCountRequest) (*types.DeptCountResponse, error) {
    // 在这里写业务逻辑
    return &types.DeptCountResponse{
        Dept:  req.Dept,
        Count: 10,
    }, nil
}
```

**步骤四：启动测试**

```bash
go run helloapi.go -f etc/helloapi.yaml
curl "http://localhost:8080/api/dept/count?dept=研发部"
```

***

### 3.9 接入 gorm 操作数据库

**1. 安装 gorm 驱动**

```bash
go get gorm.io/gorm
go get gorm.io/driver/mysql
```

**2. 定义 Model** — 创建 `internal/model/employee.go`

```go
package model

type Employee struct {
    Id   int64  `gorm:"primaryKey"`
    Name string `gorm:"size:50"`
    Dept string `gorm:"size:50"`
}

func (Employee) TableName() string {
    return "employee"
}
```

**3. 配置 gorm** — 编辑 `internal/config/config.go`

```go
package config

import (
    "github.com/zeromicro/go-zero/rest"
    "gorm.io/gorm"
)

type Config struct {
    rest.RestConf
    MysqlDB struct {
        DataSource string
    }
}
```

**4. 初始化连接** — 编辑 `internal/svc/servicecontext.go`

```go
package svc

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "helloapi/internal/config"
)

type ServiceContext struct {
    Config config.Config
    DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
    db, err := gorm.Open(mysql.Open(c.MysqlDB.DataSource), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    return &ServiceContext{
        Config: c,
        DB:     db,
    }
}
```

**5. 配置数据源** — 编辑 `etc/helloapi.yaml`

```yaml
Name: helloapi
Host: 0.0.0.0
Port: 8080

MysqlDB:
  DataSource: "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4"
```

**6. 在业务逻辑中使用** — 编辑 `internal/logic/getdeptcountlogic.go`

```go
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
```

***

### 3.10 新增接口的正确方式

**重要原则：`goctl api go`** **不会覆盖已存在的文件**

生成行为：

- `handler/*.go` — 文件已存在则跳过
- `logic/*.go` — 文件已存在则跳过
- `types/types.go` — 合并，新类型追加
- `handler.go` — 路由注册会更新

**新增接口步骤：**

1. 编辑 `helloapi.api`，追加新的接口定义

```api
type (
    AllDeptsResponse {
        Depts []string `json:"depts"`
    }
)

service helloapi {
    @handler getAllDeptsHandler
    get /api/dept/all returns (AllDeptsResponse)
}
```

1. 重新生成

```bash
goctl api format --dir .
goctl api go --api helloapi.api --dir .
```

1. 手动创建新的 logic 文件 `internal/logic/getalldeptslogic.go`（不会被覆盖）

**建议：用 git 管理代码，生成前先提交，这样有问题可以回滚。**

***

### 3.11 多文件拆分维护

当接口很多时，可以按业务拆分成多个 .api 文件，用 import 合并。

**目录结构：**

```
helloapi/
├── api/
│   ├── main.api           # 入口文件，用 import 汇总
│   ├── employee.api       # 员工相关接口
│   ├── department.api     # 部门相关接口
│   └── project.api        # 项目相关接口
├── internal/
│   └── ...
└── helloapi.go
```

**main.api（入口文件）：**

```api
syntax = "v1"

type request {
}

type response {
}

// 引入其他 api 文件
import "./api/employee.api"
import "./api/department.api"
import "./api/project.api"
```

**employee.api：**

```api
type (
    EmployeeRequest {
        Id int64 `json:"id"`
    }

    EmployeeResponse {
        Id   int64  `json:"id"`
        Name string `json:"name"`
    }
)

service helloapi {
    @handler getEmployeeHandler
    get /api/employee/get (EmployeeRequest) returns (EmployeeResponse)
}
```

**department.api：**

```api
type (
    DeptCountRequest {
        Dept string `json:"dept"`
    }

    DeptCountResponse {
        Dept  string `json:"dept"`
        Count int    `json:"count"`
    }
)

service helloapi {
    @handler getDeptCountHandler
    get /api/dept/count (DeptCountRequest) returns (DeptCountResponse)
}
```

**生成时用主文件：**

```bash
goctl api format --dir api
goctl api go --api api/main.api --dir .
```

**好处：每个业务模块一个文件，方便维护，不会在合并代码时冲突。**

***

## 四、阶段二：两个微服务互相调用

### 4.1 服务架构

```
                    ┌──────────────┐
    curl ──────►    │  gateway  │    ────grpc─────►  │  userservice  │
                    │   (端口8080)   │                      │   (端口8081)    │
                    └──────────────┘                      └───────────────┘
                         etcd                                  etcd
                     注册发现                                 注册发现
```

### 4.2 创建用户服务（userservice）

```bash
cd go-zero-learning
goctl api new userservice
cd userservice
```

### 4.3 修改 userservice 的 API 定义

编辑 `userservice/userservice.api`：

```api
type (
    UserRequest {
        Id int64 `json:"id"`
    }

    UserResponse {
        Id   int64  `json:"id"`
        Name string `json:"name"`
        Age  int    `json:"age"`
    }
)

service userservice {
    @handler getUserHandler
    psot /api/user/getUser (UserRequest) returns (UserResponse)
}
```

### 4.4 修改 userservice 配置

编辑 `userservice/etc/userservice.yaml`：

```yaml
Name: userservice
Host: 0.0.0.0
Port: 8081   # 改用 8081，避免端口冲突

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: userservice
```

### 4.5 重新生成代码

```bash
cd userservice
goctl api format --dir .
goctl api go --api userservice.api --dir .
```

### 4.6 实现业务逻辑

编辑 `internal/logic/getuserlogic.go`：

```go
package logic

import (
    "context"

    "userservice/internal/svc"
    "userservice/pb/user"

    "github.com/zeromicro/go-zero/zrpc"
)

type GetUserLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
    return &GetUserLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *GetUserLogic) GetUser(req *user.UserRequest) (*user.UserResponse, error) {
    // 模拟数据库查询
    return &user.UserResponse{
        Id:   req.Id,
        Name: "ZhangSan",
        Age:  25,
    }, nil
}
```

### 4.7 启动 userservice

```bash
go run userservice.go -f etc/userservice.yaml
```

验证：

```bash
etcdctl get --prefix ""
curl curl http://localhost:8081/api/user/getUser -X POST \
    -H "Content-Type: application/json" \
    -d '{"id":1}'
```

### 4.8 创建网关服务（gateway）

```bash
cd go-zero-learning
goctl api new gateway
cd gateway
```

### 4.9 修改 gateway 配置

编辑 `gateway/etc/gateway.yaml`：

```yaml
Name: gateway
Host: 0.0.0.0
Port: 8080   # 对外端口 8080

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: gateway
```

### 4.10 定义调用 userservice 的路由

编辑 `gateway/gateway.api`：

```api
type (
    UserRequest {
        Id int64 `json:"id"`
    }

    UserResponse {
        Id   int64  `json:"id"`
        Name string `json:"name"`
        Age  int    `json:"age"`
    }
)

service gateway {
    @handler getUserHandler
    get /api/user/getUser (UserRequest) returns (UserResponse)
}
```

### 4.11 配置 gRPC 客户端调用 userservice

编辑 `gateway/internal/config/config.go`，添加 rpc 配置：

```go
package config

import (
    "github.com/zeromicro/go-zero/rest"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
    rest.RestConf                 // API 服务配置
    UserRpc zrpc.RpcClientConf    // RPC 客户端配置（调用 userservice）
}
```

### 4.12 生成 protobuf 代码

首先确保项目根目录有 `go.work`，让各服务可以互相引用：

```bash
cd go-zero-learning
go work init
go work use ./helloapi ./userservice ./gateway
```

> 因为 helloapi、userservice、gateway 是单独的 Go module，需要 go workspace 才能跨模块 import。

删除 userservice 中旧的 API 相关文件（RPC 服务不再需要这些）：

```bash
cd userservice
rm -rf internal/handler internal/types userservice.go
```

创建 proto 文件 `userservice/pb/user/user.proto`：

```protobuf
syntax = "proto3";

package user;

option go_package = "userservice/pb/user";

message UserRequest {
    int64 id = 1;
}

message UserResponse {
    int64 id = 1;
    string name = 2;
    int32 age = 3;
}

service User {
    rpc GetUser (UserRequest) returns (UserResponse);
}
```

> 注意：`go_package` 要写成模块全路径 `userservice/pb/user`，不要用 `./pb/user`，否则生成的 pb 文件会嵌套到 `pb/pb/user/` 目录。

生成 Go 代码：

```bash
cd userservice
goctl rpc protoc pb/user/user.proto \
    --go_out=. --go-grpc_out=. \
    --go_opt=module=userservice --go-grpc_opt=module=userservice \
    --zrpc_out=.
```

> 说明：`--go_opt=module=userservice` 会让 protoc 根据 go_package 路径剥离模块前缀，将文件正确输出到 `pb/user/`（而不是 `userservice/pb/user/`）。

goctl rpc protoc 会自动生成以下文件：

```
internal/config/config.go       # 覆盖：改用 RpcServerConf
internal/server/userserver.go   # RPC server 桩，调用 logic
user.go                         # RPC 主入口（package main）
userclient/user.go              # RPC 客户端代理
pb/user/user.pb.go              # protobuf 生成
pb/user/user_grpc.pb.go         # gRPC 生成
```

### 4.13 手动修改：config 和 logic

goctl 生成骨架代码后，有两个文件需要手动修改：

**4.13.1 修改 config.go**

`internal/config/config.go` 原本是 API 服务的 `rest.RestConf`，需要改成 RPC 的 `zrpc.RpcServerConf`：

```go
package config

import (
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
    zrpc.RpcServerConf
}
```

**4.13.2 修改 getuserlogic.go**

`internal/logic/getuserlogic.go` 原本引用的是 API 的 `types.UserRequest/UserResponse`，需要改成 protobuf 类型：

```go
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
```

### 4.14 验证 RPC 配置和入口

goctl rpc protoc 已自动生成 `user.go`（RPC 入口）和 `etc/user.yaml`（RPC 配置模板）。查看生成的文件确认无误：

`user.go` 关键部分（已自动生成，无需手动编写）：

```go
var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
    // ...
    s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
        user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))
        // ...
    })
    s.Start()
}
```

`etc/user.yaml` 配置模板：

```yaml
Name: user.rpc
ListenOn: 0.0.0.0:8080
Etcd:
  Hosts:
  - 127.0.0.1:2379
  Key: user.rpc
```

根据实际需要调整端口和 etcd 配置。

### 4.15 创建 RPC 客户端代理

在 `gateway/internal/svc/servicecontext.go` 中：

```go
package svc

import (
    "github.com/zeromicro/go-zero/zrpc"
    "gateway/internal/config"
    "userservice/pb/user"
)

type ServiceContext struct {
    Config  config.Config
    UserRpc user.UserClient   // RPC 客户端
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:  c,
        UserRpc: user.NewUserClient(zrpc.MustNewClient(c.UserRpc)),  // 连接 userservice
    }
}
```

### 4.16 修改 handler 调用 RPC

编辑 `gateway/internal/handler/getuserhandler.go`：

```go
package handler

import (
    "net/http"

    "github.com/zeromicro/go-zero/rest"
    "gateway/internal/logic"
    "gateway/internal/svc"
    "gateway/internal/types"
)

func getUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.UserRequest
        if err := httpx.Parse(r, &req); err != nil {
            resp := &types.UserResponse{}
            rest.Response(w, r, resp)
            return
        }

        l := logic.NewGetUserLogic(r.Context(), svcCtx)
        // 通过 RPC 调用 userservice
        rpcReq := &user.UserRequest{Id: req.Id}
        rpcResp, err := l.svcCtx.UserRpc.GetUser(r.Context(), rpcReq)
        if err != nil {
            // 处理错误
            resp := &types.UserResponse{}
            rest.Response(w, r, resp)
            return
        }

        // 返回结果
        resp := &types.UserResponse{
            Id:   rpcResp.Id,
            Name: rpcResp.Name,
            Age:  rpcResp.Age,
        }
        rest.Response(w, r, resp)
    }
}
```

### 4.17 启动两个服务

终端1：启动 userservice

```bash
cd userservice
go run userservice.go -f etc/userservice.yaml
```

终端2：启动 gateway

```bash
cd gateway
go run gateway.go -f etc/gateway.yaml
```

### 4.18 测试完整调用链路

```bash
# 通过网关调用用户服务
curl http://localhost:8080/api/user/getUser?id=1

# 预期响应
{"id":1,"name":"ZhangSan","age":25}
```

### 4.19 验证 etcd 注册

```bash
etcdctl get --prefix ""
```

应该能看到 gateway 和 userservice 两个服务。

### 4.20 预期结果清单

- [ ] userservice 启动在 8081 端口
- [ ] gateway 启动在 8080 端口
- [ ] `etcdctl get` 能查到两个服务
- [ ] curl 网关返回正确的用户信息
- [ ] 整个链路：curl → gateway → etcd 发现 → userservice → 返回

***

## 五、阶段三：中级特性

### 5.1 中间件

go-zero 的中间件分两种：

**全局中间件**（所有路由）：

```go
// internal/handler.go
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    // 全局中间件 - 记录每个请求
    server.AddRoutes([]rest.Route{
        {
            Method:  "GET",
            Path:    "/",
            Handler: indexHandler(serverCtx),
        },
        // ...
    })

    // 添加全局中间件
    server.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // 记录开始时间
            start := time.Now()
            next.ServeHTTP(w, r)
            // 记录结束时间
            log.Printf("请求耗时: %v", time.Since(start))
        })
    })
}
```

**单个路由中间件**：

```go
server.AddRoutes([]rest.Route{
    {
        Method:  "GET",
        Path:    "/api/user/getUser",
        Handler: getUserHandler(serverCtx),
        Middlewares: []rest.Middleware{authMiddleware},  // 仅此路由
    },
})
```

### 5.2 JWT 认证

**1. 定义 JWT 配置**（`internal/config/config.go`）：

```go
type Config struct {
    rest.RestConf
    Auth struct {
        AccessSecret string
        AccessExpire int64
    }
}
```

**2. 登录时生成 Token**（`internal/logic/loginlogic.go`）：

```go
import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
    // 验证用户名密码（略）

    // 生成 JWT
    now := time.Now().Unix()
    claims := jwt.MapClaims{
        "userId": userId,
        "exp":    now + l.svcCtx.Config.Auth.AccessExpire,
    }
    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
        SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
    if err != nil {
        return nil, err
    }

    return &types.LoginResp{
        Token: token,
    }, nil
}
```

**3. 配置 JWT 校验** — go-zero 内置了 JWT 中间件，在路由注册时添加：

```go
// internal/handler/routes.go
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    server.AddRoutes(
        []rest.Route{
            {
                Method:  http.MethodGet,
                Path:    "/api/user/info",
                Handler: getUserInfoHandler(serverCtx),
            },
        },
        rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
    )
}
```

go-zero 的 `rest.WithJwt` 会自动从 Authorization header 解析和验证 JWT token，无需手动编写中间件。

**4. 应用中间件**：

```go
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    // 需要认证的路由
    server.AddRoutes([]rest.Route{
        {
            Method:  "GET",
            Path:    "/api/user/info",
            Handler: getUserInfoHandler(serverCtx),
            Middlewares: []rest.Middleware{AuthMiddleware},
        },
    }, serverCtx)
}
```

### 5.3 超时和熔断配置

**超时配置**（`yaml`）：

```yaml
Name: gateway
Host: 0.0.0.0
Port: 8080

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: gateway

# RPC 客户端超时配置
UserRpc:
  Endpoints:
    - 127.0.0.1:8081
  Timeout: 3000  # 超时时间 3 秒
  # 或者用整段配置
  # Etcd:
  #   Hosts:
  #     - 127.0.0.1:2379
  #   Key: userservice
```

**熔断配置**：

```yaml
UserRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: userservice
  CircuitBreaker:
    Enabled: true
    Qps: 100           # QPS 阈值
    RequestAmount: 10  # 熔断器打开的最小请求数
    Sleep: 1000        # 熔断持续时间（毫秒）
```

在代码中配置：

```go
UserRpc: zrpc.RpcClientConf{
    Endpoints: []string{"127.0.0.1:8081"},
    Timeout: 3000,
    CircuitBreaker: zrpc.CircuitBreaker{
        Enabled: true,
    },
}
```

### 5.4 请求参数校验

**定义校验规则**（`api` 文件）：

```api
type (
    UserRequest {
        Id   int64  `json:"id,options=required"`        // 必填
        Name string `json:"name,range=[1:20]"`         // 长度 1-20
        Age  int    `json:"age,optional"`              // 可选
        Email string `json:"email,optional"`           // 可选
    }
)
```

**启用校验**：

```go
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    server.AddRoutes(rest.Route{
        Method:  "POST",
        Path:    "/api/user/create",
        Handler: createUserHandler(serverCtx),
        Validator: func(r *http.Request, req interface{}) error {
            // 内置的 validator 可以自动校验
            // 这里可以添加自定义校验
            return nil
        },
    })
}
```

**自定义校验器**：

```go
type LoginReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func ValidateLogin(req *LoginReq) error {
    if len(req.Username) < 3 {
        return errors.New("username must be at least 3 characters")
    }
    if len(req.Password) < 6 {
        return errors.New("password must be at least 6 characters")
    }
    return nil
}
```

### 5.5 限流配置

**方式一：使用 go-zero 内置的并发控制**

在 yaml 配置中直接设置最大并发连接数：

```yaml
MaxConns: 1000  # 最大并发连接数
```

**方式二：自定义限流中间件**

使用 `github.com/zeromicro/go-zero/core/limit` 包：

```go
// internal/middleware/ratelimit.go
import "github.com/zeromicro/go-zero/core/limit"

func RateLimitMiddleware(rate int) func(http.Handler) http.Handler {
    limiter := limit.NewTokenLimiter(rate, rate, time.Second)
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                w.WriteHeader(http.StatusTooManyRequests)
                w.Write([]byte(`{"error":"rate limit exceeded"}`))
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

在路由中使用：

```go
server.AddRoutes([]rest.Route{
    {
        Method:  "GET",
        Path:    "/api/user/list",
        Handler: listUserHandler(serverCtx),
        Middlewares: []rest.Middleware{RateLimitMiddleware(100)},
    },
})
```

### 5.6 日志配置

```yaml
Log:
  ServiceName: gateway
  Mode: console  # console 或 file
  Level: info    # debug, info, error
  Path: /var/log/gateway  # 当 mode=file 时
```

代码中：

```go
import "github.com/zeromicro/go-zero/core/logx"

func (l *Logic) DoSomething() {
    logx.Info("这是一条普通日志")
    logx.Errorf("出错了: %v", err)
    logx.Debug("调试信息")
}
```

***

## 六、常见错误对照表

| 错误现象                          | 可能原因        | 解决方法                                                         |
| ----------------------------- | ----------- | ------------------------------------------------------------ |
| `etcd cluster is unavailable` | etcd 未启动    | `etcd &` 启动 etcd                                             |
| `connection refused`          | 服务未注册到 etcd | 检查 yaml 配置是否正确                                               |
| `context deadline exceeded`   | RPC 超时      | 增加 Timeout 配置                                                |
| `circuit breaker open`        | 熔断触发        | 检查下游服务，重试                                                    |
| `401 unauthorized`            | JWT 未传或已过期  | 检查 token 生成和传递逻辑                                             |
| `port already in use`         | 端口被占用       | 改用其他端口                                                       |
| `goctl: command not found`    | 未安装 goctl   | `go install github.com/zeromicro/go-zero/tools/goctl@latest` |

***

## 七、阶段验收清单

### 阶段一完成标志

- [ ] 能用 `goctl api new` 创建 API 服务
- [ ] 能启动服务并注册到 etcd
- [ ] 能 curl 到正确的响应

### 阶段二完成标志

- [ ] 两个服务能独立运行
- [ ] etcd 能查到两个服务注册
- [ ] 通过网关能调到用户服务返回数据

### 阶段三完成标志

- [ ] 能用 JWT 做登录认证
- [ ] 能配置超时和熔断
- [ ] 能用中间件做限流和日志

***

## 八、下一步学习建议

**继续深入的方向：**

1. **go-zero 源码分析**：研究 core 模块的实现
2. **数据库集成**：go-zero + gorm/MySQL 操作
3. **Redis 缓存**：go-zero 的 redis 客户端
4. **消息队列**：go-zero + kafka/rocketmq
5. **项目实战**：用一个完整项目串联所有知识点

**简历写法示例：**

```
项目经历：
- 使用 go-zero 搭建微服务框架，实现了 2 个服务（API网关 + 用户服务）
- 基于 etcd 实现服务注册与发现，支持服务横向扩展
- 实现了 JWT 认证、接口限流、熔断等中间件
- 掌握 gRPC 通信，理解 Protobuf 编解码
```

***

*手册版本：v1.0 | 最后更新：2026-05-06*
