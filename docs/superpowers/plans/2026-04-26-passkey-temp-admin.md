# Passkey Temp Admin Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Convert `cmd/passkey-temp-admin` into a create-only CLI that reads `config.yml`, writes a temp password record into SQLite, and prints the plaintext password plus expiry for the existing passkey registration flow.

**Architecture:** Keep temp-password persistence in the existing `service.CreateTempPassword` and `model.TempPassword` path. Refactor the CLI entrypoint into small helpers so tests can cover TTL parsing and successful record creation without depending on the full HTTP app. Update the README so operator-facing usage matches the new create-only command shape.

**Tech Stack:** Go, cleanenv config loading, existing zap logger wrapper, GORM + SQLite, Go `testing`

---

## File Map

- Modify: `cmd/passkey-temp-admin/main.go` — remove subcommands, add create-only argument parsing, keep config/logger/DAO wiring, and print the created password.
- Create: `cmd/passkey-temp-admin/main_test.go` — add focused tests for TTL parsing, invalid arguments, and successful temp-password persistence.
- Modify: `README.md` — replace old `create|list|revoke` examples with the new create-only usage.

## Notes Before Editing

- `cmd/passkey-temp-admin/main.go` currently requires a subcommand and mixes create/list/revoke behavior.
- `service.CreateTempPassword(ttl)` already generates the plaintext password, hashes it, and inserts into `temp_passwords`.
- `model.InitDao` is a package singleton, so tests should initialize it once and clean temp-password rows between cases instead of trying to reinitialize it repeatedly.
- No router, frontend, or passkey WebAuthn code should change for this work.

### Task 1: Lock in create-only CLI argument behavior

**Files:**
- Create: `cmd/passkey-temp-admin/main_test.go`
- Modify: `cmd/passkey-temp-admin/main.go`
- Test: `cmd/passkey-temp-admin/main_test.go`

- [ ] **Step 1: Write the failing TTL parsing test**

```go
func TestResolveTTL(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		defaultTTL int
		wantTTL    int
		wantErr    string
	}{
		{
			name:       "uses config default when ttl omitted",
			args:       nil,
			defaultTTL: 900,
			wantTTL:    900,
		},
		{
			name:       "uses positional ttl override",
			args:       []string{"1200"},
			defaultTTL: 900,
			wantTTL:    1200,
		},
		{
			name:       "rejects negative ttl",
			args:       []string{"-1"},
			defaultTTL: 900,
			wantErr:    "ttl must be >= 0",
		},
		{
			name:       "rejects extra arguments",
			args:       []string{"900", "extra"},
			defaultTTL: 900,
			wantErr:    "too many positional arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTTL, err := resolveTTL(tt.args, tt.defaultTTL)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveTTL returned error: %v", err)
			}
			if gotTTL != tt.wantTTL {
				t.Fatalf("expected ttl %d, got %d", tt.wantTTL, gotTTL)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/passkey-temp-admin -run TestResolveTTL -v`

Expected: FAIL with `undefined: resolveTTL`

- [ ] **Step 3: Implement the minimal parsing helpers**

```go
func resolveTTL(args []string, defaultTTL int) (int, error) {
	if len(args) > 1 {
		return 0, errors.New("too many positional arguments")
	}
	if len(args) == 0 {
		return defaultTTL, nil
	}

	ttl, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("invalid ttl: %w", err)
	}
	if ttl < 0 {
		return 0, errors.New("ttl must be >= 0")
	}

	return ttl, nil
}

func usage() string {
	return "Usage: go run ./cmd/passkey-temp-admin [ttl]"
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./cmd/passkey-temp-admin -run TestResolveTTL -v`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/passkey-temp-admin/main.go cmd/passkey-temp-admin/main_test.go
git commit -m "test: cover passkey temp ttl parsing" -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

### Task 2: Make the CLI create and persist temp passwords without subcommands

**Files:**
- Modify: `cmd/passkey-temp-admin/main.go`
- Modify: `cmd/passkey-temp-admin/main_test.go`
- Test: `cmd/passkey-temp-admin/main_test.go`

- [ ] **Step 1: Add a failing success-path test that proves the record is persisted**

