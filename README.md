# timelog

Yet another lyubishchev time management implementation.

Core functions:

- Time tracking
- LLM review based on time tracking data (via MCP)
- Task management (wip)

# How to use

# How to build and run

For development setup, generated-file rules, and required checks before committing, read:

- 中文（主文档）：[development-cn.md](docs/development-cn.md)
- English (AI Translated): [development.md](docs/development.md)

## API Documentation (Development/Testing Only)

Non-`prod` builds serve Redoc at `/docs/redoc.html` (legacy redirect: `/swagger`).

Regenerate the OpenAPI spec after proto or route changes:

```bash
make gen-api
```

DTO schemas come from `api/proto/timelog/v1`; REST paths live in `api/openapi/rest.yaml`. The merged spec is written to `router/docs/openapi.yaml`.

**Note:** Production builds exclude the docs routes via the `prod` build tag.

## API Contract Generation

API payload DTOs live under `api/proto/timelog/v1`.

Run:

```bash
make gen-api
make check-api
```

`make gen-api` regenerates Go DTOs under `gen/go/`, TypeScript DTOs under `web/src/gen/`, and the merged OpenAPI spec at `router/docs/openapi.yaml`.
The REST envelope `{ data, message, status }` is intentionally hand-written in Go and TypeScript and is not part of the proto contract.
The Buf CLI is installed through the frontend dev dependencies, so run `cd web && pnpm install` before using these targets in a fresh checkout.

## Passkey Setup

Passkey authentication requires HTTPS. Follow the certificate setup below, then generate a one-time temp password to bind a device in the UI.

### Certificate Setup (Required for Passkey)

Since Passkey authentication requires a secure context (HTTPS), generate a self-signed certificate:

```bash
make gen-certs
```

This creates `certs/cert.pem` and `certs/key.pem` in the project root.
The default passkey config uses `timelog.local`; for local testing, make sure it resolves to `127.0.0.1` (for example, add `127.0.0.1 timelog.local` to `/etc/hosts`).

Update your `config.yml`:

```yaml
server:
  https_enabled: true
  cert_file: ./certs/cert.pem
  key_file: ./certs/key.pem
```

**Note:** For local development, you may need to accept the self-signed certificate warning in your browser. With the default passkey config, access the application at `https://timelog.local:8080` (or update `passkey.rp_origins` if you use a different origin).

### Generate temp password for the current app instance

If you need a temp password that the current app instance can immediately accept, use the DB-backed command:

```bash
make passkey-temp
```

Optional: specify TTL seconds (defaults to config `passkey.temp_password.ttl`).

```bash
make passkey-temp PASSKEY_TEMP_TTL=900
```

Or run the CLI directly:

```bash
go run ./cmd/passkey-temp-admin
go run ./cmd/passkey-temp-admin 900
```

If your config file is not at the project root, set `TIMELOG_CONFIG_PATH` before running the command.

### Bind device

Open `https://timelog.local:8080/passkey/register` and use the temp password to create a passkey.

### Login

Open `https://timelog.local:8080/login` and complete the passkey prompt.

## Migrate

for example:

create new migration

```bash
migrate -database "sqlite3://dev.db" create -seq -ext sql --dir model/migrations/ init_xxx_table
```

forward

```bash
migrate -database "sqlite3://dev.db" --path model/migrations/ up
# or use make target (Makefile defaults to prod; use env=dev for local development)
make migrate env=dev
make migrate env=prod
```

# Launch

```bash
# Build binary
make buildx
```

# How to Deployment

- English: [deploy.md](docs/deploy.md)
- 中文: [deploy-cn.md](docs/deploy-cn.md)

### TODO

- **AI-driven summary**: Day/week productivity summaries with frontend stream output
- **Authentication upgrade**: From basic-auth to more secure authentication
- **Passkey support**: Add support for passkey authentication
- **Tag refactoring**: Convert tags into categories for better organization
- **Automated Reporting**: Daily reports generated at 4 AM showing task completion vs. time estimates
- **Advanced Analytics**: Better visualization of productivity patterns and time allocation
- **Task Templates**: Reusable task templates for common work patterns
