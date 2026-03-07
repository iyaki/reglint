// Package ignore provides ignore file parsing and matching.
package ignore

import (
	"fmt"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// IgnoreRule captures a parsed ignore rule with source metadata.
type IgnoreRule struct { //nolint:revive
	BaseDir       string
	Source        string
	Line          int
	Pattern       string
	Negated       bool
	DirectoryOnly bool
}

// Parse converts ignore file content into ordered rules.
func Parse(baseDir string, source string, content string) ([]IgnoreRule, error) {
	lines := strings.Split(normalizeLineEndings(content), "\n")
	rules := make([]IgnoreRule, 0, len(lines))
	for idx, raw := range lines {
		lineNumber := idx + 1
		rule, ok, err := parseRuleLine(baseDir, source, lineNumber, raw)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func parseRuleLine(baseDir string, source string, lineNumber int, raw string) (IgnoreRule, bool, error) {
	line := strings.TrimRight(raw, "\r")
	if line == "" {
		return IgnoreRule{}, false, nil
	}
	if strings.HasPrefix(line, "#") {
		return IgnoreRule{}, false, nil
	}

	line, escaped := trimEscapedPrefix(line)
	negated := false
	if !escaped && strings.HasPrefix(line, "!") {
		negated = true
		line = strings.TrimPrefix(line, "!")
		line, _ = trimEscapedPrefix(line)
	}
	if line == "" {
		return IgnoreRule{}, false, nil
	}

	directoryOnly := strings.HasSuffix(line, "/")
	if directoryOnly {
		line = strings.TrimSuffix(line, "/")
		if line == "" {
			return IgnoreRule{}, false, nil
		}
	}

	if err := validatePattern(line, source, lineNumber); err != nil {
		return IgnoreRule{}, false, err
	}

	return IgnoreRule{
		BaseDir:       baseDir,
		Source:        source,
		Line:          lineNumber,
		Pattern:       line,
		Negated:       negated,
		DirectoryOnly: directoryOnly,
	}, true, nil
}

func normalizeLineEndings(content string) string {
	if !strings.Contains(content, "\r") {
		return content
	}
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	return content
}

func validatePattern(pattern string, source string, lineNumber int) error {
	if _, err := doublestar.Match(pattern, "probe"); err != nil {
		return fmt.Errorf("%s:%d: invalid ignore pattern", source, lineNumber)
	}

	return nil
}

func trimEscapedPrefix(line string) (string, bool) {
	if strings.HasPrefix(line, "\\#") {
		return strings.TrimPrefix(line, "\\"), true
	}
	if strings.HasPrefix(line, "\\!") {
		return strings.TrimPrefix(line, "\\"), true
	}

	return line, false
}
