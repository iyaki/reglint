package output

import (
	"fmt"
	"strings"
)

// Registry holds available formatters keyed by format name.
type Registry struct {
	formats map[string]Formatter
}

// NewRegistry builds a registry from the provided formatters.
func NewRegistry(formatters ...Formatter) (*Registry, error) {
	formats := make(map[string]Formatter, len(formatters))
	for _, formatter := range formatters {
		if formatter == nil {
			return nil, fmt.Errorf("formatter must not be nil")
		}
		name := strings.TrimSpace(formatter.Name())
		if name == "" {
			return nil, fmt.Errorf("formatter name must not be empty")
		}
		if name != strings.ToLower(name) {
			return nil, fmt.Errorf("formatter name must be lowercase: %s", name)
		}
		if _, exists := formats[name]; exists {
			return nil, fmt.Errorf("duplicate formatter: %s", name)
		}
		formats[name] = formatter
	}

	return &Registry{formats: formats}, nil
}

// Resolve returns formatters for the requested names in order.
func (r *Registry) Resolve(names []string) ([]Formatter, error) {
	resolved := make([]Formatter, 0, len(names))
	for _, name := range names {
		formatter, ok := r.formats[name]
		if !ok {
			return nil, fmt.Errorf("invalid format: %s", name)
		}
		resolved = append(resolved, formatter)
	}

	return resolved, nil
}

// ResolveName looks up a formatter by name.
func (r *Registry) ResolveName(name string) (Formatter, error) {
	formatter, ok := r.formats[name]
	if !ok {
		return nil, fmt.Errorf("invalid format: %s", name)
	}

	return formatter, nil
}
