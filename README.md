# taurus_go

`taurus_go` 是一个面向 Go 项目的基础工具库，聚合了代码生成、实体建模、统一通知、Redis 封装、日志、模板、错误码与常用数据处理能力。仓库的模块路径是 `github.com/zodileap/taurus_go`，当前 `go.mod` 要求 `Go 1.24.1`。

如果你需要的是：

- 基于 schema 的实体定义与代码生成
- 统一的消息通知抽象，目前已支持 Telegram
- 对 `go-redis` 的轻量封装和批量事务管道
- 文件轮转日志
- 模板、错误码、文件资产和字符串/切片/地理数据等工具函数

这个仓库就是对应的基础能力集合。

## 安装

```bash
go get github.com/zodileap/taurus_go@latest
```

如果需要固定版本：

```bash
go get github.com/zodileap/taurus_go@v0.9.23
```

## 示例与测试

`taurus_go` 不再维护单独的 `taurus_go_demo` 模块。公共能力示例直接放在各包的 `*_test.go` 和本 README 中，避免示例与实现分仓后逐步漂移。

像 `asset`、`cmd`、`cache/redis` 这类纯库能力，优先查看对应包下的测试文件；需要真实数据库、外部服务或生成产物配合的场景，则在 README 和源码注释中保留最小可用说明，把完整集成用例留在使用方仓库按实际环境编写。

## 适用范围

`taurus_go` 更像一个“基础能力仓库”，不是单一框架。不同子包的定位大致如下：

| 包路径 | 说明 |
| --- | --- |
| `grpc` | gRPC server/client 注册与连接管理 |
| `console` | CLI 输出辅助函数 |
| `entity` | 实体 schema DSL、连接配置、SQL 构造、代码生成 |
| `entity/cmd` | `entity` 的初始化和生成命令 |
| `notify` | 统一通知模型与发送接口 |
| `notify/telegram` | Telegram Bot API 的通知实现 |
| `cache/lru` | 泛型 LRU 缓存 |
| `cache/redis` | 基于 `go-redis/v9` 的客户端管理、事务管道与结果聚合 |
| `tlog` | 控制台输出 + 文件轮转日志 |
| `byteutil`、`maputil`、`sliceutil`、`stringutil`、`structutil`、`geo` | 字符串、切片、字节、结构体、地理数据等工具函数 |
| `asset` | 文件与目录资产写入、格式化、复制 |
| `template` | `text/template` 的轻量封装 |
| `err` | 统一错误码结构 |
| `security/rsa` | RSA 工具 |
| `cmd` | 对 shell / Go 命令执行的封装 |
| `rand`、`listnode` | 小型公共工具 |

## 快速上手

### 1. 发送 Telegram 通知

`notify` 定义了统一通知模型，`notify/telegram` 提供当前可直接使用的发送实现。

```go
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/zodileap/taurus_go/notify"
	"github.com/zodileap/taurus_go/notify/telegram"
)

func main() {
	ctx := context.Background()

	_, _ = telegram.Send(ctx, "telegram-bot-token", "@channelusername", notify.Notification{
		Title: "Deploy Ready",
		Body:  "build and smoke tests passed",
		Attachments: []notify.Attachment{
			{
				Kind:    notify.Document,
				URL:     "https://example.com/reports/deploy-report.pdf",
				Caption: "deploy report",
			},
		},
	})

	client := telegram.New(
		"telegram-bot-token",
		"123456789",
		telegram.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
	)

	_, _ = client.Send(ctx, notify.Notification{
		Body: "latest screenshot",
		Attachments: []notify.Attachment{
			{
				Kind: notify.Photo,
				Path: "/tmp/screenshot.png",
			},
		},
	})
}
```

当前统一通知抽象支持：

- 文本消息
- 图片、视频、音频、语音、文件、动画
- Telegram 媒体组

### 2. 使用 Redis 客户端和字符串操作

`cache/redis` 维护一个命名客户端池，并在 `Client` 上提供 `Set`、`Hash`、`String` 三类操作入口。既支持单次命令，也支持先累计操作再统一 `Save()`。

```go
package main

import (
	"fmt"
	"time"

	tredis "github.com/zodileap/taurus_go/cache/redis"
)

func main() {
	tredis.SetClient("default", &tredis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	defer tredis.ClearClient()

	client, err := tredis.GetClient("default")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	str, err := client.String()
	if err != nil {
		panic(err)
	}

	if err := str.AddR("app:version", 5*time.Minute, "v0.9.23"); err != nil {
		panic(err)
	}

	value, err := str.GetR("app:version")
	if err != nil {
		panic(err)
	}

	fmt.Println(value)
}
```

如果你希望把多次操作合并成一次事务管道执行，可以先记录操作，再调用 `Save()`：

