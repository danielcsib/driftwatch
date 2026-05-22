package drift

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a key-value label attached to a service config.
type Tag struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// TagSet is an ordered collection of Tags for a service.
type TagSet struct {
	Service string `json:"service" yaml:"service"`
	Tags    []Tag  `json:"tags" yaml:"tags"`
}

// TagIndex maps service names to their TagSets for fast lookup.
type TagIndex map[string]TagSet

// NewTagIndex builds a TagIndex from a slice of TagSets.
func NewTagIndex(sets []TagSet) TagIndex {
	idx := make(TagIndex, len(sets))
	for _, ts := range sets {
		idx[strings.ToLower(ts.Service)] = ts
	}
	return idx
}

// Get returns the TagSet for the given service name (case-insensitive).
// The second return value reports whether the service was found.
func (idx TagIndex) Get(service string) (TagSet, bool) {
	ts, ok := idx[strings.ToLower(service)]
	return ts, ok
}

// HasTag reports whether the TagSet for service contains a tag with the
// given key and value (both case-insensitive).
func (idx TagIndex) HasTag(service, key, value string) bool {
	ts, ok := idx.Get(service)
	if !ok {
		return false
	}
	for _, t := range ts.Tags {
		if strings.EqualFold(t.Key, key) && strings.EqualFold(t.Value, value) {
			return true
		}
	}
	return false
}

// FilterByTag returns only those DriftResults whose service has the given
// tag key/value pair in the provided TagIndex.
func FilterByTag(results []DriftResult, idx TagIndex, key, value string) []DriftResult {
	out := make([]DriftResult, 0, len(results))
	for _, r := range results {
		if idx.HasTag(r.Service, key, value) {
			out = append(out, r)
		}
	}
	return out
}

// String returns a stable human-readable representation of a TagSet.
func (ts TagSet) String() string {
	pairs := make([]string, 0, len(ts.Tags))
	for _, t := range ts.Tags {
		pairs = append(pairs, fmt.Sprintf("%s=%s", t.Key, t.Value))
	}
	sort.Strings(pairs)
	return fmt.Sprintf("%s[%s]", ts.Service, strings.Join(pairs, ","))
}
