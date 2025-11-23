package optional_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/timhugh/optional.go"
)

func expectNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func expectEqual[T any](t *testing.T, expected, actual T) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func expectError(t *testing.T, expected error, err error, message string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: expected error %v, got none", message, expected)
	} else if !errors.Is(err, expected) {
		t.Errorf("%s: expected error %v, got %v", message, expected, err)
	}
}

func expectTrue(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Errorf("%s: expected true, got false", message)
	}
}

func expectFalse(t *testing.T, condition bool, message string) {
	t.Helper()
	if condition {
		t.Errorf("%s: expected false, got true", message)
	}
}

func expectEmpty[T any](t *testing.T, opt optional.Optional[T]) {
	t.Helper()
	expectTrue(t, opt.Empty(), "optional should be empty")
	expectFalse(t, opt.HasValue(), "optional should not have value")
	_, err := opt.Get()
	expectError(t, optional.ErrNoValue, err, "getting value from empty optional should return error")
}

func expectHasValue[T any](t *testing.T, opt optional.Optional[T], expected T) {
	t.Helper()
	expectFalse(t, opt.Empty(), "optional should not be empty")
	expectTrue(t, opt.HasValue(), "optional should have value")
	value, err := opt.Get()
	expectNoError(t, err)
	expectEqual(t, expected, value)
}
