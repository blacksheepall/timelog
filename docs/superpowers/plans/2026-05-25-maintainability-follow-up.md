# Maintainability Follow-Up Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Finish the dependency cleanup started in the bootstrap, router, and service layers by removing the remaining DAO/service singletons, then close the repo mismatches uncovered during the refactor so the documented operator workflow matches the code again.

**Architecture:** Replace package-global runtime state with explicit constructors and injected dependencies. The target shape is: `model.NewDao(...)` returns a concrete DAO, service code owns its dependencies explicitly, and router/middleware handlers receive those dependencies from bootstrap code instead of reaching into globals. Keep the existing `router -> service -> model` flow. Avoid domain rewrites, schema changes, or frontend route changes.

**Tech Stack:** Go, Gin, GORM + SQLite, go-webauthn, existing logger wrapper, Go testing, Vite/Vue only for optional API-boundary normalization follow-up

---

## Current State

- `router.RunServer(ctx, r, cfg, logger)` now owns server lifecycle and returns startup/shutdown errors explicitly.
- `router/passkey.go` already receives `*config.Config` instead of reading package-global config.
- `service` and `router/middleware` no longer keep their own `sync.Once` DAO caches; they now use provider indirection for tests.
- `model/model.go` still owns the real singleton through `InitDao` and `GetDao`.
- `service/service.go` still hides DAO and WebAuthn behind package-level provider variables.
- `mcp/init.go` still initializes and reads the global DAO directly.
- `cmd/passkey-temp-admin/main.go`, `Makefile`, and `README.md` still expose `create|list|revoke`, even though the newer design docs now target a create-only DB-backed workflow.
- `Makefile` still relies on a catch-all `%:` rule to swallow extra goals for `passkey-temp`.
- Frontend `TimeLog` field naming was fixed, but API response normalization is still thin and vulnerable to future server/client drift.

## File Map

- Modify: `model/model.go` - add a DAO constructor and phase out global singleton helpers.
- Modify: `service/service.go` - replace provider globals with an explicit service/dependency struct.
- Modify: `service/passkey.go`, `service/timelog.go`, `service/task.go`, `service/constraint.go`, `service/passkey_credentials.go` - convert package-global dependency access to explicit fields.
- Modify: `router/router.go`, `router/passkey.go`, `router/timelog.go`, `router/task.go`, `router/constraint.go` - thread explicit services into route registration.
- Modify: `router/middleware/auth.go`, `router/middleware/deps.go` - accept DAO/token validation dependencies directly instead of package globals.
- Modify: `main.go`, `mcp/init.go`, `cmd/passkey-temp-admin/main.go` - build runtime dependencies explicitly.
- Modify: `router/middleware/auth_test.go`, `service/passkey_test.go`, `test/setup_test.go`, `test/category_tree_test.go` - construct in-memory DAO/service dependencies without singleton state.
- Modify: `Makefile`, `README.md` - reconcile temp-password workflow and remove stale command guidance.
- Optional: `web/src/api/index.ts` - add centralized mapping/normalization for backend response fields.

## Notes Before Editing

- `model.GetDao()` now appears only in `mcp/init.go`, tests, and `model/model.go`, but all service packages still reach it indirectly through `daoProvider`.
- `service.InitService(logger, cfg)` currently only stores logger; it is effectively a bootstrap shim.
- There is no existing application container type; the smallest useful unit is likely a `service.Service` or `service.Dependencies` struct plus direct DAO injection for middleware.
- Route registration functions currently call package-level `service.*` functions. Converting handlers one file at a time is safer than a big-bang rewrite.
- `cmd/passkey-temp-admin/main.go` still has no tests.
- `Makefile` comments and `README.md` currently match the current multi-command CLI, not the create-only redesign docs in `docs/superpowers/specs/2026-04-26-passkey-temp-admin-design.md` and `docs/superpowers/plans/2026-04-26-passkey-temp-admin.md`.
- `config.GetConfig` already honors `TIMELOG_CONFIG_PATH`, so follow-up changes should not reintroduce hardcoded config-path assumptions.

### Task 1: Replace the model DAO singleton with a constructor

**Files:**
- Modify: `model/model.go`
- Modify: `main.go`
- Modify: `mcp/init.go`
- Modify: `cmd/passkey-temp-admin/main.go`
- Test: `go test ./model/... ./mcp/...`

- [ ] **Step 1: Add a constructor that does not touch package globals**

Add `NewDao(cfg *config.Config, logger logger.Logger) (*Dao, error)` that opens SQLite, initializes cache, and returns the DAO without mutating package state.

- [ ] **Step 2: Keep `InitDao` and `GetDao` only as temporary wrappers during the migration**

This keeps the refactor incremental and avoids a single all-or-nothing change.

- [ ] **Step 3: Update production entrypoints to use the constructor result directly**

Switch `main.go`, `mcp/init.go`, and `cmd/passkey-temp-admin/main.go` to hold onto the returned DAO instead of calling `model.GetDao()` after initialization.

- [ ] **Step 4: Remove the singleton once production call sites are gone**

Delete `dao`, `once`, and `GetDao()` only after no non-test code depends on them.

- [ ] **Step 5: Verify**

Run:

```bash
go test ./model/... ./mcp/...
go build ./...
```

### Task 2: Convert service runtime state from package globals to explicit dependencies

