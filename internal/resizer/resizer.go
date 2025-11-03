// Open source image resizer coded by kasuraSH
package resizer

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/kasuraSH/kasurarykerion/internal/interpolation"
	"github.com/kasuraSH/kasurarykerion/internal/validator"
)

var (
	ErrNilImage       = errors.New("nil image provided")
	ErrInvalidBounds  = errors.New("invalid image bounds")
	ErrResizeFailed   = errors.New("resize operation failed")
	ErrUnsupportedBit = errors.New("unsupported bit depth")
)

// Config holds resize operation parameters
type Config struct {
	TargetWidth  int
	TargetHeight int
	Quality      int // 0-100, currently unused but reserved for future
}

// Resizer handles image resizing operations
type Resizer struct {
	config Config
}

// NewResizer creates a new resizer instance
func NewResizer(cfg Config) (*Resizer, error) {
	// Assertion 1: Validate target dimensions
	if err := validator.ValidateDimensions(cfg.TargetWidth, cfg.TargetHeight); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Assertion 2: Validate quality parameter
	if cfg.Quality < 0 || cfg.Quality > 100 {
		cfg.Quality = 100 // Default to maximum quality
	}

	return &Resizer{config: cfg}, nil
}

// Resize performs the image resizing operation
func (r *Resizer) Resize(src image.Image) (image.Image, error) {
	// Assertion 1: Validate input image
	if src == nil {
		return nil, ErrNilImage
	}

	bounds := src.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	// Assertion 2: Validate source dimensions
	if err := validator.ValidateDimensions(srcWidth, srcHeight); err != nil {
		return nil, fmt.Errorf("invalid source dimensions: %w", err)
	}

	// Assertion 3: Validate resize ratio
	if err := validator.ValidateResizeRatio(srcWidth, srcHeight, r.config.TargetWidth, r.config.TargetHeight); err != nil {
		return nil, fmt.Errorf("invalid resize ratio: %w", err)
	}

	// Determine bit depth and process accordingly
	switch src.ColorModel() {
	case color.RGBAModel, color.NRGBAModel:
		return r.resizeRGBA(src, srcWidth, srcHeight)
	case color.RGBA64Model, color.NRGBA64Model:
		return r.resizeRGBA64(src, srcWidth, srcHeight)
	case color.GrayModel:
		return r.resizeGray(src, srcWidth, srcHeight)
	case color.Gray16Model:
		return r.resizeGray16(src, srcWidth, srcHeight)
	default:
		// Convert to RGBA for unsupported formats
		return r.resizeRGBA(src, srcWidth, srcHeight)
	}
}

