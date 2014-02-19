// Package newmath is a trivial example package.
package newmath

// Sqrt returns an approximation of the square root of x.
func Sqrt(x float64) float64 {
	z := 1.0
	for i := 1; i < 1000; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}
