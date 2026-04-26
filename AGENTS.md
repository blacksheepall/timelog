# AGENTS.md

This file is for agentic coding tools working in `/Users/n/Documents/Developer/timelog`.

## Project Summary
- Full-stack local time logging app.
- Backend: Go + Gin + GORM + SQLite.
- Frontend: Vue 3 + TypeScript + Vite + Tailwind + Pinia + Vue Router.
- Production serves the built frontend from the Go binary via `embed`.
- Auth includes passkey flows; most API routes are protected.

## Instruction Files Checked
- Root guidance exists in `CLAUDE.md`.
- Frontend guidance exists in `web/CLAUDE.md`.
- No repo-local `.cursorrules`, `.cursor/rules/`, or `.github/copilot-instructions.md` file was found.

## Repository Layout
- `main.go`: app entrypoint, config load, logger init, embedded frontend.
- `router/`: Gin route registration and HTTP handlers.
- `router/middleware/`: auth and CORS middleware.
- `service/`: thin business layer over model operations.
- `model/`: DAO setup, data access, migrations, generated-model integration.
- `model/gen/`: generated GORM models; treat as generated code.
- `core/config/`, `core/logger/`: config loading and zap-based logging.
- `test/`: integration tests. `web/`: Vue SPA. `mcp/`: MCP server code.

## Build And Run Commands
Run from repo root unless noted.

### Backend
- `make build` - build Go binary as `./main`.
- `make build-lite` - smaller local binary.
- `make build-linux` - Linux production binary with `prod` tag.
- `make buildx` - build frontend, then backend.
- `./main` - run backend after building.

### Frontend
- `cd web && pnpm install` - install frontend deps.
- `cd web && pnpm run dev` - start Vite dev server on `:3000`.
- `cd web && pnpm run build` - type-check and production build.
- `cd web && pnpm run preview` - preview built frontend.
- `cd web && pnpm run type-check` - TypeScript checking only.

### Other
- `make web` - install frontend deps and build frontend.
- `make mcp` - build MCP server binary.
- `make docker env=prod` - build Docker image.
- `make run env=prod` - run Dockerized app.
- `make clean` - remove binaries plus `web/dist` and `web/node_modules`.

## Test Commands
### Main entrypoints
- `go test ./...` - run all Go tests.
- `go test ./test/...` - run integration tests only.
- `go test ./router/...` - run router and middleware tests.
- `go test ./service/...` - run service tests.
- `go test ./model/...` - run model tests.

### Single-test patterns
- `go test ./service -run TestPasskeySessionKeyNamespacing -v`
- `go test ./router/middleware -run TestAuthMiddlewareRejectsPasskeySession -v`
- `go test ./model -run TestBuildCategoryTreePointers -v`
- `go test ./test -run TestCategoryTreeStructure -v`
- `go test ./test -run '^TestCategoryTreeMultipleRoots$' -v` for exact-name matching.

### Frontend verification
- No dedicated frontend unit test script exists in `web/package.json`.
- Use `cd web && pnpm run type-check` for fast validation.
- Use `cd web && pnpm run build` for the strongest frontend verification.

## Formatting And Linting
- `make fmt` is the repo formatting command.
- `make fmt` runs `go fmt ./...`, `cd mcp && go fmt ./...`, and Prettier on `web/src/`.
- There is no dedicated Go lint target and no ESLint script currently configured.
- If you touch frontend code, at minimum run `cd web && pnpm run type-check`; prefer `build` for route/layout/API-impacting changes.

## Swagger / Generated Artifacts
- Read `README.md` before changing Swagger-related backend code.
- Dev/test workflows may require `swag init` before `go mod tidy` or some builds.
- `model/gen/` is generated output; avoid manual edits unless the repo clearly expects it.
- Use `make gen-model` for model regeneration.

## Configuration Notes
- Default runtime config is `config.yml`; example config is `config-example.yml`.
- Integration tests use `config-test.yml` according to `test/setup_test.go` and `CLAUDE.md`.
- Test DB flushing is controlled by `test.flush` in config.
- Frontend dev server expects backend access on the configured local port.

## Backend Conventions (Go)
- Use `gofmt`; do not hand-format Go code.
- Keep imports in standard Go style: stdlib, blank line, third-party/internal.
- Follow existing architecture: router -> service -> model.
- Keep HTTP binding/parsing in `router/`, business rules in `service/`, persistence/tree logic in `model/`.
- Prefer existing helpers like `SuccessResponse`, `ErrorResponse`, and `parseInt32Param`.
- Most IDs are `int32`; match existing signatures instead of widening to `int`.
- Generated model fields often use pointers; preserve pointer semantics carefully.
- Use `model.GetDao().Db()` in service functions.
- When adding error context, use `fmt.Errorf("...: %w", err)` and return errors upward.
- Keep handler error paths explicit with early returns after `c.JSON(...)`.

## Go Naming And Structure
- Exported names: PascalCase. Unexported names: camelCase. Package names: short, lowercase, no underscores.
- Keep files focused by domain (`timelog.go`, `task.go`, `constraint.go`, `passkey.go`).
- Prefer a small nearby helper over a catch-all util file.
- Preserve existing comment/docstring style; many router comments are Chinese Swagger comments.

## Error Handling Expectations
- Validate request input close to the handler boundary.
- Convert parse/bind failures into `400`, missing resources into `404`, and non-client service/model failures into `500`.
- In the frontend, keep `try/catch` blocks around async API calls.
- On frontend failures, follow existing patterns: set local error state, `console.error`, and notify the user.

## Frontend Conventions (Vue + TypeScript)
- Use `<script setup lang="ts">` for Vue SFCs.
- Use 2-space indentation, single quotes, and no semicolons, matching current code.
- Prefer `@/` aliases and `import type` for type-only imports.
- Keep external imports before internal alias imports.
- Component/view filenames use PascalCase; composables use `useXxx`; stores use Pinia `defineStore` with short domain names.
- Use `ref`, `reactive`, `computed`, `watch`, and `onMounted` in the existing Composition API style.
- Keep API access centralized in `web/src/api/index.ts` and preserve shared `ApiResponse<T>` wrapper types.
- Keep router metadata titles updated when adding pages.

## Frontend UI / Data Patterns
- Prefer Tailwind utility classes in templates over ad hoc CSS.
- Follow existing loading/error/empty-state rendering patterns in views.
- Keep user feedback explicit with the injected `showNotification` pattern.
- For list loads, current code often uses `response.data || []`; keep behavior consistent unless refactoring broadly.
- Route date formatting through `web/src/utils/date`; keep auth token handling in `web/src/utils/auth` and axios interceptors.

## Testing Guidance For Agents
- Run the narrowest relevant test first, then broaden scope.
- For backend logic changes, start with a single-package `go test` command using `-run`.
- For HTTP handler changes, prefer router or integration tests depending on scope.
- Before finishing a task, run at least one verification command that directly covers your edits.

## Practical Workflow
- Read `CLAUDE.md` and `web/CLAUDE.md` before larger changes.
- Check whether your change affects generated artifacts, config, or embedded frontend output.
- Prefer minimal, pattern-matching changes over broad refactors.
- Do not edit unrelated dirty files.
- If you change frontend code intended for production serving, remember the backend serves `web/dist` output after a build.
