package drift

import (
	"testing"
)

func TestPolicy_Evaluate_DefaultAlert(t *testing.T) {
	p := NewPolicy(nil)
	if got := p.Evaluate("svc", "image"); got != PolicyActionAlert {
		t.Fatalf("expected alert, got %q", got)
	}
}

func TestPolicy_Evaluate_ExactMatch(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Service: "api", Field: "image", Action: PolicyActionIgnore},
	})
	if got := p.Evaluate("api", "image"); got != PolicyActionIgnore {
		t.Fatalf("expected ignore, got %q", got)
	}
}

func TestPolicy_Evaluate_Wildcard(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Service: "*", Field: "replicas", Action: PolicyActionFail},
	})
	if got := p.Evaluate("any-service", "replicas"); got != PolicyActionFail {
		t.Fatalf("expected fail, got %q", got)
	}
}

func TestPolicy_Evaluate_CaseInsensitive(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Service: "API", Field: "IMAGE", Action: PolicyActionIgnore},
	})
	if got := p.Evaluate("api", "image"); got != PolicyActionIgnore {
		t.Fatalf("expected ignore, got %q", got)
	}
}

func TestPolicy_Validate_Valid(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Service: "svc", Field: "image", Action: PolicyActionAlert},
	})
	if err := p.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPolicy_Validate_BadAction(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Service: "svc", Field: "image", Action: "boom"},
	})
	if err := p.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestPolicy_Validate_EmptyService(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Service: "", Field: "image", Action: PolicyActionAlert},
	})
	if err := p.Validate(); err == nil {
		t.Fatal("expected validation error for empty service")
	}
}

func TestApplyPolicy_IgnoreDropsResult(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Service: "api", Field: "image", Action: PolicyActionIgnore},
	})
	results := []DriftResult{
		{Service: "api", Field: "image", Drifted: true},
	}
	out := ApplyPolicy(results, p)
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}

func TestApplyPolicy_FailAnnotatesField(t *testing.T) {
	p := NewPolicy([]PolicyRule{
		{Service: "*", Field: "image", Action: PolicyActionFail},
	})
	results := []DriftResult{
		{Service: "api", Field: "image", Drifted: true},
	}
	out := ApplyPolicy(results, p)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Field != "[FAIL] image" {
		t.Fatalf("unexpected field: %q", out[0].Field)
	}
}

func TestApplyPolicy_NilPolicyPassthrough(t *testing.T) {
	results := []DriftResult{
		{Service: "api", Field: "image", Drifted: true},
	}
	out := ApplyPolicy(results, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
}
