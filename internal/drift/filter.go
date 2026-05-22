package drift

import "strings"

// FilterOptions controls which drift results are surfaced.
type FilterOptions struct {
	// Services restricts results to the named services (case-insensitive).
	// An empty slice means all services are included.
	Services []string

	// OnlyDrifted, when true, excludes results where no drift was detected.
	OnlyDrifted bool

	// Fields restricts drift comparison to the listed field names (e.g. "image",
	// "env"). An empty slice means all fields are compared.
	Fields []string
}

// Filter applies FilterOptions to a slice of DetectResult, returning only the
// results that satisfy every active constraint.
func Filter(results []DetectResult, opts FilterOptions) []DetectResult {
	serviceSet := toLowerSet(opts.Services)
	fieldSet := toLowerSet(opts.Fields)

	out := make([]DetectResult, 0, len(results))
	for _, r := range results {
		if len(serviceSet) > 0 {
			if _, ok := serviceSet[strings.ToLower(r.Service)]; !ok {
				continue
			}
		}

		if opts.OnlyDrifted && !r.Drifted {
			continue
		}

		if len(fieldSet) > 0 {
			r = filterFields(r, fieldSet)
			// After field filtering the result may no longer be drifted.
			if opts.OnlyDrifted && !r.Drifted {
				continue
			}
		}

		out = append(out, r)
	}
	return out
}

// filterFields returns a copy of r whose Diffs list contains only entries
// whose Field name is present in the allowed set. Drifted is recalculated.
func filterFields(r DetectResult, allowed map[string]struct{}) DetectResult {
	filtered := make([]DriftDiff, 0, len(r.Diffs))
	for _, d := range r.Diffs {
		if _, ok := allowed[strings.ToLower(d.Field)]; ok {
			filtered = append(filtered, d)
		}
	}
	r.Diffs = filtered
	r.Drifted = len(filtered) > 0
	return r
}

func toLowerSet(ss []string) map[string]struct{} {
	m := make(map[string]struct{}, len(ss))
	for _, s := range ss {
		m[strings.ToLower(s)] = struct{}{}
	}
	return m
}
