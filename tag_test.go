package sjon_test

import (
	"testing"

	"github.com/seiyab/go-sjon"
	"github.com/stretchr/testify/require"
)

func TestSerializerMarshal_StructTag(t *testing.T) {
	sj := sjon.NewSerializer()

	t.Run("key name by struct tag", func(t *testing.T) {
		type HyphenTest struct {
			Foo int
			Bar int `json:"-"`
		}
		actual, err := sj.Marshal(HyphenTest{1, 2})
		require.NoError(t, err)
		tq.Equal(t, `{"Foo":1}`, string(actual))

		type HyphenKeyTest struct {
			Foo int
			Bar int `json:"-,"`
		}
		actual, err = sj.Marshal(HyphenKeyTest{1, 2})
		require.NoError(t, err)
		tq.Equal(t, `{"Foo":1,"-":2}`, string(actual))

		type NameTest struct {
			Foo int
			Bar int `json:"abc"`
		}
		actual, err = sj.Marshal(NameTest{1, 2})
		require.NoError(t, err)
		tq.Equal(t, `{"Foo":1,"abc":2}`, string(actual))
	})

	t.Run("omit empty", func(t *testing.T) {
		type OmitEmptyTest struct {
			Foo int
			Bar int `json:",omitempty"`
			Baz int `json:"abc,omitempty"`
		}
		actual, err := sj.Marshal(OmitEmptyTest{})
		require.NoError(t, err)
		tq.Equal(t, `{"Foo":0}`, string(actual))

		type OmitEmptyForVariousTypesTest struct {
			A int            `json:",omitempty"`
			B bool           `json:",omitempty"`
			C string         `json:",omitempty"`
			D []int          `json:",omitempty"`
			E map[string]int `json:",omitempty"`
			F *int           `json:",omitempty"`
			G any            `json:",omitempty"`
			H *any           `json:",omitempty"`
			I OmitEmptyTest  `json:",omitempty"`
			J *OmitEmptyTest `json:",omitempty"`
		}

		actual, err = sj.Marshal(OmitEmptyForVariousTypesTest{})
		require.NoError(t, err)
		tq.Equal(t, `{"I":{"Foo":0}}`, string(actual))

		actual, err = sj.Marshal(&OmitEmptyForVariousTypesTest{
			A: 1,
			B: true,
			C: "abc",
			D: []int{1, 2},
			E: map[string]int{"a": 1},
			F: new(int),
			G: 1,
			H: new(any),
			I: OmitEmptyTest{1, 2, 3},
			J: &OmitEmptyTest{1, 2, 3},
		})
		require.NoError(t, err)
		tq.Equal(t, `{"A":1,"B":true,"C":"abc","D":[1,2],"E":{"a":1},"F":0,"G":1,"H":null,"I":{"Foo":1,"Bar":2,"abc":3},"J":{"Foo":1,"Bar":2,"abc":3}}`, string(actual))
	})
}
