package math

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFib(t *testing.T) {
	require := require.New(t)

	tc := []struct {
		n    int
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
		{8, 21},
	}

	for _, tt := range tc {
		name := fmt.Sprintf("fib(%d) should be %d", tt.n, tt.want)

		t.Run(name, func(t *testing.T) {
			require.Equal(tt.want, Fib(tt.n))
		})
	}
}
