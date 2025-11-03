// Open source image resizer coded by kasuraSH
package validator

import (
	"errors"
	"fmt"
)

const (
	// MaxImageDimension prevents integer overflow and memory exhaustion
	MaxImageDimension = 65535
	MinImageDimension = 1
	MaxFileSize       = 1073741824 // 1GB limit
)

var (
	ErrInvalidDimension = errors.New("dimension out of valid range")
	ErrInvalidPath      = errors.New("invalid file path")
	ErrNilPointer       = errors.New("nil pointer detected")
)

// ValidateDimensions checks if image dimensions are within safe bounds
func ValidateDimensions(width, height int) error {
	// Assertion 1: Check minimum bounds
	if width < MinImageDimension || height < MinImageDimension {
		return fmt.Errorf("%w: dimensions must be >= %d", ErrInvalidDimension, MinImageDimension)
	}

	// Assertion 2: Check maximum bounds to prevent overflow
	if width > MaxImageDimension || height > MaxImageDimension {
		return fmt.Errorf("%w: dimensions must be <= %d", ErrInvalidDimension, MaxImageDimension)
	}

	// Assertion 3: Check for potential overflow in multiplication
	if int64(width)*int64(height) > int64(MaxImageDimension*MaxImageDimension) {
		return fmt.Errorf("%w: total pixels exceed safe limit", ErrInvalidDimension)
	}

	return nil
}

// ValidatePath checks if file path is non-empty and within length limits
func ValidatePath(path string) error {
	const maxPathLength = 4096

	// Assertion 1: Check for empty path
	if path == "" {
		return fmt.Errorf("%w: path cannot be empty", ErrInvalidPath)
	}

	// Assertion 2: Check path length to prevent buffer issues
	if len(path) > maxPathLength {
		return fmt.Errorf("%w: path exceeds maximum length", ErrInvalidPath)
	}

	return nil
}

// ValidateResizeRatio checks if resize ratio is within acceptable range
func ValidateResizeRatio(originalWidth, originalHeight, newWidth, newHeight int) error {
	const maxScaleFactor = 16.0
	const minScaleFactor = 0.0625 // 1/16

	// Assertion 1: Validate all dimensions first
	if err := ValidateDimensions(originalWidth, originalHeight); err != nil {
		return err
	}

	// Assertion 2: Validate new dimensions
	if err := ValidateDimensions(newWidth, newHeight); err != nil {
		return err
	}

	// Assertion 3: Check scale factors to prevent extreme scaling
	widthRatio := float64(newWidth) / float64(originalWidth)
	heightRatio := float64(newHeight) / float64(originalHeight)

	if widthRatio > maxScaleFactor || widthRatio < minScaleFactor {
		return fmt.Errorf("%w: width scale factor out of range", ErrInvalidDimension)
	}

	if heightRatio > maxScaleFactor || heightRatio < minScaleFactor {
		return fmt.Errorf("%w: height scale factor out of range", ErrInvalidDimension)
	}

	return nil
}
