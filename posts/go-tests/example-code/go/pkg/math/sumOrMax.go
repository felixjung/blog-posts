package math

func SumOrMax(a, b, max int) int {
	sum := a + b

	if sum < max {
		return sum
	}

	return max
}
