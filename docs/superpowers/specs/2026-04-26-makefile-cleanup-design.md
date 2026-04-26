# Makefile Invalid Target Cleanup

## Problem

The current `Makefile` mixes active targets with stale or inconsistent build metadata. The immediate user goal is to remove targets that are no longer usable, but to do so aggressively enough that obviously outdated Makefile-only residue is cleaned up at the same time.

This cleanup is intentionally limited to `Makefile` itself. Documentation, scripts, and command implementations are out of scope for this change.

## Proposed Approach

Clean `Makefile` in two layers:

1. Remove targets that are currently invalid or no longer directly usable
2. Clean nearby Makefile-only residue that is clearly outdated and can be corrected without refactoring the file

This keeps the change focused while avoiding a half-finished state where the broken target is gone but the file still advertises stale variables, stale `.PHONY` declarations, or cleanup commands that no longer match actual build outputs.

## In Scope

- Remove the invalid `passkey-temp` target if it no longer matches the CLI it invokes
- Remove or correct stale Makefile-only residue tightly coupled to target cleanup, such as:
  - unused or obsolete target-related variables
  - `.PHONY` entries that are out of sync with actual targets
  - `clean` commands that reference outdated output names instead of current ones
- Preserve working targets, including currently valid temp-password build targets, unless inspection proves they are also invalid

## Out of Scope

- Editing `README.md` or other docs
- Refactoring target structure for style only
- Removing targets that still work just because they look old
- Changing command implementations under `cmd/`, `scripts/`, `service/`, or elsewhere unless Makefile verification reveals a directly coupled issue

## Cleanup Rules

- Prefer deletion over compatibility wrappers for targets that are already invalid
- Only remove a target when it is clearly unusable or inconsistent with the current repository state
- Keep working targets intact even if they are older, unless they are tightly coupled to the invalid target cleanup
- Avoid broad renaming or regrouping of the Makefile

## Verification

Verification should focus on the targets touched by the cleanup:

- confirm removed targets are no longer advertised in `.PHONY`
- confirm retained targets still have valid recipes
- run the most relevant Makefile command checks for edited targets
- inspect `make -n` output or direct command execution for any cleaned target paths as needed

## Expected Outcome

After this change, `Makefile` should no longer expose unusable targets, and its local housekeeping (`.PHONY`, related variables, and `clean`) should match the targets that actually remain.
