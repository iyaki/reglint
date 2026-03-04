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

	files, skipped, err := collectFiles([]string{root}, []string{"src/**", "docs/**"}, []string{"**/vendor/**"}, 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skipped != 0 {
		t.Fatalf("expected 0 skipped files, got %d", skipped)
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

func TestCollectFilesSkipsLargeFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFileWithContent(t, filepath.Join(root, "src", "small.txt"), "data")
	writeFileWithContent(t, filepath.Join(root, "src", "large.txt"), "0123456789")

	files, skipped, err := collectFiles([]string{root}, []string{"src/**"}, nil, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skipped != 1 {
		t.Fatalf("expected 1 skipped file, got %d", skipped)
	}

	want := []string{"src/small.txt"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("expected files %v, got %v", want, files)
	}
}

func TestCollectFilesSkipsBinaryFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeFileWithContent(t, filepath.Join(root, "src", "text.txt"), "hello")
	writeFileBytes(t, filepath.Join(root, "src", "binary.bin"), []byte{0x00, 0x01, 0x02})

	files, skipped, err := collectFiles([]string{root}, []string{"src/**"}, nil, 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if skipped != 1 {
		t.Fatalf("expected 1 skipped file, got %d", skipped)
	}

	want := []string{"src/text.txt"}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("expected files %v, got %v", want, files)
	}
}

func writeFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	writeFileWithContent(t, path, "data")
}

func writeFileWithContent(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}

func writeFileBytes(t *testing.T, path string, content []byte) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}
