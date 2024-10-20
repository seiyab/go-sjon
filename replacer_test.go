package sjon_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/seiyab/go-sjon"
	"github.com/stretchr/testify/require"
)

func TestReplacer(t *testing.T) {
	s := sjon.NewSerializer().
		With(sjon.Replacer(func(d time.Time) string {
			return d.Format("2006-01-02")
		}))

	t.Run("top level", func(t *testing.T) {
		d := time.Date(2024, 10, 20, 1, 2, 3, 4, time.UTC)
		actual, err := s.Marshal(d)
		require.NoError(t, err)

		tq.Equal(t, `"2024-10-20"`, string(actual))
	})

	t.Run("slice", func(t *testing.T) {
		ds := []time.Time{
			time.Date(2024, 10, 20, 1, 2, 3, 4, time.UTC),
			time.Date(2024, 10, 21, 1, 2, 3, 4, time.FixedZone("test1", 8)),
			time.Date(2024, 10, 22, 1, 2, 3, 4, time.FixedZone("test2", 16)),
		}
		actual, err := s.Marshal(ds)
		require.NoError(t, err)

		tq.Equal(t, `["2024-10-20","2024-10-21","2024-10-22"]`, string(actual))
	})

	t.Run("map", func(t *testing.T) {
		m := map[time.Time]time.Time{
			time.Date(2021, 2, 3, 4, 5, 6, 7, time.UTC): time.Date(2022, 3, 4, 5, 6, 7, 8, time.UTC),
		}
		actual, err := s.Marshal(m)
		require.NoError(t, err)

		tq.Equal(t, `{"2021-02-03":"2022-03-04"}`, string(actual))
	})
}

func ExampleReplacer() {
	s := sjon.NewSerializer().
		With(sjon.Replacer(func(d time.Time) string {
			return d.Format("2006-01-02")
		}))
		
	b, _ := s.Marshal([]time.Time{
			time.Date(2024, 10, 20, 1, 2, 3, 4, time.UTC),
			time.Date(2024, 10, 21, 1, 2, 3, 4, time.UTC),
			time.Date(2024, 10, 22, 1, 2, 3, 4, time.UTC),
		})
	fmt.Println(string(b))
	// Output: ["2024-10-20","2024-10-21","2024-10-22"]
}