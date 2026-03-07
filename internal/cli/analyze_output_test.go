//nolint:testpackage
package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/iyaki/reglint/internal/output"
	"github.com/iyaki/reglint/internal/rules"
	"github.com/iyaki/reglint/internal/scan"
)

var analyzeCwdMutex sync.Mutex

func TestWriteJSONOutputRequiresPathForMultipleFormats(t *testing.T) {
	t.Parallel()

	cfg := Config{Formats: []string{"console", "json"}}

	if err := writeJSONOutput(cfg, scan.Result{}, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteJSONOutputToStdout(t *testing.T) {
	t.Parallel()

	cfg := Config{Formats: []string{"json"}}
	buffer := &bytes.Buffer{}

	if err := writeJSONOutput(cfg, scan.Result{}, buffer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buffer.String(), "schemaVersion") {
		t.Fatalf("unexpected stdout output: %q", buffer.String())
	}
}

func TestWriteSARIFOutputRequiresPathForMultipleFormats(t *testing.T) {
	t.Parallel()

	cfg := Config{Formats: []string{"console", "sarif"}}

	if err := writeSARIFOutput(cfg, scan.Result{}, sampleRules(), &bytes.Buffer{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSARIFOutputToStdout(t *testing.T) {
	t.Parallel()

	cfg := Config{Formats: []string{"sarif"}}
	buffer := &bytes.Buffer{}

	if err := writeSARIFOutput(cfg, scan.Result{}, sampleRules(), buffer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buffer.String(), "RegLint") {
		t.Fatalf("unexpected stdout output: %q", buffer.String())
	}
}

func TestRenderOutputsWritesJSONFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "scan.json")
	cfg := Config{Formats: []string{"json"}, OutJSON: path}
	buffer := &bytes.Buffer{}

	if err := renderOutputs(cfg.Formats, sampleRules(), cfg, scan.Result{}, buffer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buffer.Len() != 0 {
		t.Fatalf("expected no stdout output, got %q", buffer.String())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read json output: %v", err)
	}
	if !strings.Contains(string(data), "schemaVersion") {
		t.Fatalf("expected json output, got %q", string(data))
	}
}

func TestRenderOutputsWritesSARIFFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "scan.sarif")
	cfg := Config{Formats: []string{"sarif"}, OutSARIF: path}
	buffer := &bytes.Buffer{}

	if err := renderOutputs(cfg.Formats, sampleRules(), cfg, scan.Result{}, buffer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buffer.Len() != 0 {
		t.Fatalf("expected no stdout output, got %q", buffer.String())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read sarif output: %v", err)
	}
	if !strings.Contains(string(data), "RegLint") {
		t.Fatalf("expected sarif output, got %q", string(data))
	}
}

func TestRenderOutputsWritesConsole(t *testing.T) {
	t.Parallel()

	result := scan.Result{
		Matches: []scan.Match{{Message: "msg", Severity: "error", FilePath: "file.txt", Line: 1, Column: 1}},
		Stats: scan.Stats{
			FilesScanned: 1,
			FilesSkipped: 0,
			Matches:      1,
			DurationMs:   2,
		},
	}
	cfg := Config{Formats: []string{"console"}}
	buffer := &bytes.Buffer{}

	if err := renderOutputs(cfg.Formats, sampleRules(), cfg, result, buffer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buffer.String(), "Summary:") {
		t.Fatalf("expected summary output, got %q", buffer.String())
	}
}

func TestRenderOutputsRejectsUnknownFormat(t *testing.T) {
	t.Parallel()

	cfg := Config{Formats: []string{"bogus"}}

	if err := renderOutputs(cfg.Formats, sampleRules(), cfg, scan.Result{}, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestExitCodeForFailOn(t *testing.T) {
	t.Parallel()

	matches := []scan.Match{{Severity: "warning"}}
	if exitCodeForFailOn(matches, "warning") != exitCodeFailOn {
		t.Fatalf("expected fail-on exit code")
	}
	if exitCodeForFailOn(matches, "error") != 0 {
		t.Fatalf("expected success exit code")
	}
}

func TestRunAnalyzeReturnsScanError(t *testing.T) {
	t.Parallel()

	config := "include:\n  - ''\nrules:\n  - message: 'hello'\n    regex: 'world'\n"
	configPath := writeTempConfigFile(t, config)

	result, failOn, formats, ruleset, cfg, err := runAnalyze([]string{"--config", configPath})
	_ = result
	_ = failOn
	_ = formats
	_ = ruleset
	_ = cfg
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunAnalyzeReturnsConfigLoadError(t *testing.T) {
	t.Parallel()

	configPath := writeTempConfigFile(t, "rules: []")
	if err := os.Chmod(configPath, 0o000); err != nil {
		t.Fatalf("failed to set permissions: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(configPath, 0o600)
	})

	result, failOn, formats, ruleset, cfg, err := runAnalyze([]string{"--config", configPath})
	_ = result
	_ = failOn
	_ = formats
	_ = ruleset
	_ = cfg
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

type errorFormatter struct {
	name string
}

func (e errorFormatter) Name() string {
	return e.name
}

func (e errorFormatter) Write(scan.Result, io.Writer) error {
	return errors.New("write failed")
}

func TestRunAnalyzeShortFlags(t *testing.T) {
	analyzeCwdMutex.Lock()
	defer analyzeCwdMutex.Unlock()

	config := "rules:\n  - message: 'hello'\n    regex: 'world'\n"
	configPath := writeTempConfigFile(t, config)

	current, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to read cwd: %v", err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("failed to change cwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(current)
	})

	result, failOn, formats, ruleset, cfg, err := runAnalyze([]string{"-c", configPath, "-f", "console"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if failOn != "" {
		t.Fatalf("expected empty fail-on, got %q", failOn)
	}
	if len(formats) != 1 || formats[0] != "console" {
		t.Fatalf("expected formats [console], got %v", formats)
	}
	if len(ruleset) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(ruleset))
	}
	if cfg.ConfigPath != configPath {
		t.Fatalf("expected config path %q, got %q", configPath, cfg.ConfigPath)
	}
	if result.Stats.FilesScanned < 0 {
		t.Fatalf("unexpected files scanned: %d", result.Stats.FilesScanned)
	}
}

func TestWriteJSONFileFailsOnDirectory(t *testing.T) {
	t.Parallel()

	path := t.TempDir()
	if err := writeJSONFile(path, scan.Result{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteJSONFileFailsOnMissingParent(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "missing", "scan.json")
	if err := writeJSONFile(path, scan.Result{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSARIFFileFailsOnDirectory(t *testing.T) {
	t.Parallel()

	path := t.TempDir()
	if err := writeSARIFFile(path, scan.Result{}, sampleRules()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSARIFFileFailsOnMissingParent(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "missing", "scan.sarif")
	if err := writeSARIFFile(path, scan.Result{}, sampleRules()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteJSONFileFailsOnReadOnlyDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatalf("failed to set permissions: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(dir, 0o700)
	})

	path := filepath.Join(dir, "scan.json")
	if err := writeJSONFile(path, scan.Result{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSARIFFileFailsOnReadOnlyDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatalf("failed to set permissions: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(dir, 0o700)
	})

	path := filepath.Join(dir, "scan.sarif")
	if err := writeSARIFFile(path, scan.Result{}, sampleRules()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteJSONFileFailsOnWriteError(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "scan.json")
	if err := writeJSONFile(path, scan.Result{Matches: []scan.Match{{FilePath: "", Line: 1}}}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteJSONFileFailsOnCreateError(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "missing", "scan.json")
	if err := writeJSONFile(path, scan.Result{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestWriteJSONPropagatesWriteError(t *testing.T) {
	t.Parallel()

	result := scan.Result{Matches: []scan.Match{{FilePath: "file.txt", Line: 1, Column: 1}}}
	if err := output.WriteJSON(result, failingWriter{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSARIFPropagatesWriteError(t *testing.T) {
	t.Parallel()

	result := scan.Result{Matches: []scan.Match{{FilePath: "file.txt", Line: 1, Column: 1}}}
	if err := output.WriteSARIF(result, sampleRules(), failingWriter{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteJSONFileFailsOnReadOnlyFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "scan.json")
	if err := os.WriteFile(path, []byte("data"), 0o400); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if err := writeJSONFile(path, scan.Result{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSARIFFileFailsOnReadOnlyFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "scan.sarif")
	if err := os.WriteFile(path, []byte("data"), 0o400); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if err := writeSARIFFile(path, scan.Result{}, sampleRules()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

type captureFormatter struct {
	name    string
	written bool
}

func (c *captureFormatter) Name() string {
	return c.name
}

func (c *captureFormatter) Write(scan.Result, io.Writer) error {
	c.written = true

	return nil
}

type mockFormatter struct {
	name string
}

func (m mockFormatter) Name() string {
	return m.name
}

func (m mockFormatter) Write(scan.Result, io.Writer) error {
	return nil
}

func TestRenderOutputsRequiresRegistrySetup(t *testing.T) {
	t.Parallel()

	registry, err := output.NewRegistry(mockFormatter{name: ""})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if registry != nil {
		t.Fatal("expected nil registry")
	}
}

func TestRenderOutputsPropagatesFormatterError(t *testing.T) {
	t.Parallel()

	registry, err := output.NewRegistry(errorFormatter{name: "console"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	formatter, err := registry.ResolveName("console")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := Config{Formats: []string{"console"}}
	if err := renderFormat(formatter, cfg, sampleRules(), scan.Result{}, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRenderOutputsSkipsUnknownFormatWrite(t *testing.T) {
	t.Parallel()

	formatter := &captureFormatter{name: "console"}
	registry, err := output.NewRegistry(formatter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = registry.ResolveName("console")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := Config{Formats: []string{"bogus"}}
	if err := renderOutputs(cfg.Formats, sampleRules(), cfg, scan.Result{}, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error, got nil")
	}
	if formatter.written {
		t.Fatal("expected formatter not to be written")
	}
}

type nilFormatter struct{}

func (nilFormatter) Name() string {
	return "console"
}

func (nilFormatter) Write(scan.Result, io.Writer) error {
	return nil
}

func TestRenderOutputsReturnsErrorWhenResolveFails(t *testing.T) {
	t.Parallel()

	outputRegistry = func([]rules.Rule) (*output.Registry, error) {
		return output.NewRegistry(nilFormatter{})
	}
	t.Cleanup(func() {
		outputRegistry = defaultOutputRegistry
	})

	if err := renderOutputs([]string{"missing"}, sampleRules(), Config{}, scan.Result{}, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRenderOutputsReturnsErrorWhenFormatterFails(t *testing.T) {
	t.Parallel()

	currentRegistry := outputRegistry
	outputRegistry = func([]rules.Rule) (*output.Registry, error) {
		return output.NewRegistry(errorFormatter{name: "console"})
	}
	t.Cleanup(func() {
		outputRegistry = currentRegistry
	})

	if err := renderOutputs([]string{"console"}, sampleRules(), Config{}, scan.Result{}, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRenderOutputsReturnsErrorWhenRegistryFails(t *testing.T) {
	t.Parallel()

	outputRegistry = func([]rules.Rule) (*output.Registry, error) {
		return nil, errors.New("registry failed")
	}
	t.Cleanup(func() {
		outputRegistry = defaultOutputRegistry
	})

	if err := renderOutputs([]string{"console"}, sampleRules(), Config{}, scan.Result{}, &bytes.Buffer{}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBuildEffectiveRulesSkipsOverridesWithoutCLIInput(t *testing.T) {
	t.Parallel()

	effective := rules.RuleSet{
		Include: []string{"**/*.go"},
		Exclude: []string{"**/vendor/**"},
		Rules: []rules.Rule{
			{
				Message: "rule",
				Regex:   "token",
				Paths:   []string{"custom-path"},
				Exclude: []string{"custom-exclude"},
			},
		},
	}

	got := buildEffectiveRules(Config{}, effective)
	if len(got) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(got))
	}
	if len(got[0].Paths) != 1 || got[0].Paths[0] != "custom-path" {
		t.Fatalf("expected paths to remain custom, got %v", got[0].Paths)
	}
	if len(got[0].Exclude) != 1 || got[0].Exclude[0] != "custom-exclude" {
		t.Fatalf("expected exclude to remain custom, got %v", got[0].Exclude)
	}
}

func TestResolveConcurrencyUsesConfigWhenRulesetMissing(t *testing.T) {
	t.Parallel()

	got := resolveConcurrency(Config{Concurrency: 3}, nil)
	if got != 3 {
		t.Fatalf("expected concurrency 3, got %d", got)
	}
}

func TestRenderOutputsReturnsErrorWhenRegistrySetupFails(t *testing.T) {
	t.Parallel()

	if _, err := output.NewRegistry(nilFormatter{}, nil); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func sampleRules() []rules.Rule {
	return []rules.Rule{
		{
			Message:  "rule",
			Regex:    "token",
			Severity: "warning",
			Index:    0,
		},
	}
}

func writeTempConfigFile(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "rules.yaml")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	return path
}
