# MCP 服务器配置指南

本文档介绍 TimeLog MCP 服务器的架构、配置方法以及如何安全地扩展功能。

## 架构概述

MCP 服务器是一个独立的可执行文件，它将 TimeLog 数据以 MCP 工具的形式暴露给 AI 助手。

- **传输层**：支持 stdio 或 StreamableHTTP
- **服务器核心**：使用 `mcp.NewServer` 创建，工具注册在 `mcp/server.go` 中
- **处理器**：MCP 工具的业务逻辑位于 `mcp/handlers.go`
- **日志**：MCP 日志仅写入文件，通过顶级的 `mcp` 配置项进行设置

## 配置说明

MCP 配置是 `config.yml` 中的一个顶级段落：

```yaml
mcp:
  enabled: false
  level: debug
  path: ./logs/mcp.log
  transport: stdio
  listen_addr: ":8080"
  token: "" # Authorization: Bearer <token>
```

### 环境变量覆盖

每个配置字段都可以通过环境变量设置：

- `MCP_ENABLED`
- `MCP_LEVEL`
- `MCP_PATH`
- `MCP_TRANSPORT`
- `MCP_LISTEN_ADDR`
- `MCP_TOKEN`

当配置文件和环境变量同时设置时，环境变量优先。

## 传输模式

### stdio（默认）

适用于直接启动 MCP 服务器的本地 AI 客户端。

```bash
MCP_TRANSPORT=stdio ./timelog-mcp-server
```

### http（StreamableHTTP）

适用于通过 HTTP 远程访问 MCP 服务器的场景。

```bash
MCP_TRANSPORT=http MCP_LISTEN_ADDR=":8080" ./timelog-mcp-server
```

当配置了 `token` 时，请求必须包含：

```
Authorization: Bearer <token>
```

注意：`/health` 端点不需要认证。

## 工具列表

使用以下工具与 TimeLog 数据交互（保持与 `mcp/server.go` 同步）：

- `get_timelogs_by_date_range` - 获取指定日期范围内的时间记录；设置 `"active_only": true` 可获取正在进行的计时
- `get_tasks_by_status` - 按完成状态筛选任务
- `get_active_constraints` - 了解当前自我约束和外部条件
- `get_date_info` - 获取当前日期、时间、今天、昨天以及本周日期范围
- `list_categories` - 列出所有可用分类
- `create_timelog` - 创建时间记录
- `update_timelog` - 更新时间记录

MCP 工具与 HTTP API 共用同一套 `service` 业务层；MCP handler 只负责传输层 JSON 格式化，业务规则集中在 `service/`。

## 维护规则

- 添加或修改工具时，更新本文档以及 `mcp/README.md`
- 在 MCP 模式下避免使用 stdout 输出日志，仅使用 MCP 日志记录器