```go
type fakeLogger struct{}

func (fakeLogger) Debug(fields ...interface{})                     {}
func (fakeLogger) Debugw(msg string, keysAndValues ...interface{}) {}
func (fakeLogger) Info(fields ...interface{})                      {}
func (fakeLogger) Infow(msg string, keysAndValues ...interface{})  {}
func (fakeLogger) Warn(fields ...interface{})                      {}
func (fakeLogger) Warnw(msg string, keysAndValues ...interface{})  {}
func (fakeLogger) Error(fields ...interface{})                     {}
func (fakeLogger) Errorw(msg string, keysAndValues ...interface{}) {}
func (fakeLogger) Fatal(fields ...interface{})                     {}
func (fakeLogger) Fatalw(msg string, keysAndValues ...interface{}) {}

var testCfg = &config.Config{}

func setupCommandTestEnv(t *testing.T) {
	t.Helper()

	testCfg.Database.Host = ":memory:"
	testCfg.Log.ORMLogLevel = 1
	testCfg.Passkey.TempPassword.TTL = 900

	service.InitService(fakeLogger{}, testCfg)
	model.InitDao(testCfg, fakeLogger{})

	if err := model.GetDao().Db().AutoMigrate(&model.TempPassword{}); err != nil {
		t.Fatalf("auto migrate temp_passwords: %v", err)
	}
	if err := model.GetDao().Db().Exec("DELETE FROM temp_passwords").Error; err != nil {
		t.Fatalf("clear temp_passwords: %v", err)
	}
}

func TestRunCreateUsesDefaultTTLAndPersistsRecord(t *testing.T) {
	setupCommandTestEnv(t)

	var stdout bytes.Buffer
	err := runCreate(nil, testCfg, &stdout)
	if err != nil {
		t.Fatalf("runCreate returned error: %v", err)
	}

	var count int64
	if err := model.GetDao().Db().Model(&model.TempPassword{}).Count(&count).Error; err != nil {
		t.Fatalf("count temp_passwords: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 temp password record, got %d", count)
	}
	if !strings.Contains(stdout.String(), "temp password:") {
		t.Fatalf("expected plaintext password in output, got %q", stdout.String())
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/passkey-temp-admin -run TestRunCreateUsesDefaultTTLAndPersistsRecord -v`

Expected: FAIL with `undefined: runCreate`

- [ ] **Step 3: Refactor the CLI into a create-only execution path**

```go
func runCreate(args []string, cfg *config.Config, stdout io.Writer) error {
	ttl, err := resolveTTL(args, cfg.Passkey.TempPassword.TTL)
	if err != nil {
		return err
	}

	record, password, err := service.CreateTempPassword(ttl)
	if err != nil {
		return fmt.Errorf("failed to create temp password: %w", err)
	}

	fmt.Fprintf(stdout, "temp password: %s\n", password)
	fmt.Fprintf(stdout, "expires at: %s\n", record.ExpiresAt.Format("2006-01-02 15:04:05"))
	return nil
}

func main() {
	cfg := config.GetConfig("config.yml")
	logger := log.SetZapLogger(*cfg)
	service.InitService(logger, cfg)
	model.InitDao(cfg, logger)

	if err := runCreate(os.Args[1:], cfg, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		fmt.Fprintln(os.Stderr, usage())
		os.Exit(1)
	}
}
```

Also delete the `switch` over `create|list|revoke` and remove the old `printUsage()` function.

- [ ] **Step 4: Add a failing invalid-argument test, then make it pass**

```go
func TestRunCreateRejectsExtraArgs(t *testing.T) {
	setupCommandTestEnv(t)

	var stdout bytes.Buffer
	err := runCreate([]string{"900", "extra"}, testCfg, &stdout)
	if err == nil || err.Error() != "too many positional arguments" {
		t.Fatalf("expected extra-argument error, got %v", err)
	}
}
```

Run: `go test ./cmd/passkey-temp-admin -run 'TestRunCreateUsesDefaultTTLAndPersistsRecord|TestRunCreateRejectsExtraArgs' -v`

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/passkey-temp-admin/main.go cmd/passkey-temp-admin/main_test.go
git commit -m "feat: simplify passkey temp admin cli" -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

### Task 3: Update operator-facing docs and run focused verification

**Files:**
- Modify: `README.md`
- Test: `cmd/passkey-temp-admin/main_test.go`

- [ ] **Step 1: Update the README examples to match the new command shape**

Replace the current section with examples like:

~~~md
### Generate temp password for the current app instance

If you need a temp password that the current app instance can immediately accept, use the DB-backed command:

```bash
make passkey-temp
```

Optional: specify TTL seconds (defaults to config `passkey.temp_password.ttl`).

```bash
go run ./cmd/passkey-temp-admin 900
```
~~~

Remove the `list` and `revoke` examples from the README because those behaviors are no longer supported.

- [ ] **Step 2: Run the focused tests**

Run: `go test ./cmd/passkey-temp-admin -v`

Expected: PASS

- [ ] **Step 3: Run the repo-wide verification most relevant to the touched code**

Run: `go test ./service/... ./router/... ./cmd/passkey-temp-admin`

Expected: PASS

- [ ] **Step 4: Sanity-check the command manually**

Run: `go run ./cmd/passkey-temp-admin`

Expected: stdout prints:

```text
temp password: <random value>
expires at: <timestamp>
```

and SQLite `temp_passwords` gains one new row.

- [ ] **Step 5: Commit**

```bash
git add README.md cmd/passkey-temp-admin/main.go cmd/passkey-temp-admin/main_test.go
git commit -m "docs: update passkey temp password usage" -m "Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

## Self-Review Checklist

- Spec coverage: covered by Task 1 (CLI shape), Task 2 (create-only persistence path), and Task 3 (operator docs + verification).
- Placeholder scan: no `TODO`/`TBD` placeholders remain.
- Type consistency: `resolveTTL` and `runCreate` are used consistently across all tasks.
