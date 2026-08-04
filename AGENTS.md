# AGENTS.md

## Snapshot
- Full-stack local time logging app: Go backend (`main.go`) + Vue 3/Vite frontend (`web/`) + SQLite.
- `main.go` embeds `web/dist`; production and local binary runs serve the built SPA from the Go binary.
- Passkey auth is optional, but when `passkey.enabled` is true, all main `/api` routes are behind `router/middleware/auth.go` bearer-token auth.

## Read First
- `docs/development-cn.md` is the canonical human developer guide; `docs/development.md` is the AI-translated English version.
- `README.md` for high-level usage, API docs, passkey, migration, and temp-password workflows.
- `web/CLAUDE.md` before frontend changes.

## Commands That Matter
- Full app build: `make buildx`.
- Backend-only build: `make build`.
  Use this only if `web/dist` already exists; `main.go` embeds that directory and plain `make build` does not build the frontend first.
- Frontend dev: `cd web && pnpm install && pnpm run dev`.
- Frontend verification: `cd web && pnpm run type-check`; use `cd web && pnpm run build` for stronger verification.
- Backend tests: `go test ./...`.
- Focused tests: `go test ./router/...`, `go test ./service/...`, `go test ./model/...`, `go test ./test/...`.
- Formatting: `make fmt`.

## Config And Runtime Gotchas
- App startup hardcodes `config.yml` (`main.go`, `cmd/passkey-temp-admin/main.go`). Tests hardcode `config-test.yml` (`test/setup_test.go`). If those files are missing, startup panics.
- Both `config.yml` and `config-test.yml` are gitignored. Copy from `config-example.yml` when needed.
- An optional `.env` in the working directory is loaded before config resolution (`core/config/config.go`). cleanenv injects every entry into the process environment via `os.Setenv` — overwriting pre-existing vars — so `.env` is the place for variables like `TIMELOG_CONFIG_PATH` and the `MCP_*` settings (the only struct fields with `env` tags). Copy from `.env.example`.
- Vite dev server runs on `http://localhost:5173` and proxies `/api` to `http://localhost:8080`.
- Passkeys require a secure context. In production this means serving the app behind a reverse proxy that terminates HTTPS; the backend itself runs plain HTTP. For local development, `localhost` origins are also secure contexts. See `README.md` for passkey setup.

## Verified Code Boundaries
- HTTP wiring lives in `router/`; `router/router.go` registers route groups, optional auth, API docs, and SPA fallback.
- Business logic is thin `service/` wrappers over `model/`; preserve the existing `router -> service -> model` flow.
- `model/gen/` is generated and gitignored. Do not hand-edit it.
- Frontend API access is centralized in `web/src/api/index.ts`.

## Codegen / Build Quirks
- Keep generated outputs out of commits unless explicitly requested; see `docs/development-cn.md` for source paths and generated paths.
- Non-`prod` builds compile `router/docs.go`, which embeds generated `router/docs/openapi.yaml`. `make build`, `make test`, and `make web` run `gen-api` automatically.
- `make gen-model` requires `TIMELOG_GEN_DB_PATH` to be set; `model/gentool/gormgen.go` panics otherwise.
- API DTOs are generated from `api/proto/timelog/v1` with `make gen-api`; do not hand-edit `gen/go/` or `web/src/gen/`. REST paths live in `api/openapi/rest.yaml`; `cmd/merge-openapi` merges proto JSON schemas into `router/docs/openapi.yaml`. The REST envelope remains hand-written in `router/apiresponse.go` and `web/src/types/api.ts`. The Buf CLI is provided by `web` dev dependencies, so run `cd web && pnpm install` first in a fresh checkout. `protoc-gen-jsonschema` is installed via `make install-deps`.
- CI release build runs `make install-deps gen-model buildx-linux`; `buildx-linux` runs `gen-api` via `web`/`build-linux`. Do not assume plain `go build` without codegen matches release artifacts.

## Testing Notes
- `test/setup_test.go` reloads the DB from migrations when `test.flush` is true in `config-test.yml`.
- For backend changes, prefer the narrowest package test first, then broaden.

## Frontend Notes
- Prettier style is enforced by repo config: 2 spaces, single quotes, no semicolons.
- Route titles are set from `meta.title` in `web/src/router/index.ts`; keep them updated when adding routes.
- The frontend still exports deprecated `tagAPI` as an alias of `categoryAPI`; avoid introducing new tag-based code.
