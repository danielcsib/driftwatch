package drift

import (
	"testing"
)

func makeTags() []TagSet {
	return []TagSet{
		{
			Service: "auth-service",
			Tags:    []Tag{{Key: "env", Value: "prod"}, {Key: "team", Value: "platform"}},
		},
		{
			Service: "billing-service",
			Tags:    []Tag{{Key: "env", Value: "staging"}, {Key: "team", Value: "finance"}},
		},
		{
			Service: "inventory-service",
			Tags:    []Tag{{Key: "env", Value: "prod"}, {Key: "team", Value: "ops"}},
		},
	}
}

func TestNewTagIndex_Get(t *testing.T) {
	idx := NewTagIndex(makeTags())

	ts, ok := idx.Get("auth-service")
	if !ok {
		t.Fatal("expected auth-service to be found")
	}
	if ts.Service != "auth-service" {
		t.Errorf("unexpected service name: %s", ts.Service)
	}
}

func TestTagIndex_Get_CaseInsensitive(t *testing.T) {
	idx := NewTagIndex(makeTags())

	_, ok := idx.Get("AUTH-SERVICE")
	if !ok {
		t.Fatal("expected case-insensitive lookup to succeed")
	}
}

func TestTagIndex_Get_Missing(t *testing.T) {
	idx := NewTagIndex(makeTags())

	_, ok := idx.Get("unknown-service")
	if ok {
		t.Fatal("expected missing service to return false")
	}
}

func TestTagIndex_HasTag_Found(t *testing.T) {
	idx := NewTagIndex(makeTags())

	if !idx.HasTag("auth-service", "env", "prod") {
		t.Error("expected HasTag to return true for env=prod on auth-service")
	}
}

func TestTagIndex_HasTag_CaseInsensitive(t *testing.T) {
	idx := NewTagIndex(makeTags())

	if !idx.HasTag("Auth-Service", "ENV", "PROD") {
		t.Error("expected case-insensitive HasTag to return true")
	}
}

func TestTagIndex_HasTag_NotFound(t *testing.T) {
	idx := NewTagIndex(makeTags())

	if idx.HasTag("billing-service", "env", "prod") {
		t.Error("expected HasTag to return false for env=prod on billing-service")
	}
}

func TestFilterByTag(t *testing.T) {
	idx := NewTagIndex(makeTags())

	results := []DriftResult{
		{Service: "auth-service", Drifted: true},
		{Service: "billing-service", Drifted: false},
		{Service: "inventory-service", Drifted: true},
	}

	prod := FilterByTag(results, idx, "env", "prod")
	if len(prod) != 2 {
		t.Fatalf("expected 2 prod results, got %d", len(prod))
	}
	for _, r := range prod {
		if r.Service == "billing-service" {
			t.Error("billing-service should not appear in prod results")
		}
	}
}

func TestTagSet_String(t *testing.T) {
	ts := TagSet{
		Service: "auth-service",
		Tags:    []Tag{{Key: "team", Value: "platform"}, {Key: "env", Value: "prod"}},
	}
	got := ts.String()
	want := "auth-service[env=prod,team=platform]"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
