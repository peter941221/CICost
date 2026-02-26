package policy

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Severity string

const (
	SeverityError Severity = "error"
	SeverityWarn  Severity = "warn"
	SeverityInfo  Severity = "info"
)

type Config struct {
	Rules   []Rule  `yaml:"rules"`
	Actions Actions `yaml:"actions"`
}

type Rule struct {
	ID       string   `yaml:"id"`
	When     string   `yaml:"when"`
	Severity Severity `yaml:"severity"`
}

type Actions struct {
	OnError string `yaml:"on_error"`
	OnWarn  string `yaml:"on_warn"`
}

type Finding struct {
	RuleID        string   `json:"rule_id"`
	Severity      Severity `json:"severity"`
	When          string   `json:"when"`
	EvidenceKey   string   `json:"evidence_key"`
	EvidenceValue float64  `json:"evidence_value"`
}

type expression struct {
	variable string
	operator string
	value    float64
}

var exprRe = regexp.MustCompile(`^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*(<=|>=|==|!=|>|<)\s*([0-9]+(?:\.[0-9]+)?)\s*$`)

var allowedMetrics = map[string]struct{}{
	"monthly_cost_usd": {},
	"waste_percentage": {},
	"fail_rate":        {},
	"total_runs":       {},
	"total_cost_usd":   {},
}

func LoadFromFile(path string) (Config, error) {
	var cfg Config
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("invalid policy yaml: %w", err)
	}
	return cfg, nil
}

func Lint(cfg Config) error {
	if len(cfg.Rules) == 0 {
		return fmt.Errorf("policy has no rules")
	}
	for i, rule := range cfg.Rules {
		if strings.TrimSpace(rule.ID) == "" {
			return fmt.Errorf("rule[%d] missing id", i)
		}
		if strings.TrimSpace(rule.When) == "" {
			return fmt.Errorf("rule[%s] missing when", rule.ID)
		}
		if !isValidSeverity(rule.Severity) {
			return fmt.Errorf("rule[%s] invalid severity %q", rule.ID, rule.Severity)
		}
		expr, err := parseExpression(rule.When)
		if err != nil {
			return fmt.Errorf("rule[%s] invalid expression: %w", rule.ID, err)
		}
		if _, ok := allowedMetrics[expr.variable]; !ok {
			return fmt.Errorf("rule[%s] uses unsupported variable %q", rule.ID, expr.variable)
		}
	}
	return nil
}

func Evaluate(cfg Config, metrics map[string]float64) ([]Finding, error) {
	if err := Lint(cfg); err != nil {
		return nil, err
	}
	out := make([]Finding, 0, len(cfg.Rules))
	for _, rule := range cfg.Rules {
		expr, _ := parseExpression(rule.When)
		v, ok := metrics[expr.variable]
		if !ok {
			return nil, fmt.Errorf("missing metric %q for rule %s", expr.variable, rule.ID)
		}
		if compare(v, expr.operator, expr.value) {
			out = append(out, Finding{
				RuleID:        rule.ID,
				Severity:      rule.Severity,
				When:          rule.When,
				EvidenceKey:   expr.variable,
				EvidenceValue: v,
			})
		}
	}
	return out, nil
}

func parseExpression(input string) (expression, error) {
	m := exprRe.FindStringSubmatch(strings.TrimSpace(input))
	if len(m) != 4 {
		return expression{}, fmt.Errorf("expected format `<metric> <op> <number>`, got %q", input)
	}
	value, err := strconv.ParseFloat(m[3], 64)
	if err != nil {
		return expression{}, err
	}
	return expression{
		variable: m[1],
		operator: m[2],
		value:    value,
	}, nil
}

func compare(left float64, op string, right float64) bool {
	switch op {
	case ">":
		return left > right
	case ">=":
		return left >= right
	case "<":
		return left < right
	case "<=":
		return left <= right
	case "==":
		return left == right
	case "!=":
		return left != right
	default:
		return false
	}
}

func isValidSeverity(s Severity) bool {
	switch s {
	case SeverityError, SeverityWarn, SeverityInfo:
		return true
	default:
		return false
	}
}
