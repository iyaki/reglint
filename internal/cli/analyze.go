package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	defaultConfigPath   = "regex-rules.yaml"
	defaultFormat       = "console"
	defaultMaxFileBytes = int64(5242880)
)

// Config holds parsed analyze command inputs.
type Config struct {
	ConfigPath       string
	Roots            []string
	Formats          []string
	OutJSON          string
	OutSARIF         string
	Include          []string
	Exclude          []string
	Concurrency      int
	MaxFileSizeBytes int64
	FailOnSeverity   string
}

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)

	return nil
}

// ParseAnalyzeArgs parses analyze command arguments into a Config.
func ParseAnalyzeArgs(args []string) (Config, error) {
	var cfg Config

	flagSet := flag.NewFlagSet("analyze", flag.ContinueOnError)
	flagSet.SetOutput(&strings.Builder{})

	flagSet.StringVar(&cfg.ConfigPath, "config", defaultConfigPath, "Path to YAML rules config file.")
	formatValue := flagSet.String("format", defaultFormat, "Comma-separated list of formats.")
	flagSet.StringVar(&cfg.OutJSON, "out-json", "", "Output path for JSON results.")
	flagSet.StringVar(&cfg.OutSARIF, "out-sarif", "", "Output path for SARIF results.")
	var include stringSlice
	var exclude stringSlice
	flagSet.Var(&include, "include", "Repeatable include glob.")
	flagSet.Var(&exclude, "exclude", "Repeatable exclude glob.")
	flagSet.IntVar(&cfg.Concurrency, "concurrency", runtime.GOMAXPROCS(0), "Worker count.")
	flagSet.Int64Var(&cfg.MaxFileSizeBytes, "max-file-size", defaultMaxFileBytes, "Maximum file size in bytes.")
	flagSet.StringVar(&cfg.FailOnSeverity, "fail-on", "", "Fail if matches at or above severity.")

	if err := flagSet.Parse(args); err != nil {
		return Config{}, err
	}

	cfg.Include = include
	cfg.Exclude = exclude

	if flagSet.NArg() == 0 {
		cfg.Roots = []string{"."}
	} else {
		cfg.Roots = append([]string{}, flagSet.Args()...)
	}

	formats, err := parseFormats(*formatValue)
	if err != nil {
		return Config{}, err
	}
	cfg.Formats = formats

	if err := validateAnalyzeConfig(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func parseFormats(value string) ([]string, error) {
	if strings.TrimSpace(value) == "" {
		return nil, errors.New("format must not be empty")
	}

	parts := strings.Split(value, ",")
	seen := make(map[string]struct{}, len(parts))
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		format := strings.TrimSpace(part)
		if format == "" {
			return nil, errors.New("format must not be empty")
		}
		if !isValidFormat(format) {
			return nil, fmt.Errorf("invalid format: %s", format)
		}
		if _, ok := seen[format]; ok {
			continue
		}
		seen[format] = struct{}{}
		result = append(result, format)
	}

	return result, nil
}

func isValidFormat(value string) bool {
	switch value {
	case "console", "json", "sarif":
		return true
	default:
		return false
	}
}

func validateAnalyzeConfig(cfg Config) error {
	if err := validateConfigPath(cfg.ConfigPath); err != nil {
		return err
	}
	if cfg.Concurrency <= 0 {
		return errors.New("concurrency must be positive")
	}
	if cfg.MaxFileSizeBytes <= 0 {
		return errors.New("max-file-size must be positive")
	}
	if cfg.FailOnSeverity != "" && !isValidSeverity(cfg.FailOnSeverity) {
		return fmt.Errorf("invalid fail-on value: %s", cfg.FailOnSeverity)
	}
	if err := validateOutputPaths(cfg); err != nil {
		return err
	}

	return nil
}

func validateConfigPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("config file not found: %s", path)
	}
	if info.IsDir() {
		return fmt.Errorf("config path is a directory: %s", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config file not readable: %s", path)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("config file not readable: %s", path)
	}

	return nil
}

func validateOutputPaths(cfg Config) error {
	if len(cfg.Formats) <= 1 {
		return nil
	}

	for _, format := range cfg.Formats {
		switch format {
		case "json":
			if cfg.OutJSON == "" {
				return errors.New("--out-json is required when requesting json with multiple formats")
			}
		case "sarif":
			if cfg.OutSARIF == "" {
				return errors.New("--out-sarif is required when requesting sarif with multiple formats")
			}
		}
	}

	return nil
}

func isValidSeverity(value string) bool {
	switch value {
	case "error", "warning", "notice", "info":
		return true
	default:
		return false
	}
}
