# RegLint

RegLint is a regex-based linter that scans codebases using YAML-defined rules and emits console, JSON, or SARIF outputs for local and CI workflows.

## Quickstart

```bash
reglint init
reglint analyze --config reglint-rules.yaml
```

## CLI

```bash
reglint <command> [flags]

Commands:
  analyze (alias: analyse)
  init
```

## Config

The default rules file name is `reglint-rules.yaml`.

```bash
reglint init --out configs/reglint-rules.yaml
reglint analyze --config configs/reglint-rules.yaml
```

Console output uses ANSI severity colors by default. You can disable colors in the rules file or for a single run with `NO_COLOR`.

```yaml
consoleColorsEnabled: false
rules:
  - message: "Found token $0"
    regex: "token=[a-z]+"
    severity: "error"
```

```bash
NO_COLOR=1 reglint analyze --config reglint-rules.yaml --format console
```

## Baseline

Use baseline compare mode to suppress known findings keyed by `(filePath, message)` with count-based tolerance.

```bash
reglint analyze --config testdata/rules/fail.yaml --baseline testdata/baseline/valid-equal.json testdata/fixtures
```

This command exits `0` because the fixture match is already covered by the baseline.

You can also set a default baseline in config and run without `--baseline`:

```bash
reglint analyze --config testdata/rules/baseline.yaml testdata/fixtures
```

To regenerate a baseline from full current findings:

```bash
reglint analyze --config testdata/rules/fail.yaml --baseline testdata/baseline/generated.json --write-baseline testdata/fixtures
```

This command exits `0` on successful write, even when matches exist.

To see baseline validation behavior with invalid data:

```bash
reglint analyze --config testdata/rules/fail.yaml --baseline testdata/baseline/invalid-duplicate-keys.json testdata/fixtures
```

This command exits `1` with a single baseline validation error.

## Outputs

- `console` writes to stdout.
- `json` writes to stdout when it is the only format; otherwise set `--out-json`.
- `sarif` writes to stdout when it is the only format; otherwise set `--out-sarif`.

```bash
reglint analyze --config reglint-rules.yaml --format console,json --out-json /tmp/scan.json
reglint analyze --config reglint-rules.yaml --format sarif --out-sarif /tmp/scan.sarif
```
