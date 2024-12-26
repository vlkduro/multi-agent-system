package utils

// used for comparing values of type Ordered
import (
	"math"

	"golang.org/x/exp/constraints"
)

// https://stackoverflow.com/questions/27516387/what-is-the-correct-way-to-find-the-min-between-two-integers-in-go
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Round returns the nearest integer to x
// https://dev.to/natamacm/round-numbers-in-go-5c01
func Round(x float64) int {
	return int(math.Floor(x + 0.5))
}
