# timelog

Yet another lyubishchev time management implementation.

Core functions:

- Time tracking
- LLM review based on time tracking data (via MCP)
- Task management (wip)

# How to use

# How to build and run

## Swagger Setup (Development/Testing Only)

For development and testing environments, you need to generate Swagger documentation using the `swag` tool before running `go mod tidy` or building the project:

### Install swag

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Swagger documentation

```bash
swag init
```

This will generate the necessary files in the `docs` directory that are required by the router package.

If `swag init` fails because `web/dist` is missing, build the frontend once first:

```bash
make web
make swagger
```

**Note:** For production builds, Swagger is automatically excluded via the `prod` build tag, so you don't need to generate Swagger documentation for production deployments.

## Passkey Setup

Passkey authentication requires HTTPS. Follow the certificate setup below, then generate a one-time temp password to bind a device in the UI.

### Certificate Setup (Required for Passkey)

Since Passkey authentication requires a secure context (HTTPS), generate a self-signed certificate:

```bash
make gen-certs
```

This creates `certs/cert.pem` and `certs/key.pem` in the project root.

Update your `config.yml`:

```yaml
server:
  https_enabled: true
  cert_file: ./certs/cert.pem
  key_file: ./certs/key.pem
```

**Note:** For local development, you may need to accept the self-signed certificate warning in your browser. The application will be accessible at `https://localhost:8083` (or your configured port).

### Generate temp password for the current app instance

If you need a temp password that the current app instance can immediately accept, use the DB-backed admin command:

```bash
make passkey-temp create
```

Optional: specify TTL seconds (defaults to config `passkey.temp_password.ttl`).

```bash
make passkey-temp create 900
```

If your config file is not at the project root, set `TIMELOG_CONFIG_PATH` before running the command.

### List or revoke DB-backed temp passwords

```bash
make passkey-temp list
make passkey-temp revoke <id>
```

### Bind device

Open `https://localhost:5173/passkey/register` and use the temp password to create a passkey.

### Login

Open `https://localhost:5173/login` and complete the passkey prompt.

## Migrate

for example:

create new migration

```bash
migrate -database "sqlite3://dev.db" create -seq -ext sql --dir model/migrations/ init_xxx_table
```

forward

```bash
migrate -database "sqlite3://dev.db" --path model/migrations/ up
# or use make target (defaults to dev environment)
make migrate
# or explicitly specify environment
make migrate env=prod
make migrate env=dev
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
