package testhelpers

import (
	"testing"

	"github.com/splitio/go-toolkit/datastructures/set"
)

// AssertStringSliceEquals fails is two string slices are not identical
func AssertStringSliceEquals(t *testing.T, actual []string, expected []string, message string) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Errorf(message)
		t.Errorf("Slices have different sizes. Actual: %d, expected: %d", len(actual), len(expected))
		t.Errorf("Actual: %v || Expected: %v", actual, expected)
		return
	}

	idx := 0
	for idx < len(actual) && actual[idx] == expected[idx] {
		idx++
	}

	if idx != len(actual) {
		t.Errorf(message)
		t.Errorf("Slices have different elements")
		t.Errorf("Actual: %v || Expected: %v", actual, expected)
	}
}

func AssertStringSliceEqualsNoOrder(t *testing.T, actual []string, expected []string, message string) {
	t.Helper()
	asInterfaces1 := make([]interface{}, 0, len(actual))
	for _, s := range actual {
		asInterfaces1 = append(asInterfaces1, s)
	}
	asInterfaces2 := make([]interface{}, 0, len(expected))
	for _, s := range expected {
		asInterfaces2 = append(asInterfaces2, s)
	}
	set1 := set.NewSet(asInterfaces1...)
	set2 := set.NewSet(asInterfaces2...)
	if !set1.IsEqual(set2) {
		t.Error("slices contain different elements despite order: ", message)
		t.Error("actual: ", actual)
		t.Error("expected: ", expected)
	}
}
