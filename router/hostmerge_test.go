package router

import (
	"reflect"
	"testing"
)

func TestHostMerger(t *testing.T) {
	h1 := []Host{
		{IP: "0.0.0.1", MAC: "1", Name: "name1"},
		{IP: "0.0.0.3", MAC: "3", Name: "name3"},
	}

	h2 := []Host{
		{IP: "0.0.0.1", MAC: "1", Name: "newname1"},
		{IP: "0.0.0.4", MAC: "4", Name: "name4"},
	}

	// Merge both hosts but use h2s names in case of dup
	h1h2 := []Host{
		{IP: "0.0.0.1", MAC: "1", Name: "newname1"},
		{IP: "0.0.0.3", MAC: "3", Name: "name3"},
		{IP: "0.0.0.4", MAC: "4", Name: "name4"},
	}
	if m := mergeHosts(h1, h2); !reflect.DeepEqual(m, h1h2) {
		t.Fatalf("Merge failed, got: %s", m)
	}
}