// resizeRGBA handles 8-bit RGBA images
func (r *Resizer) resizeRGBA(src image.Image, srcWidth, srcHeight int) (*image.RGBA, error) {
	// Assertion 1: Validate we can create destination image
	if err := validator.ValidateDimensions(r.config.TargetWidth, r.config.TargetHeight); err != nil {
		return nil, err
	}

	dst := image.NewRGBA(image.Rect(0, 0, r.config.TargetWidth, r.config.TargetHeight))

	xRatio := float64(srcWidth) / float64(r.config.TargetWidth)
	yRatio := float64(srcHeight) / float64(r.config.TargetHeight)

	for y := 0; y < r.config.TargetHeight; y++ {
		srcY := (float64(y) + 0.5) * yRatio

		for x := 0; x < r.config.TargetWidth; x++ {
			srcX := (float64(x) + 0.5) * xRatio

			// Process each color channel
			r, g, b, a, err := r.sampleRGBA(src, srcX, srcY, srcWidth, srcHeight)
			if err != nil {
				return nil, fmt.Errorf("sampling failed at (%d,%d): %w", x, y, err)
			}

			dst.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	return dst, nil
}

// sampleRGBA performs bicubic sampling for RGBA channels
func (r *Resizer) sampleRGBA(src image.Image, x, y float64, width, height int) (uint8, uint8, uint8, uint8, error) {
	// Calculate kernel bounds
	startX, endX, err := interpolation.CalculateKernelBounds(x, width)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	startY, endY, err := interpolation.CalculateKernelBounds(y, height)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Fractional parts for interpolation
	dx := x - float64(int(x))
	dy := y - float64(int(y))

	var rPixels, gPixels, bPixels, aPixels [interpolation.KernelSize][interpolation.KernelSize]float64

	kernelY := 0
	for srcY := startY; srcY < endY; srcY++ {
		safeY := interpolation.GetSafeIndex(srcY, height)
		kernelX := 0

		for srcX := startX; srcX < endX; srcX++ {
			safeX := interpolation.GetSafeIndex(srcX, width)

			c := src.At(safeX, safeY)
			r32, g32, b32, a32 := c.RGBA()

			// Convert from 16-bit to 8-bit
			rPixels[kernelY][kernelX] = float64(r32 >> 8)
			gPixels[kernelY][kernelX] = float64(g32 >> 8)
			bPixels[kernelY][kernelX] = float64(b32 >> 8)
			aPixels[kernelY][kernelX] = float64(a32 >> 8)

			kernelX++
		}
		kernelY++
	}

	// Perform bicubic interpolation for each channel
	rVal, err := interpolation.InterpolateBicubic(rPixels, dx, dy)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	gVal, err := interpolation.InterpolateBicubic(gPixels, dx, dy)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	bVal, err := interpolation.InterpolateBicubic(bPixels, dx, dy)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	aVal, err := interpolation.InterpolateBicubic(aPixels, dx, dy)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return interpolation.ClampUint8(rVal),
		interpolation.ClampUint8(gVal),
		interpolation.ClampUint8(bVal),
		interpolation.ClampUint8(aVal),
		nil
}

// resizeRGBA64 handles 16-bit RGBA images
func (r *Resizer) resizeRGBA64(src image.Image, srcWidth, srcHeight int) (*image.RGBA64, error) {
	// Assertion 1: Validate dimensions
	if err := validator.ValidateDimensions(r.config.TargetWidth, r.config.TargetHeight); err != nil {
		return nil, err
	}

	dst := image.NewRGBA64(image.Rect(0, 0, r.config.TargetWidth, r.config.TargetHeight))

	xRatio := float64(srcWidth) / float64(r.config.TargetWidth)
	yRatio := float64(srcHeight) / float64(r.config.TargetHeight)

	for y := 0; y < r.config.TargetHeight; y++ {
		srcY := (float64(y) + 0.5) * yRatio

		for x := 0; x < r.config.TargetWidth; x++ {
			srcX := (float64(x) + 0.5) * xRatio

			r, g, b, a, err := r.sampleRGBA64(src, srcX, srcY, srcWidth, srcHeight)
			if err != nil {
				return nil, fmt.Errorf("sampling failed at (%d,%d): %w", x, y, err)
			}

			dst.SetRGBA64(x, y, color.RGBA64{R: r, G: g, B: b, A: a})
		}
	}

	return dst, nil
}

// sampleRGBA64 performs bicubic sampling for 16-bit RGBA
func (r *Resizer) sampleRGBA64(src image.Image, x, y float64, width, height int) (uint16, uint16, uint16, uint16, error) {
	startX, endX, err := interpolation.CalculateKernelBounds(x, width)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	startY, endY, err := interpolation.CalculateKernelBounds(y, height)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	dx := x - float64(int(x))
	dy := y - float64(int(y))

	var rPixels, gPixels, bPixels, aPixels [interpolation.KernelSize][interpolation.KernelSize]float64

	kernelY := 0
	for srcY := startY; srcY < endY; srcY++ {
		safeY := interpolation.GetSafeIndex(srcY, height)
		kernelX := 0

		for srcX := startX; srcX < endX; srcX++ {
			safeX := interpolation.GetSafeIndex(srcX, width)

			c := src.At(safeX, safeY)
			r32, g32, b32, a32 := c.RGBA()

			rPixels[kernelY][kernelX] = float64(r32)
			gPixels[kernelY][kernelX] = float64(g32)
			bPixels[kernelY][kernelX] = float64(b32)
			aPixels[kernelY][kernelX] = float64(a32)

			kernelX++
		}
		kernelY++
	}

	rVal, err := interpolation.InterpolateBicubic(rPixels, dx, dy)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	gVal, err := interpolation.InterpolateBicubic(gPixels, dx, dy)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	bVal, err := interpolation.InterpolateBicubic(bPixels, dx, dy)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	aVal, err := interpolation.InterpolateBicubic(aPixels, dx, dy)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return interpolation.ClampUint16(rVal),
		interpolation.ClampUint16(gVal),
		interpolation.ClampUint16(bVal),
		interpolation.ClampUint16(aVal),
		nil
}

// resizeGray handles 8-bit grayscale images
func (r *Resizer) resizeGray(src image.Image, srcWidth, srcHeight int) (*image.Gray, error) {
	if err := validator.ValidateDimensions(r.config.TargetWidth, r.config.TargetHeight); err != nil {
		return nil, err
	}

	dst := image.NewGray(image.Rect(0, 0, r.config.TargetWidth, r.config.TargetHeight))

	xRatio := float64(srcWidth) / float64(r.config.TargetWidth)
	yRatio := float64(srcHeight) / float64(r.config.TargetHeight)

	for y := 0; y < r.config.TargetHeight; y++ {
		srcY := (float64(y) + 0.5) * yRatio

		for x := 0; x < r.config.TargetWidth; x++ {
			srcX := (float64(x) + 0.5) * xRatio

			grayVal, err := r.sampleGray(src, srcX, srcY, srcWidth, srcHeight)
			if err != nil {
				return nil, fmt.Errorf("sampling failed at (%d,%d): %w", x, y, err)
			}

			dst.SetGray(x, y, color.Gray{Y: grayVal})
		}
	}

	return dst, nil
}

// sampleGray performs bicubic sampling for grayscale
func (r *Resizer) sampleGray(src image.Image, x, y float64, width, height int) (uint8, error) {
	startX, endX, err := interpolation.CalculateKernelBounds(x, width)
	if err != nil {
		return 0, err
	}

	startY, endY, err := interpolation.CalculateKernelBounds(y, height)
	if err != nil {
		return 0, err
	}

	dx := x - float64(int(x))
	dy := y - float64(int(y))

	var pixels [interpolation.KernelSize][interpolation.KernelSize]float64

	kernelY := 0
	for srcY := startY; srcY < endY; srcY++ {
		safeY := interpolation.GetSafeIndex(srcY, height)
		kernelX := 0

		for srcX := startX; srcX < endX; srcX++ {
			safeX := interpolation.GetSafeIndex(srcX, width)
			c := src.At(safeX, safeY)
			gray, _, _, _ := c.RGBA()
			pixels[kernelY][kernelX] = float64(gray >> 8)
			kernelX++
		}
		kernelY++
	}

	val, err := interpolation.InterpolateBicubic(pixels, dx, dy)
	if err != nil {
		return 0, err
	}

	return interpolation.ClampUint8(val), nil
}

// resizeGray16 handles 16-bit grayscale images
func (r *Resizer) resizeGray16(src image.Image, srcWidth, srcHeight int) (*image.Gray16, error) {
	if err := validator.ValidateDimensions(r.config.TargetWidth, r.config.TargetHeight); err != nil {
		return nil, err
	}

	dst := image.NewGray16(image.Rect(0, 0, r.config.TargetWidth, r.config.TargetHeight))

	xRatio := float64(srcWidth) / float64(r.config.TargetWidth)
	yRatio := float64(srcHeight) / float64(r.config.TargetHeight)

	for y := 0; y < r.config.TargetHeight; y++ {
		srcY := (float64(y) + 0.5) * yRatio

		for x := 0; x < r.config.TargetWidth; x++ {
			srcX := (float64(x) + 0.5) * xRatio

			grayVal, err := r.sampleGray16(src, srcX, srcY, srcWidth, srcHeight)
			if err != nil {
				return nil, fmt.Errorf("sampling failed at (%d,%d): %w", x, y, err)
			}

			dst.SetGray16(x, y, color.Gray16{Y: grayVal})
		}
	}

	return dst, nil
}

// sampleGray16 performs bicubic sampling for 16-bit grayscale
func (r *Resizer) sampleGray16(src image.Image, x, y float64, width, height int) (uint16, error) {
	startX, endX, err := interpolation.CalculateKernelBounds(x, width)
	if err != nil {
		return 0, err
	}

	startY, endY, err := interpolation.CalculateKernelBounds(y, height)
	if err != nil {
		return 0, err
	}

	dx := x - float64(int(x))
	dy := y - float64(int(y))

	var pixels [interpolation.KernelSize][interpolation.KernelSize]float64

	kernelY := 0
	for srcY := startY; srcY < endY; srcY++ {
		safeY := interpolation.GetSafeIndex(srcY, height)
		kernelX := 0

		for srcX := startX; srcX < endX; srcX++ {
			safeX := interpolation.GetSafeIndex(srcX, width)
			c := src.At(safeX, safeY)
			gray, _, _, _ := c.RGBA()
			pixels[kernelY][kernelX] = float64(gray)
			kernelX++
		}
		kernelY++
	}

	val, err := interpolation.InterpolateBicubic(pixels, dx, dy)
	if err != nil {
		return 0, err
	}

	return interpolation.ClampUint16(val), nil
}
