# RegLint

RegLint is a regex-based linter for source repositories. It scans files using YAML-defined rules and emits `console`, `json`, or `sarif` output for local development and CI pipelines.

## Quickstart

Create a starter rules file, then run a scan:

```bash
reglint init
reglint analyze --config reglint-rules.yaml
```

Use `reglint --help` or `reglint analyze --help` for the full command reference.

## CLI Overview

```bash
reglint <command> [flags]

Commands:
  analyze (alias: analyse)
  init
```

Common usage patterns:

```bash
# Analyze current directory with default config path
reglint analyze

# Analyze specific roots
reglint analyze services/api services/web

# Generate config at a custom location
reglint init --out configs/reglint-rules.yaml
reglint analyze --config configs/reglint-rules.yaml
```

## Exit Codes

- `0`: command succeeded and no `--fail-on` threshold was triggered.
- `1`: command/config/runtime error (for example invalid flags or invalid baseline file).
- `2`: command succeeded, but at least one finding met `--fail-on` severity.

## Configuration

Default config path is `reglint-rules.yaml`.

Top-level fields:

- `rules` (required): list of regex rules.
- `include` / `exclude`: repository-level glob controls.
- `failOn`: one of `error`, `warning`, `notice`, `info`.
- `concurrency`: worker count override.
- `baseline`: default baseline file path.
- `git`: default Git scan settings.
- `consoleColorsEnabled`: enable or disable ANSI color in console output.
- `ignoreFilesEnabled`: enable or disable ignore file processing.
- `ignoreFiles`: custom ignore file list.

Minimal example:

```yaml
include:
  - "**/*"
exclude:
  - "**/.git/**"
  - "**/node_modules/**"
failOn: "error"
consoleColorsEnabled: true
rules:
  - message: "Avoid hardcoded token: $1"
    regex: "token\\s*[:=]\\s*([A-Za-z0-9_-]+)"
    severity: "error"
    paths:
      - "src/**"
```

Console output uses ANSI severity colors by default. You can disable colors in config or for a single run with `NO_COLOR`:

```bash
NO_COLOR=1 reglint analyze --config reglint-rules.yaml --format console
```

## Output Formats

- `console` writes to stdout.
- `json` writes to stdout only when it is the single selected format.
- `sarif` writes to stdout only when it is the single selected format.
- When combining multiple formats, use `--out-json` and/or `--out-sarif` as needed.

```bash
reglint analyze --format console,json --out-json /tmp/scan.json
reglint analyze --format sarif --out-sarif /tmp/scan.sarif
```

## Baseline Workflow

Baseline compare mode suppresses known findings using `(filePath, message)` with count-based tolerance.

Compare against an existing baseline:

```bash
reglint analyze --config testdata/rules/fail.yaml --baseline testdata/baseline/valid-equal.json testdata/fixtures
```

Use baseline from config (without passing `--baseline`):

```bash
reglint analyze --config testdata/rules/baseline.yaml testdata/fixtures
```

Generate or refresh a baseline from current findings:

```bash
reglint analyze --config testdata/rules/fail.yaml --baseline testdata/baseline/generated.json --write-baseline testdata/fixtures
```

`--write-baseline` exits `0` on successful write, even if findings exist.

## Git-Scoped Scans

Git integration is optional and defaults to `off`.

```bash
reglint analyze --config reglint-rules.yaml
```

Scan staged files only:

```bash
reglint analyze --config reglint-rules.yaml --git-mode staged
```

Scan files selected by a diff target:

```bash
reglint analyze --config reglint-rules.yaml --git-mode diff --git-diff HEAD~1..HEAD
```

`--git-diff` implies `--git-mode diff` if `--git-mode` is not provided.

Restrict reporting to matches on added lines:

```bash
reglint analyze --config reglint-rules.yaml --git-mode diff --git-diff HEAD~1..HEAD --git-added-lines-only
```

Ignore-file behavior:

- Ignore files are enabled by default in all modes (`off`, `staged`, `diff`).
- Default evaluation order is `.gitignore`, then `.ignore`, then `.reglintignore`.
- Use `--no-gitignore` to disable only `.gitignore`.
- Use `--no-ignore-files` to disable all ignore-file processing.

## Development

Run quick end-to-end smoke coverage:

```bash
make test-e2e-smoke
```

Run full end-to-end matrix:

```bash
make test-e2e
```

Run full local quality checks:

```bash
make quality
```
