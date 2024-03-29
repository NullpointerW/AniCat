package errs

import "testing"

func NoError(t *testing.T, err error, message ...string) {
	t.Helper()
	if err != nil {
		t.Error(err, message)
	}
}

func HasError(t *testing.T, err error, message ...string) {
	t.Helper()
	if err == nil {
		t.Error(err, message)
	}
}
