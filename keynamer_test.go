package sjon_test

import (
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/seiyab/go-sjon"
	"github.com/seiyab/teq"
	"github.com/stretchr/testify/require"
)

type Fixture struct {
	Foo       int
	FooBar    bool
	FooBarBaz []int
}

var tq = teq.New()

func TestSerializerKeyNamer(t *testing.T) {
	t.Run("lowerCamel", func(t *testing.T) {
		sj := sjon.NewSerializer()
		sj.KeyNamer = strcase.ToLowerCamel

		actual, err := sj.Marshal(Fixture{
			Foo:       1,
			FooBar:    true,
			FooBarBaz: []int{1, 2, 3},
		})
		require.NoError(t, err)
		tq.Equal(t, `{"foo":1,"fooBar":true,"fooBarBaz":[1,2,3]}`, string(actual))
	})
}
