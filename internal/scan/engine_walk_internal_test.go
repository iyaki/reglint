package scan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollectRootFilesUsesConfiguredWalker(t *testing.T) {
	root := t.TempDir()
	called := false

	originalWalker := walkRootTree
	walkRootTree = func(rootPath string, fn filepath.WalkFunc) error {
		called = true
		info, err := os.Stat(rootPath)
		if err != nil {
			return err
		}

		return fn(rootPath, info, nil)
	}
	t.Cleanup(func() {
		walkRootTree = originalWalker
	})

	files, skipped, err := collectRootFiles(root, []string{"**/*"}, nil, nil, nil, 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected configured walker to be called")
	}
	if len(files) != 0 {
		t.Fatalf("expected no files, got %v", files)
	}
	if skipped != 0 {
		t.Fatalf("expected no skipped files, got %d", skipped)
	}
}
