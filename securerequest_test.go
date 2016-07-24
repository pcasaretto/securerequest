package securerequest_test

import "testing"

func assertPanic(t *testing.T, f func()) {
	if !didPanic(f, nil) {
		t.Error("Expected function to panic, it did not")
	}
}

func assertNoPanic(t *testing.T, f func()) {
	var m interface{}
	if didPanic(f, m) {
		t.Error("Expected function not to panic, it did panic with arg: %s", m)
	}
}

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f func(), m interface{}) bool {
	didPanic := false
	func() {
		defer func() {
			if m = recover(); m != nil {
				didPanic = true
			}
		}()
		// call the target function
		f()
	}()
	return didPanic
}
