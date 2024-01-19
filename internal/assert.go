package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	// Test helpers won't be cited in test output.
	t.Helper()

	if actual != expected {
		t.Errorf("got %v; want %v", actual, expected)
	}
}