```go
setClient, _ := client.Set()
setClient.Add("service:regions", 0, "ap-southeast-1")
setClient.Add("service:regions", 0, "eu-central-1")

res, err := client.Save()
if err != nil {
	panic(err)
}

fmt.Println(res.GetSet("service:regions").AddNum)
```

### 3. 使用 tlog 输出控制台和文件日志

`tlog` 支持命名 logger、等级控制、调用位置信息，以及带轮转的文件输出。

```go
package main

import "github.com/zodileap/taurus_go/tlog"

func main() {
	tlog.Get("app").
		SetLevel(tlog.InfoLevel).
		SetCaller(false).
		SetOutputPath("logs/app.log", 50, 10, 7)

	tlog.Info(
		"app",
		"service started",
		tlog.String("env", "dev"),
		tlog.Int("port", 8080),
	)
}
```

`SetOutputPath(path, maxSize, maxBackups, maxAge)` 的含义分别是：

- `path`: 日志文件路径
- `maxSize`: 单文件最大大小
- `maxBackups`: 轮转文件保留数量
- `maxAge`: 日志保留天数

### 4. 定义 entity schema 并生成代码

`entity` 是仓库里最核心的一块能力之一，包含：

- schema DSL
- 数据库连接注册
- SQL 构造器
- 代码生成器

先初始化一个 entity 目录：

```bash
go run github.com/zodileap/taurus_go/entity/cmd new Blog -e Post -t ./entity
```

这条命令会生成：

- `./entity/generate.go`
- `./entity/schema/db.go`
- `./entity/schema/post.go`

你可以在 schema 中定义自己的实体。一个最小示例如下：

```go
package schema

import (
	"github.com/zodileap/taurus_go/entity"
	"github.com/zodileap/taurus_go/entity/field"
)

type BlogEntity struct {
	entity.Entity
	Id          *field.Int64
	Title       *field.Varchar
	CreatedTime *field.Timestamptz
}

func (e *BlogEntity) Config() entity.EntityConfig {
	return entity.EntityConfig{
		AttrName: "blog",
		Comment:  "博客表",
	}
}

func (e *BlogEntity) Fields() []entity.FieldBuilder {
	return []entity.FieldBuilder{
		e.Id.Name("id").Primary(1).Sequence(entity.NewSequence("blog_id_seq")).Locked(),
		e.Title.Required().Name("title").MaxLen(128).Comment("标题"),
		e.CreatedTime.Default("CURRENT_TIMESTAMP").Name("created_time"),
	}
}
```

定义好 schema 后执行生成：

```bash
go run github.com/zodileap/taurus_go/entity/cmd generate ./entity/schema
```

或者直接使用 `go:generate`：

```go
package entity

//go:generate go run github.com/zodileap/taurus_go/entity/cmd generate ./schema
```

如果需要注册数据库连接，可以直接使用 `entity.AddConnection`：

```go
package main

import (
	"github.com/zodileap/taurus_go/entity"
	"github.com/zodileap/taurus_go/entity/dialect"
)

func main() {
	_ = entity.AddConnection(entity.ConnectionConfig{
		Driver:   dialect.PostgreSQL,
		Tag:      "main",
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "app",
	})
}
```

## 代码组织

仓库当前主要结构如下：

```text
taurus_go/
├── asset
├── cache/redis
├── cmd
├── byteutil
├── geo
├── maputil
├── sliceutil
├── stringutil
├── structutil
├── entity/{cmd,codegen,dialect,entitysql,field}
├── err
├── listnode
├── notify/telegram
├── rand
├── security/rsa
├── template
└── tlog
```

## 开发与测试

在仓库根目录执行：

```bash
go test ./...
```

如果你只想验证某个子系统，可以按包执行：

```bash
go test ./notify/...
go test ./entity/...
go test ./cache/redis
```

`entity` 相关开发常用命令：

```bash
go run github.com/zodileap/taurus_go/entity/cmd new User -e Profile -t ./entity
go run github.com/zodileap/taurus_go/entity/cmd generate ./entity/schema
```

## 版本与发布

仓库版本遵循 `Major.Minor.Patch` 三段式语义化版本。

- `Major`: 不兼容变更
- `Minor`: 新增功能
- `Patch`: 修复或小幅更新

如果你在项目中锁定版本，建议使用明确 tag，例如 `v0.9.23`。

## 说明

- README 中的示例基于当前仓库实际导出 API 整理，不依赖额外脚手架。
- `entity` 的生成链路较长，实际项目里通常会配合独立的业务仓库或 demo 工程使用。
- 更多仓库协作约束、文档规范和 AI 协作说明见 [AGENTS.md](./AGENTS.md)。
