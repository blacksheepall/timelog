# TimeLog 开发指南

这是本项目面向开发者的主文档。中文版本是事实来源；英文版本 `docs/development.md` 由 AI 翻译维护。

## 项目概览

TimeLog 是一个本地优先的全栈时间记录应用：

- 后端：Go，入口在 `main.go`。
- 前端：Vue 3 + Vite，目录是 `web/`。
- 数据库：SQLite，迁移文件在 `model/migrations/`。
- API：HTTP 路由在 `router/`，业务逻辑在 `service/`，持久化在 `model/`。
- API 合约：DTO 来源是 `api/proto/timelog/v1/`，REST 路径来源是 `api/openapi/rest.yaml`。
- 文档：非 `prod` 构建会暴露 Redoc，路径是 `/docs/redoc.html`，旧入口 `/swagger` 会重定向过去。
- MCP：MCP server 在 `mcp/`，说明见 `docs/mcp/`。

生产和本地二进制运行时，Go 会从 `main.go` embed `web/dist` 并同时服务 API 和 SPA。因此第一次编译 Go 主程序前，必须先生成 `web/dist`。

## 第一次拉仓库

需要先准备这些工具：

- Go：版本以 `go.mod` 为准。
- Node.js + pnpm：前端和 Buf CLI 都依赖 `web/` 的 dev dependencies。
- SQLite migration CLI：通过 `make install-deps` 安装。

推荐初始化流程：

```bash
cp config-example.yml config.yml
cp config-example.yml config-test.yml
make gen-certs
make install-deps
cd web && pnpm install && cd ..
make gen-api
make migrate env=dev
make buildx
```

说明：

- `config.yml` 和 `config-test.yml` 被 git ignore，本地必须自己创建。
- `config-example.yml` 默认启用 HTTPS，并引用 `certs/cert.pem` / `certs/key.pem`。如果沿用默认配置，需要先跑 `make gen-certs`；如果只做普通 HTTP 开发，也可以把本地 `config.yml` 的 `server.https_enabled` 改成 `false`。
- 默认前端代理和 passkey RP origin 使用 `timelog.local`。本地开发时建议把 `timelog.local` 解析到 `127.0.0.1`，例如在 `/etc/hosts` 增加 `127.0.0.1 timelog.local`。
- `make install-deps` 安装 Go 侧开发工具，包括 `protoc-gen-jsonschema`、`gentool` 和带 SQLite 驱动的 `migrate`。
- Buf CLI 来自 `web/node_modules/.bin/buf`，所以新仓库需要先安装前端依赖；`make gen-api` 在缺少 Buf 时也会尝试执行 `cd web && pnpm install`。
- `Makefile` 的 `env` 默认是 `prod`。开发库迁移请显式使用 `make migrate env=dev`，否则会操作 `prod.db`。
- `make buildx` 会先构建前端，再编译 Go 主程序，是最稳妥的完整构建入口。

## 日常开发

### 后端 + 前端开发模式

常见开发方式是后端跑在 `8080`，前端 Vite dev server 跑在 `5173` 并代理 `/api`：

```bash
make buildx
./main
```

另开终端：

```bash
cd web
pnpm run dev
```

访问 `http://localhost:5173`。默认 Vite 代理指向 `https://timelog.local:8080`，如果代理失败，先确认 `timelog.local` 能解析到本机并且后端 HTTPS 证书已生成。

如果开启 passkey，推荐使用完整构建后的后端 HTTPS 地址访问，例如 `https://timelog.local:8080`。Passkey 的 RP origin 必须和 `config.yml` 中的 `passkey.rp_origins` 匹配。

### 构建

```bash
make buildx        # 完整本机构建：前端 + API 生成 + Go 二进制
make buildx-linux  # Linux 发布构建
make web           # 只构建前端，同时会运行 gen-api
make build         # 编译 Go 主程序；要求 web/dist 已存在
```

`make build` 会运行 `gen-api`，但不会构建 `web/dist`。如果是干净仓库，优先用 `make buildx`。

### 测试与检查

按改动范围选择最小必要检查：

```bash
go test ./router/...
go test ./service/...
go test ./model/...
go test ./test/...
make test
```

前端改动：

```bash
cd web
pnpm run type-check
pnpm run build
```

API 合约改动：

```bash
make gen-api
make check-api
make test
```

格式化：

```bash
make fmt
```

## 生成文件规则

本项目选择不把主要生成物提交到 git。这样可以让 review 聚焦在源文件上，避免生成代码占满 diff。

绝对不要手改这些路径：

