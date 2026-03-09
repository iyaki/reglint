//nolint:testpackage
package output

import (
	"bytes"
	"testing"
)

func assertNoANSIControlSequences(t *testing.T, data []byte) {
	t.Helper()

	if bytes.Contains(data, []byte("\x1b[")) {
		t.Fatalf("expected output without ANSI control sequences, got:\n%s", string(data))
	}
}
