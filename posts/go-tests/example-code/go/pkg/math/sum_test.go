package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSum(t *testing.T) {
	a := 1
	b := 2

	want := 3

	got := Sum(a, b)

	if want != got {
		t.Fatalf("Sum: wanted %d, got %d", want, got)
	}
}

func TestSumWithTestify(t *testing.T) {
	a := 1
	b := 2

	want := 3

	got := Sum(a, b)

	require.Equal(t, want, got)
}

func TestSumWithAssertionsInstance(t *testing.T) {
	require := require.New(t)

	a := 1
	b := 2

	want := 3

	got := Sum(a, b)

	require.Equal(want, got)
}

func TestSumWithSubtest(t *testing.T) {
}
