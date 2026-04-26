# Remove Portable Temp Password CLI

## Problem

The repository still exposes a standalone temp-password CLI (`cmd/temp-password`) and Makefile build targets for producing a binary that generates password material without writing to SQLite. That behavior is no longer wanted. The desired state is to keep only the DB-backed temp-password flow and remove the portable CLI surface that can mislead operators into producing non-persisted passwords.

## Proposed Approach

Remove only the portable CLI surface:

1. Delete the standalone command entrypoint under `cmd/temp-password`
2. Remove the Makefile targets that build the portable binary
3. Remove README guidance that documents the portable CLI

Keep `pkg/temppassword` intact because the DB-backed flow in `service/passkey.go` still depends on it for password generation and hashing.

## In Scope

- Remove `build-temp-password` and `build-temp-password-lite` from `Makefile`
- Remove the `cmd/temp-password/` command
- Remove README sections and examples that advertise the standalone non-DB temp-password CLI
- Update nearby Makefile housekeeping (`.PHONY`, `clean`) so it no longer references the removed binary output

## Out of Scope

- Refactoring or deleting `pkg/temppassword` wholesale
- Changing the DB-backed temp-password workflow
- Changing passkey registration, service, model, or router behavior
- Broad documentation cleanup beyond the portable CLI references

## Cleanup Rules

- Prefer deleting the unused CLI surface instead of turning it into a wrapper or stub
- Preserve any shared password-generation logic still used by DB-backed code
- Keep changes surgical: remove only references tied to the portable CLI

## Verification

Verification should confirm:

- `Makefile` no longer advertises or builds the portable binary
- `README.md` no longer documents the standalone temp-password flow
- the removed command directory is gone
- remaining temp-password-related code still references only the DB-backed path

## Expected Outcome

After this change, the repository should no longer offer a non-persisting temp-password command path. Operators should only see the DB-backed flow, while shared password-generation helpers remain available internally for the code that still needs them.
