package scan

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

const (
	binaryProbeSize = 8000
	nullByte        = 0x00
)

func collectFiles(roots []string, include []string, exclude []string, maxFileSizeBytes int64) ([]string, int, error) {
	var files []string
	skipped := 0

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

			selected, fileSkipped, err := evaluateFile(path, relPath, entry, include, exclude, maxFileSizeBytes)
			if err != nil {
				return err
			}
			if fileSkipped {
				skipped++

				return nil
			}
			if selected {
				files = append(files, relPath)
			}

			return nil
		}); err != nil {
			return nil, 0, err
		}
	}

	sort.Strings(files)

	return files, skipped, nil
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

func evaluateFile(
	path string,
	relPath string,
	entry os.DirEntry,
	include []string,
	exclude []string,
	maxFileSizeBytes int64,
) (bool, bool, error) {
	match, err := matchesPath(relPath, include, exclude)
	if err != nil {
		return false, false, err
	}
	if !match {
		return false, false, nil
	}

	skip, err := shouldSkipFile(path, entry, maxFileSizeBytes)
	if err != nil {
		return false, false, err
	}
	if skip {
		return false, true, nil
	}

	return true, false, nil
}

func shouldSkipFile(path string, entry os.DirEntry, maxFileSizeBytes int64) (bool, error) {
	info, err := entry.Info()
	if err != nil {
		return false, err
	}
	if maxFileSizeBytes > 0 && info.Size() > maxFileSizeBytes {
		return true, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return false, err
	}

	buffer := make([]byte, binaryProbeSize)
	read, readErr := file.Read(buffer)
	closeErr := file.Close()
	if closeErr != nil {
		return false, closeErr
	}
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return false, readErr
	}
	for _, value := range buffer[:read] {
		if value == nullByte {
			return true, nil
		}
	}

	return false, nil
}
