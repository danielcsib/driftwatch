package drift

import (
	"fmt"
	"strings"
)

// PolicyAction defines what action to take when drift is detected.
type PolicyAction string

const (
	PolicyActionAlert  PolicyAction = "alert"
	PolicyActionIgnore PolicyAction = "ignore"
	PolicyActionFail   PolicyAction = "fail"
)

// PolicyRule defines a rule that matches services and fields.
type PolicyRule struct {
	Service string       `json:"service" yaml:"service"`
	Field   string       `json:"field"   yaml:"field"`
	Action  PolicyAction `json:"action"  yaml:"action"`
}

// Policy holds a set of rules applied during drift evaluation.
type Policy struct {
	Rules []PolicyRule `json:"rules" yaml:"rules"`
}

// NewPolicy returns a Policy with the given rules.
func NewPolicy(rules []PolicyRule) *Policy {
	return &Policy{Rules: rules}
}

// Evaluate returns the action that applies to the given service and field.
// If no rule matches, PolicyActionAlert is returned as the default.
func (p *Policy) Evaluate(service, field string) PolicyAction {
	for _, r := range p.Rules {
		if matchesGlob(r.Service, service) && matchesGlob(r.Field, field) {
			return r.Action
		}
	}
	return PolicyActionAlert
}

// Validate checks that all rules contain valid actions.
func (p *Policy) Validate() error {
	for i, r := range p.Rules {
		switch r.Action {
		case PolicyActionAlert, PolicyActionIgnore, PolicyActionFail:
			// valid
		default:
			return fmt.Errorf("rule %d: unknown action %q", i, r.Action)
		}
		if strings.TrimSpace(r.Service) == "" {
			return fmt.Errorf("rule %d: service must not be empty", i)
		}
		if strings.TrimSpace(r.Field) == "" {
			return fmt.Errorf("rule %d: field must not be empty", i)
		}
	}
	return nil
}

// matchesGlob returns true if pattern is "*" or equals value (case-insensitive).
func matchesGlob(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	return strings.EqualFold(pattern, value)
}
