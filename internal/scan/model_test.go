//nolint:testpackage
package scan

import (
	"testing"
)

func TestScanModelsHoldFields(t *testing.T) {
	t.Parallel()

	match := Match{
		Message:   "Avoid hardcoded token: abc123",
		Severity:  "error",
		FilePath:  "src/auth/token.go",
		Line:      12,
		Column:    5,
		MatchText: "token=abc123",
	}
	stats := Stats{
		FilesScanned: 10,
		FilesSkipped: 1,
		Matches:      1,
		DurationMs:   120,
	}
	result := Result{
		Matches: []Match{match},
		Stats:   stats,
	}

	if result.Matches[0].Message != "Avoid hardcoded token: abc123" {
		t.Fatalf("unexpected match message: %s", result.Matches[0].Message)
	}
	if result.Stats.FilesScanned != 10 {
		t.Fatalf("unexpected files scanned: %d", result.Stats.FilesScanned)
	}
}
