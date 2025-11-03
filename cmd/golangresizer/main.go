// Open source image resizer coded by kasuraSH
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kasuraSH/golangresizer/internal/resizer"
	"github.com/kasuraSH/golangresizer/internal/validator"
	"github.com/kasuraSH/golangresizer/pkg/imageio"
)

const (
	// ExitSuccess indicates successful execution
	ExitSuccess = 0
	// ExitError indicates an error occurred
	ExitError = 1
	// Version of the application
	Version = "1.0.0"
)

// Config holds application configuration
type Config struct {
	InputPath  string
	OutputPath string
	Width      int
	Height     int
	ShowHelp   bool
	ShowVer    bool
}

// parseFlags parses command line flags
func parseFlags() (*Config, error) {
	cfg := &Config{}

	// Define flags
	flag.StringVar(&cfg.InputPath, "input", "", "Input image file path (required)")
	flag.StringVar(&cfg.InputPath, "i", "", "Input image file path (shorthand)")
	flag.StringVar(&cfg.OutputPath, "output", "", "Output image file path (required)")
	flag.StringVar(&cfg.OutputPath, "o", "", "Output image file path (shorthand)")
	flag.IntVar(&cfg.Width, "width", 0, "Target width in pixels (required)")
	flag.IntVar(&cfg.Width, "w", 0, "Target width in pixels (shorthand)")
	flag.IntVar(&cfg.Height, "height", 0, "Target height in pixels (required)")
	flag.IntVar(&cfg.Height, "h", 0, "Target height in pixels (shorthand)")
	flag.BoolVar(&cfg.ShowHelp, "help", false, "Show help message")
	flag.BoolVar(&cfg.ShowVer, "version", false, "Show version information")

	flag.Parse()

	// Assertion 1: Check if help or version requested
	if cfg.ShowHelp {
		return cfg, nil
	}

	if cfg.ShowVer {
		return cfg, nil
	}

	// Assertion 2: Validate required parameters
	if cfg.InputPath == "" {
		return nil, fmt.Errorf("input path is required")
	}

	if cfg.OutputPath == "" {
		return nil, fmt.Errorf("output path is required")
	}

	if cfg.Width <= 0 {
		return nil, fmt.Errorf("width must be greater than 0")
	}

	if cfg.Height <= 0 {
		return nil, fmt.Errorf("height must be greater than 0")
	}

	// Assertion 3: Validate paths
	if err := validator.ValidatePath(cfg.InputPath); err != nil {
		return nil, fmt.Errorf("invalid input path: %w", err)
	}

	if err := validator.ValidatePath(cfg.OutputPath); err != nil {
		return nil, fmt.Errorf("invalid output path: %w", err)
	}

	// Assertion 4: Validate dimensions
	if err := validator.ValidateDimensions(cfg.Width, cfg.Height); err != nil {
		return nil, fmt.Errorf("invalid dimensions: %w", err)
	}

	return cfg, nil
}

// printHelp displays usage information
func printHelp() {
	fmt.Println("GolangResizer - High-Quality Image Resizer")
	fmt.Println("Open source image resizer coded by kasuraSH")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  golangresizer -input <file> -output <file> -width <pixels> -height <pixels>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -input, -i     Input image file path (required)")
	fmt.Println("  -output, -o    Output image file path (required)")
	fmt.Println("  -width, -w     Target width in pixels (required)")
	fmt.Println("  -height, -h    Target height in pixels (required)")
	fmt.Println("  -help          Show this help message")
	fmt.Println("  -version       Show version information")
	fmt.Println()
	fmt.Println("Supported formats:")
	fmt.Println("  Input:  JPEG, PNG, BMP, TIFF, WebP")
	fmt.Println("  Output: JPEG, PNG, BMP, TIFF")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  golangresizer -i input.jpg -o output.png -w 1920 -h 1080")
	fmt.Println("  golangresizer -input photo.png -output resized.jpg -width 800 -height 600")
}

// printVersion displays version information
func printVersion() {
	fmt.Printf("GolangResizer version %s\n", Version)
	fmt.Println("Open source image resizer coded by kasuraSH")
	fmt.Println("Built with NASA Power of 10 safety-critical coding rules")
}

// run executes the main application logic
func run(cfg *Config) error {
	// Assertion 1: Validate configuration
	if cfg == nil {
		return fmt.Errorf("configuration is nil")
	}

	// Load input image
	fmt.Printf("Loading image: %s\n", cfg.InputPath)
	img, err := imageio.LoadImage(cfg.InputPath)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	// Assertion 2: Validate loaded image
	if img == nil {
		return fmt.Errorf("loaded image is nil")
	}

	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	fmt.Printf("Source dimensions: %dx%d\n", srcWidth, srcHeight)
	fmt.Printf("Target dimensions: %dx%d\n", cfg.Width, cfg.Height)

	// Assertion 3: Validate resize ratio
	if err := validator.ValidateResizeRatio(srcWidth, srcHeight, cfg.Width, cfg.Height); err != nil {
		return fmt.Errorf("invalid resize parameters: %w", err)
	}

	// Create resizer
	resizerCfg := resizer.Config{
		TargetWidth:  cfg.Width,
		TargetHeight: cfg.Height,
		Quality:      100,
	}

	r, err := resizer.NewResizer(resizerCfg)
	if err != nil {
		return fmt.Errorf("failed to create resizer: %w", err)
	}

	// Assertion 4: Validate resizer was created
	if r == nil {
		return fmt.Errorf("resizer is nil")
	}

	// Perform resize operation
	fmt.Println("Resizing image using bicubic interpolation...")
	resizedImg, err := r.Resize(img)
	if err != nil {
		return fmt.Errorf("resize failed: %w", err)
	}

	// Assertion 5: Validate resized image
	if resizedImg == nil {
		return fmt.Errorf("resized image is nil")
	}

	// Verify output dimensions
	outBounds := resizedImg.Bounds()
	if outBounds.Dx() != cfg.Width || outBounds.Dy() != cfg.Height {
		return fmt.Errorf("output dimensions mismatch: got %dx%d, expected %dx%d",
			outBounds.Dx(), outBounds.Dy(), cfg.Width, cfg.Height)
	}

	// Save output image
	fmt.Printf("Saving image: %s\n", cfg.OutputPath)
	if err := imageio.SaveImage(cfg.OutputPath, resizedImg); err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	fmt.Println("Resize completed successfully!")
	return nil
}

// main is the entry point
func main() {
	// Parse command line flags
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		printHelp()
		os.Exit(ExitError)
	}

	// Handle help flag
	if cfg.ShowHelp {
		printHelp()
		os.Exit(ExitSuccess)
	}

	// Handle version flag
	if cfg.ShowVer {
		printVersion()
		os.Exit(ExitSuccess)
	}

	// Execute main logic
	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitError)
	}

	os.Exit(ExitSuccess)
}
