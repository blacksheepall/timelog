# AGENTS.md

## Snapshot
- Full-stack local time logging app: Go backend (`main.go`) + Vue 3/Vite frontend (`web/`) + SQLite.
- `main.go` embeds `web/dist`; production and local binary runs serve the built SPA from the Go binary.
- Passkey auth is optional, but when `passkey.enabled` is true, all main `/api` routes are behind `router/middleware/auth.go` bearer-token auth.

## Read First
- `README.md` for Swagger, passkey, migration, and temp-password workflows.
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
- Vite dev server runs on `http://localhost:5173` and proxies `/api` to `http://localhost:8083`.
- Passkeys require HTTPS. Cert generation lives in `scripts/generate-certs.sh`; see `README.md` for the exact setup.

## Verified Code Boundaries
- HTTP wiring lives in `router/`; `router/router.go` registers route groups, optional auth, Swagger, and SPA fallback.
- Business logic is thin `service/` wrappers over `model/`; preserve the existing `router -> service -> model` flow.
- `model/gen/` is generated and gitignored. Do not hand-edit it.
- Frontend API access is centralized in `web/src/api/index.ts`.

## Codegen / Build Quirks
- Non-`prod` builds compile `router/swagger.go`, which imports `github.com/blacksheepaul/timelog/docs`. If Swagger files are missing, run `swag init` before `go mod tidy` or non-prod builds.
- `make gen-model` requires `TIMELOG_GEN_DB_PATH` to be set; `model/gentool/gormgen.go` panics otherwise.
- CI release build runs `make install-deps gen-model buildx-linux`; do not assume `make build-linux` alone matches release artifacts.

## Testing Notes
- `test/setup_test.go` reloads the DB from migrations when `test.flush` is true in `config-test.yml`.
- For backend changes, prefer the narrowest package test first, then broaden.

## Frontend Notes
- Prettier style is enforced by repo config: 2 spaces, single quotes, no semicolons.
- Route titles are set from `meta.title` in `web/src/router/index.ts`; keep them updated when adding routes.
- The frontend still exports deprecated `tagAPI` as an alias of `categoryAPI`; avoid introducing new tag-based code.
