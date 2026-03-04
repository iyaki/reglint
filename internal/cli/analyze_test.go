package cli_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/iyaki/regex-checker/internal/cli"
)

func TestParseAnalyzeDefaults(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	got, err := cli.ParseAnalyzeArgs([]string{"--config", configPath})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assertDefaultConfig(t, got, configPath)
}

func TestParseAnalyzeFormatsDedupAndTrim(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	got, err := cli.ParseAnalyzeArgs([]string{"--config", configPath, "--format", "json, json"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := []string{"json"}
	if len(got.Formats) != len(want) {
		t.Fatalf("expected formats %v, got %v", want, got.Formats)
	}
	for i := range want {
		if got.Formats[i] != want[i] {
			t.Fatalf("expected formats %v, got %v", want, got.Formats)
		}
	}
}

func TestParseAnalyzeIncludeExclude(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	got, err := cli.ParseAnalyzeArgs([]string{
		"--config", configPath,
		"--include", "**/*.go",
		"--include", "**/*.md",
		"--exclude", "**/vendor/**",
		"--exclude", "**/.git/**",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(got.Include) != 2 || got.Include[0] != "**/*.go" || got.Include[1] != "**/*.md" {
		t.Fatalf("unexpected include list: %v", got.Include)
	}
	if len(got.Exclude) != 2 || got.Exclude[0] != "**/vendor/**" || got.Exclude[1] != "**/.git/**" {
		t.Fatalf("unexpected exclude list: %v", got.Exclude)
	}
}

func TestParseAnalyzeInvalidFormat(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	_, err := cli.ParseAnalyzeArgs([]string{"--config", configPath, "--format", "bogus"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseAnalyzeInvalidFailOn(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	_, err := cli.ParseAnalyzeArgs([]string{"--config", configPath, "--fail-on", "fatal"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseAnalyzeRequiresOutPathForMultiFormat(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	_, err := cli.ParseAnalyzeArgs([]string{"--config", configPath, "--format", "console,json"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseAnalyzeAllowsSingleJsonToStdout(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	_, err := cli.ParseAnalyzeArgs([]string{"--config", configPath, "--format", "json"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestParseAnalyzeAllowsSingleSarifToStdout(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	_, err := cli.ParseAnalyzeArgs([]string{"--config", configPath, "--format", "sarif"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func assertDefaultConfig(t *testing.T, got cli.Config, configPath string) {
	t.Helper()

	if got.ConfigPath != configPath {
		t.Fatalf("expected config path %q, got %q", configPath, got.ConfigPath)
	}
	if len(got.Roots) != 1 || got.Roots[0] != "." {
		t.Fatalf("expected default root '.', got %v", got.Roots)
	}
	if len(got.Formats) != 1 || got.Formats[0] != "console" {
		t.Fatalf("expected default format [console], got %v", got.Formats)
	}
	if got.Concurrency != runtime.GOMAXPROCS(0) {
		t.Fatalf("expected concurrency %d, got %d", runtime.GOMAXPROCS(0), got.Concurrency)
	}
	if got.MaxFileSizeBytes != 5242880 {
		t.Fatalf("expected max file size 5242880, got %d", got.MaxFileSizeBytes)
	}
	if got.OutJSON != "" || got.OutSARIF != "" {
		t.Fatalf("expected empty output paths, got out-json=%q out-sarif=%q", got.OutJSON, got.OutSARIF)
	}
}

func TestParseAnalyzeRejectsEmptyFormatValue(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	_, err := cli.ParseAnalyzeArgs([]string{"--config", configPath, "--format", ""})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseAnalyzeRequiresConfigFile(t *testing.T) {
	t.Parallel()

	missingPath := filepath.Join(t.TempDir(), "missing.yaml")

	_, err := cli.ParseAnalyzeArgs([]string{"--config", missingPath})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseAnalyzeRejectsZeroConcurrency(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	_, err := cli.ParseAnalyzeArgs([]string{"--config", configPath, "--concurrency", "0"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseAnalyzeRejectsZeroMaxFileSize(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfig(t)

	_, err := cli.ParseAnalyzeArgs([]string{"--config", configPath, "--max-file-size", "0"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func writeTempConfig(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "rules.yaml")
	if err := os.WriteFile(path, []byte("rules: []"), 0o600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	return path
}
