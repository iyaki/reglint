# Implementation Plan (short-flags)

**Status:** Short flags for analyze complete; mutation coverage expanded.
**Last Updated:** 2026-03-07
**Primary Specs:** `specs/cli-help.md`, `specs/cli-analyze.md`, `specs/cli.md`

## Quick Reference

| System / Subsystem              | Specs                              | Modules / Packages                                               | Artifacts |
| ------------------------------- | ---------------------------------- | ---------------------------------------------------------------- | --------- |
| CLI flag parsing (analyze)      | `specs/cli-analyze.md`             | `internal/cli/analyze.go` ✅                                     | N/A       |
| CLI help topics + renderer      | `specs/cli-help.md`                | `internal/cli/help.go` ✅, `internal/cli/cli.go` ✅              | N/A       |
| CLI short-flag coverage (tests) | `specs/testing-and-validations.md` | `internal/cli/analyze_test.go` ✅, `internal/cli/cli_test.go` ✅ | N/A       |

## Phase 9: Analyze short flags

**Goal:** Support `-c`/`-f` short flags for analyze command parsing and help output.
**Status:** Complete
**Paths:** `internal/cli/analyze.go`, `internal/cli/help.go`, `internal/cli/*_test.go`
**Reference pattern:** `internal/cli/analyze.go`

### 9.1 Parse short flags for analyze

- [x] Add `-c` alias for `--config` in analyze flag parsing.
- [x] Add `-f` alias for `--format` in analyze flag parsing.

### 9.2 Help output includes short flags

- [x] Include `-c` short flag on `--config` in analyze help flags.
- [x] Include `-f` short flag on `--format` in analyze help flags.

### 9.3 Tests for short flags

- [x] Add/extend tests for parsing `-c` and `-f` (analyze args).
- [x] Update help output expectations to include `-c` and `-f`.

**Definition of Done**

- `reglint analyze -c <path>` is accepted and equivalent to `--config`.
- `reglint analyze -f json` is accepted and equivalent to `--format`.
- `reglint analyze --help` includes `-c, --config` and `-f, --format`.

**Risks/Dependencies**

- Keep help topics aligned with analyze spec to avoid drift.

## Verification Log

- 2026-03-06: `git log -n 10 -- specs/cli.md specs/cli-help.md specs/cli-analyze.md specs/cli-init.md` - reviewed recent CLI spec changes for short flags.
- 2026-03-06: Read `specs/cli-help.md`, `specs/cli-analyze.md`, `specs/cli.md` - confirmed short-flag requirements.
- 2026-03-06: Read `internal/cli/analyze.go`, `internal/cli/help.go`, `internal/cli/cli.go` - verified short flags absent in parsing and help.
- 2026-03-06: Read `internal/cli/cli_test.go`, `internal/cli/analyze_test.go` - verified tests do not cover `-c`/`-f`.
- 2026-03-06: `go test ./internal/cli` - passed.
- 2026-03-07: `go test ./internal/cli` - passed.

## Summary

| Phase                        | Status   |
| ---------------------------- | -------- |
| Phase 9: Analyze short flags | Complete |

**Remaining effort:** None.

## Known Existing Work

- Analyze parsing and help rendering are implemented for long-form flags in `internal/cli/analyze.go` and `internal/cli/help.go`.
- Help routing already honors `-h` for root and subcommands via `internal/cli/cli.go`.

## Manual Deployment Tasks

None.
