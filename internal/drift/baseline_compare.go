package drift

import "fmt"

// BaselineDiff describes the drift between a live config and a baseline entry.
type BaselineDiff struct {
	Service  string
	Field    string
	Baseline string
	Live     string
}

// CompareToBaseline detects drift between live configs and a saved Baseline.
// It returns one BaselineDiff per field that has changed.
func CompareToBaseline(baseline *Baseline, live []ServiceConfig) []BaselineDiff {
	index := make(map[string]ServiceConfig, len(baseline.Configs))
	for _, c := range baseline.Configs {
		index[c.Name] = c
	}

	var diffs []BaselineDiff
	for _, lc := range live {
		bc, ok := index[lc.Name]
		if !ok {
			diffs = append(diffs, BaselineDiff{
				Service:  lc.Name,
				Field:    "existence",
				Baseline: "<missing>",
				Live:     "<present>",
			})
			continue
		}
		if bc.Image != lc.Image {
			diffs = append(diffs, BaselineDiff{
				Service:  lc.Name,
				Field:    "image",
				Baseline: bc.Image,
				Live:     lc.Image,
			})
		}
		for k, bv := range bc.Env {
			lv, exists := lc.Env[k]
			if !exists {
				diffs = append(diffs, BaselineDiff{
					Service:  lc.Name,
					Field:    fmt.Sprintf("env.%s", k),
					Baseline: bv,
					Live:     "<missing>",
				})
			} else if bv != lv {
				diffs = append(diffs, BaselineDiff{
					Service:  lc.Name,
					Field:    fmt.Sprintf("env.%s", k),
					Baseline: bv,
					Live:     lv,
				})
			}
		}
	}
	return diffs
}