**Files:**
- Modify: `service/service.go`
- Modify: `service/passkey.go`
- Modify: `service/timelog.go`
- Modify: `service/task.go`
- Modify: `service/constraint.go`
- Modify: `service/passkey_credentials.go`
- Test: `go test ./service/...`

- [ ] **Step 1: Introduce a small dependency owner for the service layer**

Add a `service.Service` or `service.Dependencies` struct that owns `logger`, `dao`, and optional `webAuthn`.

- [ ] **Step 2: Replace `getDao()` and `getWebAuthn()` with explicit fields**

Convert internal call sites away from `daoProvider` and `webAuthnProvider`.

- [ ] **Step 3: Prefer method receivers for service operations used by HTTP, MCP, and CLI entrypoints**

Keep thin package-level wrappers only if they materially reduce churn during the migration.

- [ ] **Step 4: Remove `InitService` and the testing-only provider shims when unused**

Once all call sites construct explicit service instances, delete the bootstrap/test indirection.

- [ ] **Step 5: Verify**

Run:

```bash
go test ./service/...
```

### Task 3: Thread explicit services into router and middleware

**Files:**
- Modify: `router/router.go`
- Modify: `router/passkey.go`
- Modify: `router/timelog.go`
- Modify: `router/task.go`
- Modify: `router/constraint.go`
- Modify: `router/middleware/auth.go`
- Modify: `router/middleware/deps.go`
- Modify: `main.go`
- Test: `go test ./router/...`

- [ ] **Step 1: Change route registration helpers to accept concrete dependencies**

Route setup should receive the service object instead of relying on package-level `service.*` functions.

- [ ] **Step 2: Replace middleware DAO providers with constructor-style middleware**

For example, `middleware.Auth(dao *model.Dao)` or an equally small token-store dependency.

- [ ] **Step 3: Keep `router.Register` responsible for wiring only**

Do not move business logic out of `service` during this refactor.

- [ ] **Step 4: Update `main.go` to construct and pass the shared dependencies once**

Bootstrap should create the DAO/service/WebAuthn instance and inject them into router registration.

- [ ] **Step 5: Verify**

Run:

```bash
go test ./router/... ./service/...
```

### Task 4: Update tests to stop depending on process-wide singleton state

**Files:**
- Modify: `router/middleware/auth_test.go`
- Modify: `service/passkey_test.go`
- Modify: `test/setup_test.go`
- Modify: `test/category_tree_test.go`
- Test: `go test ./test/... ./router/... ./service/...`

- [ ] **Step 1: Replace test helpers that call `model.InitDao`**

Construct a fresh in-memory DAO directly instead of relying on shared package state.

- [ ] **Step 2: Isolate any unavoidable shared state with explicit reset helpers during the transition**

Do not let tests depend on initialization order across packages.

- [ ] **Step 3: Add focused constructor/wiring tests**

Cover the new dependency construction path so singleton regressions are caught early.

- [ ] **Step 4: Verify**

Run:

```bash
go test ./test/... ./router/... ./service/...
```

### Task 5: Reconcile the temp-password operator workflow and Makefile behavior

**Files:**
- Modify: `cmd/passkey-temp-admin/main.go`
- Create: `cmd/passkey-temp-admin/main_test.go`
- Modify: `Makefile`
- Modify: `README.md`
- Reference: `docs/superpowers/specs/2026-04-26-passkey-temp-admin-design.md`
- Reference: `docs/superpowers/plans/2026-04-26-passkey-temp-admin.md`
- Test: `go test ./cmd/passkey-temp-admin`

- [ ] **Step 1: Decide whether the repo still wants the create-only CLI redesign**

Current design docs say yes; current implementation still says no. Do not keep both stories.

- [ ] **Step 2: If create-only remains the desired behavior, finish the earlier plan**

That means removing `list` and `revoke`, adding command tests, and updating `README.md` and `Makefile` examples.

- [ ] **Step 3: Replace the global `%:` catch-all in `Makefile`**

Use a narrower pattern or variable-driven argument passing so unknown targets are not silently swallowed across the whole repo.

- [ ] **Step 4: Verify**

Run:

```bash
go test ./cmd/passkey-temp-admin
make passkey-temp
```

Also run any revised documented invocation if the Makefile interface changes.

### Task 6: Optional guard against future frontend/backend schema drift

**Files:**
- Modify: `web/src/api/index.ts`
- Possibly modify: `web/src/types/index.ts`
- Test: `cd web && pnpm run type-check`

- [ ] **Step 1: Add lightweight response normalization at the API boundary**

If the backend is likely to keep mixed field names, keep the mapping centralized in the API layer rather than spreading fallback field logic through Vue components.

- [ ] **Step 2: Keep this optional and scoped**

Do not re-open the frontend unless backend field drift remains a realistic risk after the service/refactor work.

- [ ] **Step 3: Verify**

Run:

```bash
cd web && pnpm run type-check
cd web && pnpm run build
```

## Verification Sweep

```bash
go test ./...
go build ./...
cd web && pnpm run build
make buildx
```

## Expected Outcome

After this follow-up work:

- bootstrap code owns all runtime dependencies explicitly
- `model` no longer requires a process-global DAO
- router, middleware, CLI, MCP, and tests use the same constructor-based wiring model
- temp-password docs and Makefile commands describe exactly one supported workflow
- frontend API integration is less likely to break on simple JSON field drift
