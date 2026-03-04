package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/iyaki/regex-checker/internal/cli"
)

func TestRunShowsHelpWhenNoArgs(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	code := run([]string{}, &output)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}

	if !strings.Contains(output.String(), "Usage:") {
		t.Fatalf("expected usage help, got %q", output.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	code := run([]string{"bogus"}, &output)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}

	if output.String() != "Unknown command: bogus\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestRunRoutesAnalyze(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	code := run([]string{"analyze"}, &output)

	if code != 1 {
		t.Fatalf("expected exit code 1 for missing handler, got %d", code)
	}
	if output.String() != "Unknown command: analyze\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestRunRoutesAnalyseAlias(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	code := run([]string{"analyse"}, &output)

	if code != 1 {
		t.Fatalf("expected exit code 1 for missing handler, got %d", code)
	}
	if output.String() != "Unknown command: analyse\n" {
		t.Fatalf("unexpected output: %q", output.String())
	}
}

func TestRunUsesProvidedOutputWriter(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	code := cli.Run([]string{}, map[string]cli.Handler{}, &output)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(output.String(), "Usage:") {
		t.Fatalf("expected usage help, got %q", output.String())
	}
}

func TestMainExitsWithRunCode(t *testing.T) {
	var output bytes.Buffer
	var exitCode int
	var exited bool

	originalArgs := args
	originalOutput := outputWriter
	originalExit := exitFunc
	defer func() {
		args = originalArgs
		outputWriter = originalOutput
		exitFunc = originalExit
	}()

	args = []string{"regex-checker"}
	outputWriter = &output
	exitFunc = func(code int) {
		exited = true
		exitCode = code
	}

	main()

	if !exited {
		t.Fatal("expected main to call exit")
	}
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(output.String(), "Usage:") {
		t.Fatalf("expected usage help, got %q", output.String())
	}
}
