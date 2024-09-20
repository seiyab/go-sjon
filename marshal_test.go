package sjon_test

import (
	"testing"

	"github.com/seiyab/go-sjon"
)

func TestMarshalPrimitive(t *testing.T) {
	type TestCase struct {
		value    any
		expected string
	}
	tests := []TestCase{
		{nil, "null"},
		{true, "true"},
		{false, "false"},
	}

	s := sjon.New()
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			actual, err := s.Marshal(tt.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(actual) != tt.expected {
				t.Errorf("expected %q, but got %q", tt.expected, string(actual))
			}
		})
	}

}
