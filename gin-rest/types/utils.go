package types

import (
	"golang.org/x/exp/constraints"
)

// Max returns the maximum of two values of any ordered type.
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
