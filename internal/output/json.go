package output

import (
	"encoding/json"
	"io"
	"sort"

	"github.com/iyaki/regex-checker/internal/scan"
)

type jsonResult struct {
	SchemaVersion int          `json:"schemaVersion"`
	Matches       []scan.Match `json:"matches"`
	Stats         scan.Stats   `json:"stats"`
}

// WriteJSON renders a scan result to the provided writer.
func WriteJSON(result scan.Result, out io.Writer) error {
	matches := append([]scan.Match{}, result.Matches...)
	sort.Slice(matches, func(i, j int) bool {
		left := matches[i]
		right := matches[j]
		if left.FilePath != right.FilePath {
			return left.FilePath < right.FilePath
		}
		if left.Line != right.Line {
			return left.Line < right.Line
		}
		if left.Column != right.Column {
			return left.Column < right.Column
		}
		if left.Severity != right.Severity {
			return severityRank(left.Severity) < severityRank(right.Severity)
		}

		return left.Message < right.Message
	})

	payload := jsonResult{
		SchemaVersion: 1,
		Matches:       matches,
		Stats:         result.Stats,
	}

	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "\t")

	return encoder.Encode(payload)
}
