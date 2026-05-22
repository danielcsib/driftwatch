package drift

import (
	"testing"
)

func makeResults() []DetectResult {
	return []DetectResult{
		{
			Service: "api",
			Drifted: true,
			Diffs: []DriftDiff{
				{Field: "image", Desired: "v1", Actual: "v2"},
				{Field: "env", Desired: "DEBUG=true", Actual: "DEBUG=false"},
			},
		},
		{
			Service: "worker",
			Drifted: false,
			Diffs:   nil,
		},
		{
			Service: "gateway",
			Drifted: true,
			Diffs: []DriftDiff{
				{Field: "replicas", Desired: "3", Actual: "1"},
			},
		},
	}
}

func TestFilter_NoOptions(t *testing.T) {
	results := makeResults()
	got := Filter(results, FilterOptions{})
	if len(got) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(got))
	}
}

func TestFilter_OnlyDrifted(t *testing.T) {
	got := Filter(makeResults(), FilterOptions{OnlyDrifted: true})
	if len(got) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(got))
	}
	for _, r := range got {
		if !r.Drifted {
			t.Errorf("non-drifted result included: %s", r.Service)
		}
	}
}

func TestFilter_ByService(t *testing.T) {
	got := Filter(makeResults(), FilterOptions{Services: []string{"api", "worker"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestFilter_ByService_CaseInsensitive(t *testing.T) {
	got := Filter(makeResults(), FilterOptions{Services: []string{"API"}})
	if len(got) != 1 || got[0].Service != "api" {
		t.Fatalf("expected api result, got %+v", got)
	}
}

func TestFilter_ByField(t *testing.T) {
	got := Filter(makeResults(), FilterOptions{Fields: []string{"image"}})
	// worker has no diffs so passes through; api keeps only image diff; gateway has no image diff
	var apiResult *DetectResult
	for i := range got {
		if got[i].Service == "api" {
			apiResult = &got[i]
		}
	}
	if apiResult == nil {
		t.Fatal("api result missing after field filter")
	}
	if len(apiResult.Diffs) != 1 || apiResult.Diffs[0].Field != "image" {
		t.Errorf("unexpected diffs after field filter: %+v", apiResult.Diffs)
	}
}

func TestFilter_ByFieldAndOnlyDrifted(t *testing.T) {
	// Only image field + only drifted: gateway has no image diff so excluded
	got := Filter(makeResults(), FilterOptions{
		Fields:      []string{"image"},
		OnlyDrifted: true,
	})
	if len(got) != 1 || got[0].Service != "api" {
		t.Fatalf("expected only api, got %+v", got)
	}
}
