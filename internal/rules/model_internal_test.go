package rules

import "testing"

func TestInterpolateMessageInternalCases(t *testing.T) {
	t.Parallel()

	for _, tc := range interpolateMessageCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := InterpolateMessage(tc.message, tc.captures)
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

type interpolateCase struct {
	name     string
	message  string
	captures []string
	want     string
}

func interpolateMessageCases() []interpolateCase {
	return []interpolateCase{
		{name: "empty message", message: "", want: ""},
		{name: "literal text", message: "plain", captures: []string{"m0"}, want: "plain"},
		{name: "replace full capture", message: "$0", captures: []string{"full"}, want: "full"},
		{
			name:     "replace numbered capture",
			message:  "value:$1",
			captures: []string{"full", "one"},
			want:     "value:one",
		},
		{name: "missing capture is empty", message: "value:$9", captures: []string{"full"}, want: "value:"},
		{name: "escaped dollar", message: "$$5", captures: []string{"full"}, want: "$5"},
		{name: "non digit after dollar", message: "$abc", captures: []string{"full"}, want: "$abc"},
		{name: "trailing dollar", message: "value$", captures: []string{"full"}, want: "value$"},
		{
			name:     "multi digit capture",
			message:  "$10",
			captures: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
			want:     "10",
		},
	}
}

func TestParseDigitsInternalCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		value      string
		wantIndex  int
		wantDigits int
	}{
		{name: "empty", value: "", wantIndex: 0, wantDigits: 0},
		{name: "single digit", value: "7", wantIndex: 7, wantDigits: 1},
		{name: "multiple digits", value: "123x", wantIndex: 123, wantDigits: 3},
		{name: "non digit prefix", value: "x123", wantIndex: 0, wantDigits: 0},
		{name: "leading zero", value: "01", wantIndex: 1, wantDigits: 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			index, digits := parseDigits(tc.value)
			if index != tc.wantIndex || digits != tc.wantDigits {
				t.Fatalf(
					"expected (index=%d, digits=%d), got (index=%d, digits=%d)",
					tc.wantIndex,
					tc.wantDigits,
					index,
					digits,
				)
			}
		})
	}
}