- `gen/go/`
- `web/src/gen/`
- `gen/openapi/schemas/`
- `router/docs/openapi.yaml`
- `model/gen/`
- `web/dist/`

应该修改对应的源文件：

| 目标 | 修改哪里 | 生成命令 |
| --- | --- | --- |
| Go / TypeScript API DTO | `api/proto/timelog/v1/*.proto` | `make gen-api` |
| REST 路径、operation、响应 envelope 说明 | `api/openapi/rest.yaml` | `make gen-api` |
| 合并 OpenAPI 的逻辑 | `cmd/merge-openapi/` | `make gen-api` |
| Redoc 静态页 | `router/docs/redoc.html` | 不需要生成 |
| GORM model 代码 | 数据库结构和 `model/gentool/gormgen.go` | `make gen-model` |
| 前端生产静态文件 | `web/src/` | `make web` 或 `make buildx` |

`router/docs/openapi.yaml` 是由 `api/openapi/rest.yaml` 和 `gen/openapi/schemas/` 合并出来的。非 `prod` 构建会 embed 它，所以如果缺失，运行 `make gen-api`。

## 改代码时写在哪里

### API 与后端

- 新增或修改 API DTO：改 `api/proto/timelog/v1/`。
- 新增或修改 REST 路径：改 `api/openapi/rest.yaml`。
- 路由注册和 HTTP handler：改 `router/`。
- 请求 / 响应 envelope：Go 侧在 `router/apiresponse.go`，前端类型在 `web/src/types/api.ts`。
- 业务逻辑：改 `service/`。
- 数据访问和模型：改 `model/`。
- API DTO 与数据库模型之间的转换：改 `internal/api/mapper/`。

保持现有分层：`router -> service -> model`。不要把数据库访问直接塞进 handler，也不要在 service 里绕开 model 层。

### 前端

- API 调用集中在 `web/src/api/index.ts`。
- 生成的 DTO 从 `web/src/gen/` 导入，但不要手改该目录。
- 路由标题来自 `web/src/router/index.ts` 的 `meta.title`。
- 样式使用现有 Tailwind / Element Plus 风格。
- 前端仍保留 `tagAPI` 作为 `categoryAPI` 的兼容别名，新代码不要再引入 tag 语义。

前端改动前先读 `web/CLAUDE.md`，但以实际代码和本开发指南为准；该文件里有部分历史描述可能不再完全准确。

### 数据库迁移

创建迁移：

```bash
migrate -database "sqlite3://dev.db" create -seq -ext sql --dir model/migrations/ add_xxx
```

执行开发库迁移：

```bash
make migrate env=dev
```

执行生产库迁移：

```bash
make migrate env=prod
```

`make migrate` 不带参数会使用默认 `env=prod`。开发时请明确写 `env=dev`。

## 提交前检查

提交前按改动类型检查：

- 只改后端业务逻辑：先跑相关包测试，再视影响范围跑 `make test`。
- 改 API 合约、REST spec 或 mapper：跑 `make gen-api && make check-api && make test`，前端有使用时再跑 `cd web && pnpm run type-check`。
- 改前端：跑 `cd web && pnpm run type-check`，必要时跑 `cd web && pnpm run build`。
- 改数据库 schema / migration：跑 `make migrate env=dev`，必要时跑相关 model/service 测试。
- 改发布构建或 embed 相关逻辑：跑 `make buildx`。

如果生成物出现在 `git status` 里，先确认它是否应该被 ignore。通常应该提交的是源文件，而不是生成输出。

## 常见问题

### 启动时报缺少 `config.yml`

从样例复制：

```bash
cp config-example.yml config.yml
```

测试需要：

```bash
cp config-example.yml config-test.yml
```

### `make gen-api` 找不到 Buf

先安装前端依赖：

```bash
cd web
pnpm install
```

### `make migrate` 找不到 `migrate`

运行：

```bash
make install-deps
```

并确认 Go 安装目录的 `bin` 在 `PATH` 里。

### 非 `prod` 构建找不到 `router/docs/openapi.yaml`

运行：

```bash
make gen-api
```

### Go build 提示找不到 `web/dist`

运行完整构建：

```bash
make buildx
```

`make build` 只编译 Go，不负责生成前端 dist。

### Passkey 本地不可用

Passkey 需要 HTTPS secure context。先生成证书：

```bash
make gen-certs
```

然后按 README 更新 `config.yml`，并在浏览器里接受本地自签证书。

默认 passkey 配置使用 `timelog.local`，因此还需要确保它解析到本机，并通过 `https://timelog.local:8080` 这类匹配 RP origin 的地址访问应用。
