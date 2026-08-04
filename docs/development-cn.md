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
make install-deps
cd web && pnpm install && cd ..
make gen-api
make migrate env=dev
make buildx
```

说明：

- `config.yml` 和 `config-test.yml` 被 git ignore，本地必须自己创建。
- `config-example.yml` 默认以纯 HTTP 运行。生产环境若需要 HTTPS（例如启用 passkey），请在应用前部署反向代理（Nginx/Traefik等）来终止 TLS。
- 本地开发时 Vite dev server 代理到 `http://127.0.0.1:8080`。
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

访问 `http://localhost:5173`。默认 Vite 代理指向 `http://127.0.0.1:8080`，如果代理失败，先确认后端服务已在 `8080` 端口运行。

如果开启 passkey，请确保浏览器访问的 origin 属于 `passkey.rp_origins`。生产环境通过反向代理提供 HTTPS 时，将 `passkey.rp_origins` 设为对应的公共 HTTPS origin；本地 Vite 开发可使用 `http://localhost:5173`（localhost 也是 secure context）。

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

## API 请求/响应约定

前后端共享同一份冻结契约，后端由 `router.TestEnvelopeContract` 与
`router/contract_test.go` 的 golden 测试强制约束：

- **信封**：所有 `/api` 响应（含 handler、auth 中间件、NoRoute 404）统一为
  `{data, message, status}`，定义在 `router/apiresponse.go`，前端镜像在
  `web/src/types/api.ts`。
- **成功**：HTTP 200，`data` 为 proto DTO（或 command 类接口的 `null`），
  `message` 为人类可读文案（默认 `"success"`），`status` 恒为 200。
- **错误**：非 2xx，`data` 恒为 `null`，`message` 为人类可读错误（前端通过
  `err.response?.data?.message` 展示），`status` 与 HTTP 状态码一致。
  状态码语义：400 绑定/校验失败，401 未认证（前端清 token 跳登录），
  404 资源不存在，500 服务端错误，502 上游数据源失败（仅 datasource sync）。
- **字段命名**：wire 上统一 snake_case（Go 生成结构体的 `json` tag +
  ts-proto `snakeToCamel=false`），请求与响应双向一致。
- **请求绑定**：handler 用 `c.ShouldBindJSON` 绑定 `gen/go` 的 proto DTO，
  绑定失败返回 400；语义校验在 `internal/api/mapper`。
- 修改信封时两侧文件与本节必须同步更新。

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

Passkey 需要 secure context。本地开发时 `http://localhost:5173` 和 `http://localhost:8080` 均属于 secure context，请把对应 origin 加入 `passkey.rp_origins`。

生产环境则必须通过反向代理提供 HTTPS，并把公共 HTTPS origin 加入 `passkey.rp_origins`。后端本身保持 HTTP 即可。
