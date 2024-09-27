package sjon_test

import (
	"testing"
	"time"

	"github.com/seiyab/go-sjon"
	"github.com/stretchr/testify/require"
)

func TestJSONMarshaler(t *testing.T) {
	s := sjon.NewSerializer()
	t.Run("time.Time", func(t *testing.T) {
		d := time.Date(2024, 9, 27, 1, 2, 3, 4, time.UTC)
		actual, err := s.Marshal(d)
		require.NoError(t, err)
		tq.Equal(t, `"2024-09-27T01:02:03.000000004Z"`, string(actual))
		compareStandard(t, d)
	})

	t.Run("time.Time in struct", func(t *testing.T) {
		type S struct {
			T  time.Time
			TP *time.Time
		}
		d := S{
			time.Date(2024, 9, 20, 5, 4, 3, 2, time.UTC),
			ref(time.Date(2024, 3, 18, 5, 4, 3, 2, time.UTC)),
		}
		actual, err := s.Marshal(d)
		require.NoError(t, err)
		tq.Equal(t, `{"T":"2024-09-20T05:04:03.000000002Z","TP":"2024-03-18T05:04:03.000000002Z"}`, string(actual))
		compareStandard(t, d)
	})
}
