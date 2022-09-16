package utils

func Pow(x float64, n int) float64 {
	if n == 0 {
		return 1
	}

	return x * Pow(x, n-1)
}
