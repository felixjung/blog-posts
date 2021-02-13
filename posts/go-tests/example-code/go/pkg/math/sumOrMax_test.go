package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSumOrMax(t *testing.T) {
	t.Run("when the sum is below the maximum", func(t *testing.T) {
		a := 1
		b := 2
		max := 4

		want := 3

		have := SumOrMax(a, b, max)

		require.Equal(t, want, have, "should return the sum")
	})

	t.Run("when the sum is at or above the maximum", func(t *testing.T) {
		a := 5
		b := 10
		max := 4

		want := max

		have := SumOrMax(a, b, max)

		require.Equal(t, want, have, "should return the maximum")
	})
}
