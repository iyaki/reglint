# Console Formatter

Status: Proposed

## Overview

### Purpose

- Provide human-readable scan results on stdout for local usage.
- Keep output stable and grep-friendly across platforms.

### Goals

- Deterministic ordering of files and matches.
- Clear severity labeling and location formatting.
- Single-line summary with scan stats.

### Non-Goals

- Colored or interactive output.
- Machine-optimized formats (see JSON/SARIF formatter specs).
- Streaming output per file while scanning (render after scan completes).

### Scope

- Render a `ScanResult` to stdout when `--format console` is selected.
- Use only data already present in `ScanResult`.

## Architecture

### Module/package layout (tree format)

```
internal/
  output/
    console.go
```

### Component diagram (ASCII)

```
[ScanResult] -> [Console Formatter] -> stdout
```

### Data flow summary

1. Receive `ScanResult` from the scan service.
2. Sort matches deterministically (see Ordering).
3. Group matches by `filePath`.
4. Render each file section and its matches.
5. Render the summary line.

## Data model

### Core Entities

ConsoleMatchLine

- Definition: A rendered line for a single match.
- Fields:
  - `filePath` (string): Relative file path used in the header.
  - `severityLabel` (string): One of `ERROR|WARN|NOTICE|INFO`.
  - `line` (int): 1-based line number.
  - `column` (int): 1-based column number (rune index).
  - `message` (string): Interpolated match message.
  - `fileUri` (string): Absolute file URI suffix in the format `file://<abs-path>:<line>`.

ConsoleSummary

- Definition: Summary of the scan.
- Fields:
  - `filesScanned` (int)
  - `filesSkipped` (int)
  - `matches` (int)
  - `durationMs` (int64)

### Relationships

- `ConsoleMatchLine` derives from `Match` in `specs/data-model.md`.
- `ConsoleSummary` derives from `ScanStats` in `specs/data-model.md`.

### Persistence Notes

- No persistence. Output is written to stdout.

## Workflows

### Render matches (happy path)

1. Sort matches (Ordering).
2. For each `filePath`, print the file header line.
3. For each match in that file, print a match line with the `fileUri` suffix.
4. After all files, print the summary line.

### No matches

1. Print `No matches found.`
2. Print the summary line.

## APIs

- Internal writer interface only. No network APIs.
- Suggested signature: `WriteConsole(result ScanResult, out io.Writer) error`.

## Client SDK Design

- No client SDK. Formatter is internal only.

## Configuration

- CLI flag: `--format console` (default).
- Output destination: stdout only.

## Permissions

- No permissions or authentication.

## Security Considerations

- Output contains interpolated messages that may include sensitive data.
- Do not emit raw `matchText` in console output.
- Output includes absolute file URIs, which may reveal local filesystem layout.

## Dependencies

- Standard library only.

## Open Questions / Risks

- None.

## Verifications

- Console output is deterministic across runs with identical inputs.
- `No matches found.` appears when `matches == 0`.
- Summary line includes files scanned, skipped, total matches, and duration.

## Appendices

### Ordering

- Sort by `filePath` (ascending, byte-wise).
- Within a file, sort by `line` then `column`.
- If ties remain, sort by severity order `error > warning > notice > info`.
- If ties remain, sort by `message` (ascending, byte-wise).

### Output format

```
path/to/file.ext
  ERROR 12:5 Avoid hardcoded token: abc123 file:///abs/path/to/file.ext:12
  WARN  42:1 Unexpected debug flag file:///abs/path/to/file.ext:42

another/file.go
  NOTICE 3:9 Use of TODO comment file:///abs/another/file.go:3

Summary: files=10 skipped=1 matches=3 durationMs=120
```

### File URI format

- Scheme is fixed to `file` (no configuration).
- Path is absolute and URL-encoded (spaces become `%20`).
- The line number is appended as `:<line>` (no column).
- Windows drive letters use a leading slash (example: `file:///C:/path/to/file.go:12`).
