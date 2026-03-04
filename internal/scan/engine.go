package scan

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func collectFiles(roots []string, include []string, exclude []string) ([]string, error) {
	var files []string

	for _, root := range roots {
		root = filepath.Clean(root)
		if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			relPath = filepath.ToSlash(relPath)
			if relPath == "." {
				return nil
			}

			match, err := matchesPath(relPath, include, exclude)
			if err != nil {
				return err
			}
			if match {
				files = append(files, relPath)
			}

			return nil
		}); err != nil {
			return nil, err
		}
	}

	sort.Strings(files)

	return files, nil
}

func matchesPath(path string, include []string, exclude []string) (bool, error) {
	if len(include) == 0 {
		return false, errors.New("include patterns required")
	}

	for _, pattern := range exclude {
		ok, err := doublestar.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if ok {
			return false, nil
		}
	}

	for _, pattern := range include {
		ok, err := doublestar.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

func normalizePatterns(patterns []string) []string {
	if len(patterns) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		trimmed := strings.TrimSpace(pattern)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}

	return normalized
}
