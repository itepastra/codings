package testinghelpers

import "testing"

func ExpectEqual[T comparable](t *testing.T, result T, correct T, context string) {
	if result != correct {
		t.Logf("result %v is not correct (%v), context: %s", result, correct, context)
		t.Fail()
	}
}
