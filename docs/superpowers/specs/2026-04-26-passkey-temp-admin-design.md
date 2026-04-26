# Passkey Temp Admin CLI Redesign

## Problem

The repository currently has two separate temp-password CLIs:

- `cmd/temp-password`: standalone generator that produces password material only and does not persist anything
- `cmd/passkey-temp-admin`: DB-backed admin CLI with `create|list|revoke`

The desired outcome is a single-purpose, independently compiled tool that reads `config.yml`, connects to the configured SQLite database, generates a temp password, persists it to `temp_passwords`, and prints the plaintext password for immediate use in the passkey registration flow. Pure generation without persistence is out of scope.

## Proposed Approach

Redesign `cmd/passkey-temp-admin` into a create-only CLI.

The command will:

1. Load `config.yml`
2. Initialize logging and the DAO
3. Resolve TTL from `passkey.temp_password.ttl`, optionally overridden by a CLI argument
4. Generate a temp password and persist it through the existing service/model flow
5. Print the plaintext temp password and expiry time

This keeps the tool independently compilable while reducing it to the one workflow that matters for passkey registration.

## Scope

### In scope

- Remove `list` and `revoke` subcommands from `cmd/passkey-temp-admin`
- Make the command create-only
- Preserve DB-backed temp password creation using the existing `service.CreateTempPassword` path
- Continue loading default TTL from `config.yml`
- Allow an explicit TTL override from the command line
- Keep the output suitable for manual copy/paste into the passkey registration UI

### Out of scope

- Generating passkey/WebAuthn credentials offline
- Changing `/passkey/register` or frontend registration behavior
- Refactoring the portable `cmd/temp-password` command as part of this change
- Adding broader temp password management features

## Command Interface

The command becomes:

```bash
go run ./cmd/passkey-temp-admin
go run ./cmd/passkey-temp-admin 900
```

Behavior:

- No positional argument: use `config.yml` `passkey.temp_password.ttl`
- One positional argument: interpret as TTL seconds
- More than one positional argument: fail with usage output
- Non-integer or negative TTL: fail with a clear error

## Data Flow

1. `config.GetConfig("config.yml")` loads configuration
2. Logger initialization follows the current command pattern
3. `model.InitDao(cfg, logger)` opens the configured SQLite database
4. `service.CreateTempPassword(ttl)` generates the plaintext password, hashes it, inserts the record, and returns both the record and plaintext
5. CLI prints:
   - temp password
   - expiry timestamp

The existing `/passkey/register` flow remains unchanged because it already validates against `temp_passwords`.

## Error Handling

The CLI should explicitly fail on:

- missing or invalid `config.yml`
- DAO initialization failure
- invalid TTL input
- temp password generation failure
- database insert failure

Errors should be written to stderr and exit non-zero, following existing command behavior.

## Testing and Verification

Verification should focus on the command entrypoint and unchanged integration with the current temp-password flow:

- command succeeds with default TTL from config
- command succeeds with explicit TTL override
- invalid TTL is rejected
- inserted temp password record is present in SQLite

Existing service/model behavior is reused rather than redesigned.
