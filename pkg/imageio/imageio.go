// Open source image resizer coded by kasuraSH
package imageio

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/kasurarykerion/golangresizer/internal/validator"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported image format")
	ErrFileOpen          = errors.New("failed to open file")
	ErrFileCreate        = errors.New("failed to create file")
	ErrDecode            = errors.New("failed to decode image")
	ErrEncode            = errors.New("failed to encode image")
)

const (
	// JPEGQuality defines the JPEG encoding quality (1-100)
	JPEGQuality = 95
	// PNGCompression defines PNG compression level
	PNGCompression = png.DefaultCompression
)

// SupportedFormats lists all supported image formats
var SupportedFormats = []string{".jpg", ".jpeg", ".png", ".bmp", ".tiff", ".tif", ".webp"}

// LoadImage loads an image from the specified file path
func LoadImage(path string) (image.Image, error) {
	// Assertion 1: Validate path
	if err := validator.ValidatePath(path); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFileOpen, err)
	}

	// Assertion 2: Open file with error checking
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFileOpen, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log error but don't override return error
		}
	}()

	// Assertion 3: Get file info to validate size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("%w: cannot stat file: %v", ErrFileOpen, err)
	}

	// Assertion 4: Check file size is within limits
	if fileInfo.Size() > validator.MaxFileSize {
		return nil, fmt.Errorf("%w: file too large", ErrFileOpen)
	}

	// Determine format from extension
	ext := strings.ToLower(filepath.Ext(path))
	
	var img image.Image
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	case ".bmp":
		img, err = bmp.Decode(file)
	case ".tiff", ".tif":
		img, err = tiff.Decode(file)
	case ".webp":
		img, err = webp.Decode(file)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, ext)
	}

	// Assertion 5: Check decode result
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecode, err)
	}

	// Assertion 6: Validate decoded image
	if img == nil {
		return nil, fmt.Errorf("%w: decoded image is nil", ErrDecode)
	}

	return img, nil
}

// SaveImage saves an image to the specified file path
func SaveImage(path string, img image.Image) error {
	// Assertion 1: Validate path
	if err := validator.ValidatePath(path); err != nil {
		return fmt.Errorf("%w: %v", ErrFileCreate, err)
	}

	// Assertion 2: Validate image is not nil
	if img == nil {
		return fmt.Errorf("%w: image is nil", ErrFileCreate)
	}

	// Assertion 3: Validate image dimensions
	bounds := img.Bounds()
	if err := validator.ValidateDimensions(bounds.Dx(), bounds.Dy()); err != nil {
		return fmt.Errorf("%w: invalid dimensions: %v", ErrFileCreate, err)
	}

	// Create output directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("%w: cannot create directory: %v", ErrFileCreate, err)
	}

	// Assertion 4: Create output file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFileCreate, err)
	}

	// Ensure file is closed and check for errors
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Error during close
		}
	}()

	// Determine format from extension
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".jpg", ".jpeg":
		// Assertion 5: Check JPEG encode
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: JPEGQuality})
	case ".png":
		// Assertion 6: Check PNG encode
		encoder := &png.Encoder{CompressionLevel: PNGCompression}
		err = encoder.Encode(file, img)
	case ".bmp":
		// Assertion 7: Check BMP encode
		err = bmp.Encode(file, img)
	case ".tiff", ".tif":
		// Assertion 8: Check TIFF encode
		err = tiff.Encode(file, img, &tiff.Options{Compression: tiff.Deflate})
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedFormat, ext)
	}

	// Assertion 9: Check encode result
	if err != nil {
		return fmt.Errorf("%w: %v", ErrEncode, err)
	}

	return nil
}

// GetImageFormat returns the format of an image file
func GetImageFormat(path string) (string, error) {
	// Assertion 1: Validate path
	if err := validator.ValidatePath(path); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(path))

	// Assertion 2: Check if format is supported
	supported := false
	for i := 0; i < len(SupportedFormats); i++ {
		if ext == SupportedFormats[i] {
			supported = true
			break
		}
	}

	if !supported {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedFormat, ext)
	}

	return ext, nil
}
