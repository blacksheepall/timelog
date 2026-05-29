# MCP Architecture & Configuration

This document describes how the TimeLog MCP server is structured, how it is configured, and how to extend it safely.

## Architecture

The MCP server is a standalone binary that exposes the TimeLog data as MCP tools.

- **Transport**: stdio or StreamableHTTP
- **Server**: `mcp.NewServer` with tool registrations in `mcp/server.go`
- **Handlers**: MCP tool logic lives in `mcp/handlers.go`
- **Logging**: MCP logs are file-only and configured via top-level `mcp` settings

## Configuration

MCP configuration is a top-level section in `config.yml`:

```yaml
mcp:
  enabled: false
  level: debug
  path: ./logs/mcp.log
  transport: stdio
  listen_addr: ":8080"
  token: "" # Authorization: Bearer <token>
```

### Environment Overrides

Each field can also be provided via environment variables:

- `MCP_ENABLED`
- `MCP_LEVEL`
- `MCP_PATH`
- `MCP_TRANSPORT`
- `MCP_LISTEN_ADDR`
- `MCP_TOKEN`

When both are set, environment variables take precedence.

## Transport Modes

### stdio (default)

Use for local MCP clients that spawn the MCP server directly.

```bash
MCP_TRANSPORT=stdio ./timelog-mcp-server
```

### http (StreamableHTTP)

Use when the MCP server is accessed remotely over HTTP.

```bash
MCP_TRANSPORT=http MCP_LISTEN_ADDR=":8080" ./timelog-mcp-server
```

When `token` is configured, requests must include:

```
Authorization: Bearer <token>
```

The `/health` endpoint does not require authentication.

## Tool List

Keep this list synchronized with `mcp/server.go`:

- `get_timelogs_by_date_range` — use `"active_only": true` to fetch ongoing time logs
- `get_tasks_by_status`
- `get_active_constraints`
- `get_date_info`
- `list_categories`
- `create_timelog`
- `update_timelog`

MCP tools call the same `service` layer as the HTTP API. Transport-specific JSON shaping stays in MCP handlers; business rules live in `service/`.

## Maintenance Rules

- When adding or modifying tools, update this document and `mcp/README.md`.
- Avoid stdout logging in MCP mode; use the MCP logger only.
