package sjon_test

import (
	"math"
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
		{0, "0"},
		{1, "1"},
		{100, "100"},
		{1234567890, "1234567890"},
		{math.MinInt8, "-128"},
		{math.MaxInt8, "127"},
		{math.MinInt16, "-32768"},
		{math.MaxInt16, "32767"},
		{math.MinInt32, "-2147483648"},
		{math.MaxInt32, "2147483647"},
		{math.MinInt64, "-9223372036854775808"},
		{math.MaxInt64, "9223372036854775807"},
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

func TestMarshalArray(t *testing.T) {
	type TestCase struct {
		value    any
		expected string
	}
	tests := []TestCase{
		{[]int{1, 2, 3}, "[1,2,3]"},
		{[]bool{false, true, true, false}, `[false,true,true,false]`},
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
