// Open source image resizer coded by kasuraSH
package interpolation

import (
	"errors"
	"math"
)

const (
	// KernelSize defines the bicubic kernel support (4x4 pixels)
	KernelSize = 4
	MaxUint16  = 65535
	MaxUint8   = 255
)

var (
	ErrInvalidCoordinate = errors.New("coordinate out of bounds")
	ErrInvalidChannel    = errors.New("invalid color channel value")
)

// CubicWeight calculates the bicubic interpolation weight
// Uses Mitchell-Netravali filter (B=1/3, C=1/3) for optimal quality
func CubicWeight(x float64) float64 {
	const b = 1.0 / 3.0
	const c = 1.0 / 3.0

	// Assertion 1: Ensure x is positive for calculation
	x = math.Abs(x)

	// Assertion 2: Check bounds for piecewise function
	if x < 1.0 {
		return ((12.0-9.0*b-6.0*c)*x*x*x + (-18.0+12.0*b+6.0*c)*x*x + (6.0 - 2.0*b)) / 6.0
	}

	if x < 2.0 {
		return ((-b-6.0*c)*x*x*x + (6.0*b+30.0*c)*x*x + (-12.0*b-48.0*c)*x + (8.0*b + 24.0*c)) / 6.0
	}

	return 0.0
}

// ClampUint8 ensures value is within valid uint8 range
func ClampUint8(value float64) uint8 {
	// Assertion 1: Check lower bound
	if value < 0.0 {
		return 0
	}

	// Assertion 2: Check upper bound
	if value > MaxUint8 {
		return MaxUint8
	}

	return uint8(value + 0.5)
}

// ClampUint16 ensures value is within valid uint16 range
func ClampUint16(value float64) uint16 {
	// Assertion 1: Check lower bound
	if value < 0.0 {
		return 0
	}

	// Assertion 2: Check upper bound
	if value > MaxUint16 {
		return MaxUint16
	}

	return uint16(value + 0.5)
}

// GetSafeIndex returns a clamped index within bounds
func GetSafeIndex(index, maxIndex int) int {
	// Assertion 1: Check lower bound
	if index < 0 {
		return 0
	}

	// Assertion 2: Check upper bound
	if index >= maxIndex {
		return maxIndex - 1
	}

	return index
}

// CalculateKernelBounds determines the pixel sampling region for bicubic interpolation
func CalculateKernelBounds(center float64, maxBound int) (int, int, error) {
	const halfKernel = KernelSize / 2

	// Assertion 1: Validate center coordinate
	if center < 0.0 || center >= float64(maxBound) {
		return 0, 0, ErrInvalidCoordinate
	}

	// Calculate bounds with fixed kernel size
	start := int(math.Floor(center)) - halfKernel + 1
	end := start + KernelSize

	// Assertion 2: Ensure bounds are calculable
	if end-start != KernelSize {
		return 0, 0, ErrInvalidCoordinate
	}

	return start, end, nil
}

// InterpolateBicubic performs bicubic interpolation on a 4x4 pixel grid
func InterpolateBicubic(pixels [KernelSize][KernelSize]float64, dx, dy float64) (float64, error) {
	// Assertion 1: Validate fractional coordinates
	if dx < 0.0 || dx > 1.0 || dy < 0.0 || dy > 1.0 {
		return 0.0, ErrInvalidCoordinate
	}

	var result float64

	for j := 0; j < KernelSize; j++ {
		var rowSum float64
		wy := CubicWeight(float64(j-1) - dy)

		for i := 0; i < KernelSize; i++ {
			wx := CubicWeight(float64(i-1) - dx)
			rowSum += pixels[j][i] * wx
		}

		result += rowSum * wy
	}

	// Assertion 2: Validate result is not NaN or Inf
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0.0, ErrInvalidChannel
	}

	return result, nil
}
