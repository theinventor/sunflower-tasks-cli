package cmd

import "testing"

func TestEndpointCatalogIsComplete(t *testing.T) {
	if len(endpoints) < 100 {
		t.Fatalf("catalog has %d operations", len(endpoints))
	}
	seen := map[string]bool{}
	for _, e := range endpoints {
		k := e.Method + " " + e.Path
		if seen[k] {
			t.Fatalf("duplicate endpoint %s", k)
		}
		seen[k] = true
	}
}
