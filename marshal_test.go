package sjon_test

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/seiyab/go-sjon"
	"github.com/stretchr/testify/require"
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
		{"abc", `"abc"`},
		{"", `""`},
		{"\n", `"\n"`},
		{true, "true"},
		{false, "false"},
		{float32(0), "0"},
		{float32(1.23), "1.23"},
		{float32(-1.23), "-1.23"},
		{float64(0), "0"},
		{float64(1.23), "1.23"},
		{float64(-1.23), "-1.23"},
	}

	s := sjon.NewSerializer()
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

func TestMarshalInvalidPrimitive(t *testing.T) {
	ts := []any{
		math.NaN(),
		math.Inf(-1),
		math.Inf(1),
	}
	s := sjon.NewSerializer()
	for _, tt := range ts {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			_, err := s.Marshal(tt)
			require.Error(t, err)
		})
	}
}

func FuzzMarshalString(f *testing.F) {
	sj := sjon.NewSerializer()
	f.Add("")
	f.Fuzz(func(t *testing.T, s string) {
		out, err := sj.Marshal(s)
		require.NoError(t, err)

		expected, err := json.Marshal(s)
		require.NoError(t, err)

		tq.Equal(t, string(expected), string(out))
		tq.Equal(t, expected, out)
	})
}

func FuzzMarshalInt(f *testing.F) {
	sj := sjon.NewSerializer()
	f.Add(0)
	f.Fuzz(func(t *testing.T, i int) {
		out, err := sj.Marshal(i)
		require.NoError(t, err)

		expected, err := json.Marshal(i)
		require.NoError(t, err)

		tq.Equal(t, string(expected), string(out))
		tq.Equal(t, expected, out)
	})
}

func FuzzMarshalUint(f *testing.F) {
	sj := sjon.NewSerializer()
	f.Add(uint(0))
	f.Fuzz(func(t *testing.T, i uint) {
		out, err := sj.Marshal(i)
		require.NoError(t, err)

		expected, err := json.Marshal(i)
		require.NoError(t, err)

		tq.Equal(t, string(expected), string(out))
		tq.Equal(t, expected, out)
	})
}

func FuzzMarshalFloat32(f *testing.F) {
	sj := sjon.NewSerializer()
	f.Add(float32(0))
	f.Fuzz(func(t *testing.T, i float32) {
		out, err := sj.Marshal(i)
		require.NoError(t, err)

		expected, err := json.Marshal(i)
		require.NoError(t, err)

		tq.Equal(t, string(expected), string(out))
		tq.Equal(t, expected, out)
	})
}

func FuzzMarshalFloat64(f *testing.F) {
	sj := sjon.NewSerializer()
	f.Add(float64(0))
	f.Fuzz(func(t *testing.T, i float64) {
		out, err := sj.Marshal(i)
		require.NoError(t, err)

		expected, err := json.Marshal(i)
		require.NoError(t, err)

		tq.Equal(t, string(expected), string(out))
		tq.Equal(t, expected, out)
	})
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
		{[]*int{ref(1), ref(2), ref(3), nil, nil}, "[1,2,3,null,null]"},
	}

	s := sjon.NewSerializer()
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

func TestMarshalStruct(t *testing.T) {
	type TestCase struct {
		value    any
		expected string
	}
	tests := []TestCase{
		{struct{}{}, "{}"},
		{struct{ A int }{1}, `{"A":1}`},
		{struct {
			A int
			B bool
		}{1, true}, `{"A":1,"B":true}`},
		{struct {
			A int
			B bool
			C *int
		}{1, true, ref(100)}, `{"A":1,"B":true,"C":100}`},
	}

	s := sjon.NewSerializer()
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

func ref[T any](v T) *T {
	return &v
}
