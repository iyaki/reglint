//nolint:testpackage
package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/iyaki/regex-checker/internal/scan"
)

func TestWriteJSONNoMatches(t *testing.T) {
	t.Parallel()

	result := scan.Result{
		Matches: nil,
		Stats: scan.Stats{
			FilesScanned: 0,
			FilesSkipped: 0,
			Matches:      0,
			DurationMs:   0,
		},
	}

	var buffer bytes.Buffer
	if err := WriteJSON(result, &buffer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got jsonResult
	if err := json.Unmarshal(buffer.Bytes(), &got); err != nil {
		t.Fatalf("failed to parse json output: %v", err)
	}

	if got.SchemaVersion != 1 {
		t.Fatalf("unexpected schema version: %d", got.SchemaVersion)
	}
	if len(got.Matches) != 0 {
		t.Fatalf("expected no matches, got %d", len(got.Matches))
	}
	if got.Stats.FilesScanned != 0 || got.Stats.FilesSkipped != 0 || got.Stats.Matches != 0 || got.Stats.DurationMs != 0 {
		t.Fatalf("unexpected stats: %+v", got.Stats)
	}
}

func TestWriteJSONOrdersMatches(t *testing.T) {
	t.Parallel()

	result := scan.Result{
		Matches: sampleMatches(),
		Stats: scan.Stats{
			FilesScanned: 2,
			FilesSkipped: 1,
			Matches:      5,
			DurationMs:   12,
		},
	}

	var buffer bytes.Buffer
	if err := WriteJSON(result, &buffer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got jsonResult
	if err := json.Unmarshal(buffer.Bytes(), &got); err != nil {
		t.Fatalf("failed to parse json output: %v", err)
	}

	if got.SchemaVersion != 1 {
		t.Fatalf("unexpected schema version: %d", got.SchemaVersion)
	}
	if len(got.Matches) != 5 {
		t.Fatalf("unexpected match count: %d", len(got.Matches))
	}

	assertFirstMatch(t, got.Matches[0])
	assertWarningOrdering(t, got.Matches[1], got.Matches[2])
	assertInfoMatch(t, got.Matches[3])
	assertLastMatch(t, got.Matches[4])
}

func sampleMatches() []scan.Match {
	return []scan.Match{
		{
			Message:   "Warn msg",
			Severity:  "warning",
			FilePath:  "b/file.go",
			Line:      2,
			Column:    3,
			MatchText: "warn",
		},
		{
			Message:   "Info msg",
			Severity:  "info",
			FilePath:  "a/file.go",
			Line:      10,
			Column:    1,
			MatchText: "info",
		},
		{
			Message:   "Error msg",
			Severity:  "error",
			FilePath:  "a/file.go",
			Line:      2,
			Column:    5,
			MatchText: "error",
		},
		{
			Message:   "Zulu warn",
			Severity:  "warning",
			FilePath:  "a/file.go",
			Line:      2,
			Column:    5,
			MatchText: "zulu",
		},
		{
			Message:   "Alpha warn",
			Severity:  "warning",
			FilePath:  "a/file.go",
			Line:      2,
			Column:    5,
			MatchText: "alpha",
		},
	}
}

func assertFirstMatch(t *testing.T, match scan.Match) {
	t.Helper()

	if match.FilePath != "a/file.go" {
		t.Fatalf("unexpected first file path: %s", match.FilePath)
	}
	if match.Line != 2 || match.Column != 5 {
		t.Fatalf("unexpected first location: %d:%d", match.Line, match.Column)
	}
	if match.Severity != "error" {
		t.Fatalf("unexpected first severity: %s", match.Severity)
	}
}

func assertWarningOrdering(t *testing.T, first scan.Match, second scan.Match) {
	t.Helper()

	if first.Message != "Alpha warn" || second.Message != "Zulu warn" {
		t.Fatalf("unexpected warning ordering: %s, %s", first.Message, second.Message)
	}
}

func assertInfoMatch(t *testing.T, match scan.Match) {
	t.Helper()

	if match.Message != "Info msg" {
		t.Fatalf("unexpected info ordering: %+v", match)
	}
}

func assertLastMatch(t *testing.T, match scan.Match) {
	t.Helper()

	if match.FilePath != "b/file.go" {
		t.Fatalf("unexpected last match: %+v", match)
	}
}
