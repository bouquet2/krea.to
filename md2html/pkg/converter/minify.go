package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
)

// MinifyCSS minifies a CSS file and writes it to the output path
func MinifyCSS(inputPath, outputPath string) error {
	logger := log.With().Str("input", inputPath).Str("output", outputPath).Logger()
	logger.Debug().Msg("Starting CSS minification")

	// Read input file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading CSS file: %v", err)
	}
	logger.Debug().Int("original_size", len(content)).Msg("Read CSS file")

	// Create minifier
	m := minify.New()
	m.AddFunc("text/css", css.Minify)

	// Minify CSS
	minified, err := m.String("text/css", string(content))
	if err != nil {
		return fmt.Errorf("error minifying CSS: %v", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	// Write minified content
	if err := os.WriteFile(outputPath, []byte(minified), 0644); err != nil {
		return fmt.Errorf("error writing minified CSS: %v", err)
	}

	originalSize := len(content)
	minifiedSize := len(minified)
	savings := float64(originalSize-minifiedSize) / float64(originalSize) * 100

	logger.Debug().
		Int("minified_size", minifiedSize).
		Float64("savings_percent", savings).
		Msg("CSS minification complete")

	fmt.Printf("Minified CSS: %s -> %s (%.1f%% reduction)\n", inputPath, outputPath, savings)
	return nil
}

// MinifyCSSDir minifies all CSS files in a directory
func MinifyCSSDir(inputDir, outputDir string) error {
	logger := log.With().Str("input_dir", inputDir).Str("output_dir", outputDir).Logger()
	logger.Debug().Msg("Starting CSS directory minification")

	// Walk through input directory
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process CSS files
		if !strings.HasSuffix(strings.ToLower(path), ".css") {
			return nil
		}

		// Calculate output path
		relPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return fmt.Errorf("error calculating relative path: %v", err)
		}
		outputPath := filepath.Join(outputDir, relPath)

		// Minify the file
		return MinifyCSS(path, outputPath)
	})
}

// MinifyJS minifies a JavaScript file and writes it to the output path
func MinifyJS(inputPath, outputPath string) error {
	logger := log.With().Str("input", inputPath).Str("output", outputPath).Logger()
	logger.Debug().Msg("Starting JS minification")

	// Read input file
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading JS file: %v", err)
	}
	logger.Debug().Int("original_size", len(content)).Msg("Read JS file")

	// Create minifier
	m := minify.New()
	m.AddFunc("application/javascript", js.Minify)

	// Minify JS
	minified, err := m.String("application/javascript", string(content))
	if err != nil {
		return fmt.Errorf("error minifying JS: %v", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	// Write minified content
	if err := os.WriteFile(outputPath, []byte(minified), 0644); err != nil {
		return fmt.Errorf("error writing minified JS: %v", err)
	}

	originalSize := len(content)
	minifiedSize := len(minified)
	savings := float64(originalSize-minifiedSize) / float64(originalSize) * 100

	logger.Debug().
		Int("minified_size", minifiedSize).
		Float64("savings_percent", savings).
		Msg("JS minification complete")

	fmt.Printf("Minified JS: %s -> %s (%.1f%% reduction)\n", inputPath, outputPath, savings)
	return nil
}

// MinifyJSDir minifies all JavaScript files in a directory
func MinifyJSDir(inputDir, outputDir string) error {
	logger := log.With().Str("input_dir", inputDir).Str("output_dir", outputDir).Logger()
	logger.Debug().Msg("Starting JS directory minification")

	// Walk through input directory
	return filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process JS files
		if !strings.HasSuffix(strings.ToLower(path), ".js") {
			return nil
		}

		// Calculate output path
		relPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return fmt.Errorf("error calculating relative path: %v", err)
		}
		outputPath := filepath.Join(outputDir, relPath)

		// Minify the file
		return MinifyJS(path, outputPath)
	})
}
