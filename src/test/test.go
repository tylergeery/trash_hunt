package test

import (
	"testing"
)

// FatalOnError to be reused by testing package
func FatalOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected err: %s", err.Error())
	}
}

// ExpectEqualString expects  values to be equal
func ExpectEqualString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Expected equal string values: (%s) (%s)", expected, actual)
	}
}

// ExpectNotEmptyString expects string to not be empty
func ExpectNotEmptyString(t *testing.T, actual string) {
	if actual == "" {
		t.Fatalf("Expected not empty string value")
	}
}

// ExpectEqualInt64s expects int64 values to be equal
func ExpectEqualInt64s(t *testing.T, expected, actual int64) {
	if expected != actual {
		t.Fatalf("Expected equal int64 values: (%d) (%d)", expected, actual)
	}
}
