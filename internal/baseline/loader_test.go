//nolint:testpackage // Validate internal behavior directly.
package baseline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadParsesDocument(t *testing.T) {
	t.Parallel()

	path := writeBaselineFile(t, `{
		"schemaVersion": 1,
		"entries": [
			{"filePath": "src/a.go", "message": "m1", "count": 2},
			{"filePath": "src/b.go", "message": "m2", "count": 1}
		]
	}`)

	document, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if document.SchemaVersion != 1 {
		t.Fatalf("expected schema version 1, got %d", document.SchemaVersion)
	}
	if len(document.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(document.Entries))
	}
	if document.Entries[0].FilePath != "src/a.go" {
		t.Fatalf("expected first filePath src/a.go, got %q", document.Entries[0].FilePath)
	}
}

func TestLoadRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	path := writeBaselineFile(t, `{
		"schemaVersion": 1,
		"entries": [
	`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "parse baseline") {
		t.Fatalf("expected parse baseline error, got %v", err)
	}
}

func TestLoadRejectsSchemaViolations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name: "missing schema version",
			content: `{
				"entries": []
			}`,
			wantErr: "schemaVersion is required",
		},
		{
			name: "invalid schema version",
			content: `{
				"schemaVersion": 2,
				"entries": []
			}`,
			wantErr: "schemaVersion must be 1",
		},
		{
			name: "missing entries",
			content: `{
				"schemaVersion": 1
			}`,
			wantErr: "entries is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := writeBaselineFile(t, tt.content)
			_, err := Load(path)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

type invalidEntryCase struct {
	name    string
	content string
	wantErr string
}

func TestLoadRejectsInvalidEntries(t *testing.T) {
	t.Parallel()

	for _, tt := range allInvalidEntryCases() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := writeBaselineFile(t, tt.content)
			_, err := Load(path)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func allInvalidEntryCases() []invalidEntryCase {
	all := invalidPathEntryCases()
	all = append(all, invalidMessageAndCountEntryCases()...)
	all = append(all, duplicateEntryCases()...)

	return all
}

func invalidPathEntryCases() []invalidEntryCase {
	return []invalidEntryCase{
		{
			name: "empty file path",
			content: `{
				"schemaVersion": 1,
				"entries": [
					{"filePath": "", "message": "m1", "count": 1}
				]
			}`,
			wantErr: "filePath",
		},
		{
			name: "absolute file path",
			content: `{
				"schemaVersion": 1,
				"entries": [
					{"filePath": "/src/a.go", "message": "m1", "count": 1}
				]
			}`,
			wantErr: "normalized relative path",
		},
		{
			name: "path traversal",
			content: `{
				"schemaVersion": 1,
				"entries": [
					{"filePath": "../a.go", "message": "m1", "count": 1}
				]
			}`,
			wantErr: "normalized relative path",
		},
		{
			name: "non normalized file path",
			content: `{
				"schemaVersion": 1,
				"entries": [
					{"filePath": "src/./a.go", "message": "m1", "count": 1}
				]
			}`,
			wantErr: "normalized relative path",
		},
	}
}

func invalidMessageAndCountEntryCases() []invalidEntryCase {
	return []invalidEntryCase{
		{
			name: "empty message",
			content: `{
				"schemaVersion": 1,
				"entries": [
					{"filePath": "src/a.go", "message": " ", "count": 1}
				]
			}`,
			wantErr: "message is required",
		},
		{
			name: "non positive count",
			content: `{
				"schemaVersion": 1,
				"entries": [
					{"filePath": "src/a.go", "message": "m1", "count": 0}
				]
			}`,
			wantErr: "count must be positive",
		},
	}
}

func duplicateEntryCases() []invalidEntryCase {
	return []invalidEntryCase{
		{
			name: "duplicate key",
			content: `{
				"schemaVersion": 1,
				"entries": [
					{"filePath": "src/a.go", "message": "m1", "count": 1},
					{"filePath": "src/a.go", "message": "m1", "count": 2}
				]
			}`,
			wantErr: "duplicate baseline entry",
		},
	}
}

func writeBaselineFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write baseline: %v", err)
	}

	return path
}
