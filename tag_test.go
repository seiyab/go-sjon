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
	})
}
