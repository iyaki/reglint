//nolint:testpackage
package scan

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCollectFilesFiltersByIncludeExclude(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFile(t, filepath.Join(root, "src", "main.go"))
	writeFile(t, filepath.Join(root, "src", "vendor", "skip.go"))
	writeFile(t, filepath.Join(root, "docs", "notes.md"))
	writeFile(t, filepath.Join(root, "README.md"))

	files, err := collectFiles([]string{root}, []string{"src/**", "docs/**"}, []string{"**/vendor/**"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"docs/notes.md", "src/main.go"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("expected files %v, got %v", want, files)
	}
}

func TestMatchesPathHonorsExclude(t *testing.T) {
	t.Parallel()

	matched, err := matchesPath("src/vendor/skip.go", []string{"src/**"}, []string{"**/vendor/**"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matched {
		t.Fatal("expected excluded path to be skipped")
	}
}

func TestMatchesPathRequiresInclude(t *testing.T) {
	t.Parallel()

	_, err := matchesPath("src/main.go", nil, nil)
	if err == nil {
		t.Fatal("expected error for empty include patterns")
	}
}

func TestMatchesPathRejectsInvalidPattern(t *testing.T) {
	t.Parallel()

	_, err := matchesPath("src/main.go", []string{"["}, nil)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNormalizePatternsTrimsAndDropsEmpty(t *testing.T) {
	t.Parallel()

	got := normalizePatterns([]string{" src/** ", "", "  ", "docs/**"})
	want := []string{"src/**", "docs/**"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected patterns %v, got %v", want, got)
	}
}

func writeFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("data"), 0o600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}
