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
		{int8(math.MinInt8), "-128"},
		{int8(math.MaxInt8), "127"},
		{int16(math.MinInt16), "-32768"},
		{int16(math.MaxInt16), "32767"},
		{int32(math.MinInt32), "-2147483648"},
		{int32(math.MaxInt32), "2147483647"},
		{int64(math.MinInt64), "-9223372036854775808"},
		{int64(math.MaxInt64), "9223372036854775807"},
		{uint8(math.MaxUint8), "255"},
		{uint16(math.MaxUint16), "65535"},
		{uint32(math.MaxUint32), "4294967295"},
		{uint64(math.MaxUint64), "18446744073709551615"},
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
		{[]int{}, "[]"},
		{[]int{1, 2, 3}, "[1,2,3]"},
		{[]bool{false, true, true, false}, `[false,true,true,false]`},
		{[0]int{}, "[]"},
		{[2]int{1, 2}, "[1,2]"},
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
