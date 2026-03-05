# Init Command

Status: Implemented

## Overview

### Purpose

- Generate a default regex rules config file for quick setup.
- Provide a safe, deterministic starting point for new users.

### Goals

- Create a valid YAML rules config that matches `specs/configuration.md` and `specs/regex-rules.md`.
- Avoid overwriting existing files unless explicitly requested.
- Emit a short success message with the output path.

### Non-Goals

- Interactive prompting or wizard flows.
- Validating or scanning files (use `analyze`).
- Managing multiple config profiles.

### Scope

- Command syntax, flags, defaults, and validation rules for init.
- Default file contents and overwrite behavior.

## Architecture

### Module/package layout (tree format)

```
cmd/
  reglint/
    main.go
internal/
  cli/
    init.go
  config/
```

### Component diagram (ASCII)

```
[Init Command] -> [Config Template] -> [Filesystem Writer]
```

### Data flow summary

1. Parse flags.
2. Resolve output path (default `reglint-rules.yaml`).
3. If output exists and `--force` is not set, error.
4. Render default config content.
5. Write file and print a success message.

## Data model

### Core Entities

InitConfig

- Definition: Parsed CLI inputs and resolved defaults for init.
- Fields:
  - `outputPath` (string, required)
  - `force` (bool, required)

Reference struct (Go):

```go
type InitConfig struct {
    OutputPath string
    Force      bool
}
```

## Workflows

### Generate default config (happy path)

1. Parse flags.
2. Resolve output path (default `reglint-rules.yaml`).
3. Check if the output path exists.
4. If exists and `--force` is not set, error and exit 1.
5. Write the default YAML.
6. Print: `Wrote default config to <path>`.

### Validation and errors

- `--out` must be a valid file path.
- If the output file already exists and `--force` is not set, exit with code 1 and explain how to override.
- Any write failure exits with code 1.

### Exit codes

- `0`: config file written successfully.
- `1`: configuration or filesystem error.

## Configuration

### Command syntax

```
reglint init [flags]
```

### Flags

| Flag      | Type   | Required | Default              | Purpose                               |
| --------- | ------ | -------- | -------------------- | ------------------------------------- |
| `--out`   | string | no       | `reglint-rules.yaml` | Output path for the config file.      |
| `--force` | bool   | no       | `false`              | Overwrite if the file already exists. |

### Default config content

```yaml
include:
  - "**/*"
exclude:
  - "**/.git/**"
  - "**/node_modules/**"
  - "**/vendor/**"
failOn: "error"
rules:
  - message: "Avoid hardcoded token: $1"
    regex: "token\\s*[:=]\\s*([A-Za-z0-9_-]+)"
    severity: "error"
    paths:
      - "src/**"
```

Notes:

- The default rule is intentionally simple and safe to run on most repositories.

## Verifications

- `reglint init` writes `reglint-rules.yaml` in the current directory.
- `reglint init --out configs/reglint-rules.yaml` writes the file at the specified path.
- `reglint init` fails if the file exists; `reglint init --force` overwrites it.

## Appendices

### Examples

```
reglint init
reglint init --out configs/reglint-rules.yaml
reglint init --force
```
