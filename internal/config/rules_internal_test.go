package config

import "testing"

func TestRuleToRulesRulePreservesExplicitValues(t *testing.T) {
	t.Parallel()

	rule := Rule{
		Message:  "message",
		Regex:    "value",
		Severity: "info",
		Paths:    []string{"src/**"},
		Exclude:  []string{"src/vendor/**"},
	}

	converted := rule.toRulesRule()
	if converted.Severity != "info" {
		t.Fatalf("expected severity info, got %q", converted.Severity)
	}
	if len(converted.Paths) != 1 || converted.Paths[0] != "src/**" {
		t.Fatalf("expected paths [src/**], got %v", converted.Paths)
	}
	if len(converted.Exclude) != 1 || converted.Exclude[0] != "src/vendor/**" {
		t.Fatalf("expected exclude [src/vendor/**], got %v", converted.Exclude)
	}

	converted.Paths[0] = "changed/**"
	converted.Exclude[0] = "changed/**"
	if rule.Paths[0] != "src/**" {
		t.Fatalf("expected original paths to remain unchanged, got %v", rule.Paths)
	}
	if rule.Exclude[0] != "src/vendor/**" {
		t.Fatalf("expected original exclude to remain unchanged, got %v", rule.Exclude)
	}
}

func TestRuleSetToRulesLeavesRuleSpecificPathsAndExclude(t *testing.T) {
	t.Parallel()

	ruleSet := RuleSet{
		Include: []string{"global/**"},
		Exclude: []string{"global-exclude/**"},
		Rules: []Rule{
			{
				Message: "rule",
				Regex:   "value",
				Paths:   []string{"rule/**"},
				Exclude: []string{"rule-exclude/**"},
			},
		},
	}

	converted := ruleSet.ToRules()
	if len(converted.Rules) != 1 {
		t.Fatalf("expected 1 converted rule, got %d", len(converted.Rules))
	}
	if len(converted.Rules[0].Paths) != 1 || converted.Rules[0].Paths[0] != "rule/**" {
		t.Fatalf("expected rule-specific paths, got %v", converted.Rules[0].Paths)
	}
	if len(converted.Rules[0].Exclude) != 1 || converted.Rules[0].Exclude[0] != "rule-exclude/**" {
		t.Fatalf("expected rule-specific exclude, got %v", converted.Rules[0].Exclude)
	}
}

func TestRuleSetToRulesHandlesEmptyRules(t *testing.T) {
	t.Parallel()

	ruleSet := RuleSet{}
	converted := ruleSet.ToRules()

	if len(converted.Rules) != 0 {
		t.Fatalf("expected zero converted rules, got %d", len(converted.Rules))
	}
	if len(converted.Include) != 1 || converted.Include[0] != "**/*" {
		t.Fatalf("expected default include, got %v", converted.Include)
	}
	if len(converted.Exclude) != 3 {
		t.Fatalf("expected default exclude list, got %v", converted.Exclude)
	}
	if converted.Concurrency == nil {
		t.Fatal("expected default concurrency to be set")
	}
}
