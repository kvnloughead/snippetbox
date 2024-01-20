package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	// Test helpers won't be cited in test output.
	t.Helper()

	if actual != expected {
		t.Errorf("got %v; want %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstr string) {
	t.Helper()

	if !strings.Contains(actual, expectedSubstr) {
		t.Errorf("got %q; expected to contain %q", actual, expectedSubstr)
	}
}
