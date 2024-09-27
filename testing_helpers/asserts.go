package testinghelpers

import "testing"

func ExpectEqual(t *testing.T, result byte, correct byte, context string) {
	if result != correct {
		t.Logf("result %+q is not correct (%+q), context: %s", result, correct, context)
		t.Fail()
	}
}
