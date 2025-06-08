package testutil

import "testing"

type TestFunc func()

func AssertPanic(t *testing.T, fn TestFunc, args ...any) (panicValue any) {

	t.Helper()

	defer func() {
		if panicValue = recover(); panicValue == nil {
			t.Errorf("Panic expected but no panic occurred. %v", args...)
		}
	}()
	fn()

	return
}
