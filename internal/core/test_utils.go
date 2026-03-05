package core

import (
	"testing"
)

// AssertEqual is a simple helper for testing equality.
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
