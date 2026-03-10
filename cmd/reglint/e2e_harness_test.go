package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestE2EHarnessRunExecutesCompiledBinary(t *testing.T) {
	harness := newE2EHarness(t)

	workspace := t.TempDir()
	configPath := filepath.Join(workspace, "rules.yaml")
	scenario := e2EScenario{
		ID:           "E2E-HARNESS-001",
		Tier:         "smoke",
		Name:         "init command creates a config file",
		Fixture:      workspace,
		Command:      []string{"init", "--out", configPath},
		ExpectedExit: 0,
	}

	result := harness.mustRunScenario(t, scenario)
	harness.assertScenarioStdoutContains(t, scenario, result, "Wrote default config to "+configPath)
	harness.assertScenarioStderrEmpty(t, scenario, result)
	harness.assertScenarioPathExists(t, scenario, result, configPath)
}

func TestE2EHarnessBuildsBinaryOnlyOnce(t *testing.T) {
	harness1 := newE2EHarness(t)
	workspace1 := t.TempDir()
	path1 := filepath.Join(workspace1, "rules.yaml")
	scenario1 := e2EScenario{
		ID:           "E2E-HARNESS-002",
		Tier:         "smoke",
		Name:         "first run uses built binary",
		Fixture:      workspace1,
		Command:      []string{"init", "--out", path1},
		ExpectedExit: 0,
	}
	harness1.mustRunScenario(t, scenario1)

	buildCountAfterFirst := e2eBinaryBuildInvocations()

	harness2 := newE2EHarness(t)
	workspace2 := t.TempDir()
	path2 := filepath.Join(workspace2, "rules.yaml")
	scenario2 := e2EScenario{
		ID:           "E2E-HARNESS-003",
		Tier:         "smoke",
		Name:         "second run reuses built binary",
		Fixture:      workspace2,
		Command:      []string{"init", "--out", path2},
		ExpectedExit: 0,
	}
	harness2.mustRunScenario(t, scenario2)

	buildCountAfterSecond := e2eBinaryBuildInvocations()
	if buildCountAfterSecond != buildCountAfterFirst {
		t.Fatalf("expected binary build count to stay at %d, got %d", buildCountAfterFirst, buildCountAfterSecond)
	}
}

func TestE2EHarnessFailureDiagnosticsIncludeScenarioMetadata(t *testing.T) {
	t.Parallel()

	harness := &e2EHarness{binaryPath: filepath.Join(t.TempDir(), "reglint binary")}
	scenario := e2EScenario{
		ID:           "E2E-SMOKE-999",
		Tier:         "smoke",
		Name:         "metadata diagnostics",
		Fixture:      filepath.Join(t.TempDir(), "fixture workspace"),
		Command:      []string{"analyze", "--config", "rules.yaml", "./fixtures"},
		Env:          map[string]string{"NO_COLOR": "1", "TEST_MODE": "e2e"},
		ExpectedExit: 2,
	}
	result := e2EProcessResult{ExitCode: 1, Stdout: "stdout-body", Stderr: "stderr-body"}

	diagnostic := harness.scenarioFailureDiagnostic(scenario, result, "expected exit code 2, got 1")

	if !strings.Contains(diagnostic, "scenario: E2E-SMOKE-999") {
		t.Fatalf("expected scenario id in diagnostic, got %q", diagnostic)
	}
	if !strings.Contains(diagnostic, "fixture: "+scenario.Fixture) {
		t.Fatalf("expected fixture path in diagnostic, got %q", diagnostic)
	}
	if !strings.Contains(diagnostic, "replay: ") {
		t.Fatalf("expected replay command in diagnostic, got %q", diagnostic)
	}
	if !strings.Contains(diagnostic, "NO_COLOR=") || !strings.Contains(diagnostic, "TEST_MODE=") {
		t.Fatalf("expected env overrides in replay command, got %q", diagnostic)
	}
	if !strings.Contains(diagnostic, "stdout: \"stdout-body\"") {
		t.Fatalf("expected stdout capture in diagnostic, got %q", diagnostic)
	}
	if !strings.Contains(diagnostic, "stderr: \"stderr-body\"") {
		t.Fatalf("expected stderr capture in diagnostic, got %q", diagnostic)
	}
}
