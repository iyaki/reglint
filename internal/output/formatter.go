package output

import (
	"io"

	"github.com/iyaki/regex-checker/internal/scan"
)

// Formatter renders a scan result to the provided writer.
type Formatter interface {
	Name() string
	Write(result scan.Result, out io.Writer) error
}
