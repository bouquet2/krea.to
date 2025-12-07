package converter

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// processWikiLinks replaces wiki-style links [[path/to/page]] with HTML links
func processWikiLinks(content []byte, basePath string, inputRoot string, outputRoot string) []byte {
	// Regular expression to match wiki-style links: [[path/to/page]] or [[path/to/folder/]]
	re := regexp.MustCompile(`\[\[(.*?)\]\]`)

	return re.ReplaceAllFunc(content, func(match []byte) []byte {
		// Extract the link path (remove [[ and ]])
		linkPath := string(match[2 : len(match)-2])
		displayText := linkPath

		// Determine if it's a folder link (ends with /)
		isFolder := strings.HasSuffix(linkPath, "/")

		// Normalize the link path - remove trailing slash for path calculation
		normalizedPath := linkPath
		if isFolder && len(normalizedPath) > 1 {
			normalizedPath = normalizedPath[:len(normalizedPath)-1]
		}

		// Calculate target file path
		var targetPath string
		if strings.HasPrefix(normalizedPath, "/") {
			// Absolute path (relative to input root)
			targetPath = normalizedPath[1:] // Remove leading slash
		} else {
			// Relative path from the current file's directory
			targetPath = filepath.Join(filepath.Dir(basePath), normalizedPath)
			targetPath, _ = filepath.Rel(inputRoot, targetPath)
		}

		// Determine the relative path from current file to target
		sourceDir := filepath.Dir(basePath)
		sourceRelToRoot, _ := filepath.Rel(inputRoot, sourceDir)
		sourceOutputDir := filepath.Join(outputRoot, sourceRelToRoot)

		targetOutputDir := filepath.Join(outputRoot, filepath.Dir(targetPath))

		relPath, err := filepath.Rel(sourceOutputDir, targetOutputDir)
		if err != nil {
			// Fallback to the original path if there's an error
			relPath = targetPath
		}

		// Convert backslashes to forward slashes for URLs
		relPath = strings.ReplaceAll(relPath, string(filepath.Separator), "/")

		// If not empty and doesn't end with a slash, add one
		if relPath != "" && !strings.HasSuffix(relPath, "/") {
			relPath += "/"
		}

		// Get the file name from the target path
		fileName := filepath.Base(targetPath)

		// Build the HTML link
		htmlLink := ""
		if isFolder {
			// For simple folder references like [[Nim/]]
			if fileName == normalizedPath {
				// Simple folder reference - directly use the folder name
				htmlLink = fmt.Sprintf("<a href=\"./%s/index.html\">%s</a>", normalizedPath, displayText)
			} else {
				// More complex path for nested folders
				if relPath == "./" || relPath == "" {
					// Same directory case
					htmlLink = fmt.Sprintf("<a href=\"./%s/index.html\">%s</a>", filepath.Base(normalizedPath), displayText)
				} else {
					// Nested directory case
					htmlLink = fmt.Sprintf("<a href=\"%s%s/index.html\">%s</a>", relPath, filepath.Base(normalizedPath), displayText)
				}
			}
		} else {
			// Link to the specific page
			htmlLink = fmt.Sprintf("<a href=\"%s%s.html\">%s</a>", relPath, fileName, displayText)
		}

		return []byte(htmlLink)
	})
}

// calculatePathDepth determines the directory depth based on input and output directories
func calculatePathDepth(inputDir, inputRoot string) (int, error) {
	// Get absolute paths to ensure accurate comparison
	absInputRoot, err := filepath.Abs(inputRoot)
	if err != nil {
		return 0, fmt.Errorf("error getting absolute path for root: %v", err)
	}

	absCurrentDir, err := filepath.Abs(inputDir)
	if err != nil {
		return 0, fmt.Errorf("error getting absolute path for current dir: %v", err)
	}

	// Compare the paths to determine depth
	if absInputRoot == absCurrentDir {
		return 0, nil
	}

	// Get the relative path from root to current directory
	relPath, err := filepath.Rel(absInputRoot, absCurrentDir)
	if err != nil {
		// Handle special case: on some platforms, paths may not be directly relatable

		// Try using string manipulation as fallback
		absRootStr := strings.ReplaceAll(absInputRoot, string(filepath.Separator), "/")
		absCurrentStr := strings.ReplaceAll(absCurrentDir, string(filepath.Separator), "/")

		if strings.HasPrefix(absCurrentStr, absRootStr) {
			// If current dir starts with root, count the separators in the remaining path
			remaining := strings.TrimPrefix(absCurrentStr, absRootStr)
			remaining = strings.TrimPrefix(remaining, "/")
			if remaining == "" {
				return 0, nil
			}

			depth := strings.Count(remaining, "/") + 1
			return depth, nil
		}

		return 0, fmt.Errorf("error calculating relative path: %v", err)
	}

	if relPath == "." {
		return 0, nil
	}

	// Count path components to determine depth
	components := strings.Split(relPath, string(filepath.Separator))
	// Filter out empty components
	var nonEmptyComponents []string
	for _, comp := range components {
		if comp != "" {
			nonEmptyComponents = append(nonEmptyComponents, comp)
		}
	}

	depth := len(nonEmptyComponents)

	// Additional check for forward slashes on Windows
	if runtime.GOOS == "windows" && depth == 0 {
		components = strings.Split(relPath, "/")
		nonEmptyComponents = nil
		for _, comp := range components {
			if comp != "" {
				nonEmptyComponents = append(nonEmptyComponents, comp)
			}
		}
		depth = len(nonEmptyComponents)
	}

	// Check for special case with blog subdirectories
	if strings.Contains(absCurrentDir, "blog/") || strings.Contains(absCurrentDir, "blog"+string(filepath.Separator)) {
		segments := strings.Split(absCurrentDir, "blog"+string(filepath.Separator))
		if len(segments) > 1 && segments[1] != "" {
			subSegments := strings.Split(segments[1], string(filepath.Separator))
			nonEmptySubSegments := 0
			for _, comp := range subSegments {
				if comp != "" {
					nonEmptySubSegments++
				}
			}

			if nonEmptySubSegments > 0 {
				// We're in a blog subfolder
				if nonEmptySubSegments == 1 {
					// Direct child of blog folder
					return 1, nil
				} else {
					// Nested folder within blog structure
					return nonEmptySubSegments, nil
				}
			}
		}
	}

	return depth, nil
}

// adjustPaths modifies CSS and JS paths based on the current directory depth
func adjustPaths(config Config, depth int, outputDir string) Config {
	adjustedConfig := config

	// Special case for blog/category structure - explicit detection
	// Find blog in the output path components (may be prefixed with dist/ or similar)
	outComponents := strings.Split(outputDir, string(filepath.Separator))
	blogIndex := -1
	for i, comp := range outComponents {
		if comp == "blog" {
			blogIndex = i
			break
		}
	}

	if blogIndex >= 0 {
		blogComponents := outComponents[blogIndex:]
		// For first-level blog categories (blog/X), ALWAYS use ../../ prefix regardless of depth calculation
		if len(blogComponents) == 2 {
			adjustedConfig.CSSPath = "../../css/style.css"
			adjustedConfig.JSPath = "../../js/script.js"
			return adjustedConfig
		} else if len(blogComponents) > 2 {
			// For nested blog structure (blog/X/Y), count the components to determine proper path
			nestingDepth := len(blogComponents) - 1 // blog/Nim/SubDir would be depth 2
			prefix := ""
			for i := 0; i < nestingDepth; i++ {
				prefix += "../"
			}
			adjustedConfig.CSSPath = prefix + "css/style.css"
			adjustedConfig.JSPath = prefix + "js/script.js"
			return adjustedConfig
		}
	}

	// For any non-blog structure, process normally
	if !strings.HasPrefix(config.CSSPath, "http") {
		if depth > 0 {
			prefix := ""
			for i := 0; i < depth; i++ {
				prefix += "../"
			}

			// If it's an absolute path (starts with /), remove the leading slash
			if strings.HasPrefix(config.CSSPath, "/") {
				adjustedConfig.CSSPath = prefix + config.CSSPath[1:]
			} else {
				// This is a relative path, reference from root
				adjustedConfig.CSSPath = prefix + config.CSSPath
			}
		} else if strings.HasPrefix(config.CSSPath, "/") {
			// If we're at the root level but have an absolute path, just remove the slash
			adjustedConfig.CSSPath = config.CSSPath[1:]
		}
	}

	if !strings.HasPrefix(config.JSPath, "http") {
		if depth > 0 {
			prefix := ""
			for i := 0; i < depth; i++ {
				prefix += "../"
			}

			// If it's an absolute path (starts with /), remove the leading slash
			if strings.HasPrefix(config.JSPath, "/") {
				adjustedConfig.JSPath = prefix + config.JSPath[1:]
			} else {
				// This is a relative path, reference from root
				adjustedConfig.JSPath = prefix + config.JSPath
			}
		} else if strings.HasPrefix(config.JSPath, "/") {
			// If we're at the root level but have an absolute path, just remove the slash
			adjustedConfig.JSPath = config.JSPath[1:]
		}
	}

	return adjustedConfig
}

// calculateBackURL determines the URL for going back to the parent directory
func calculateBackURL(depth int) string {
	if depth <= 0 {
		return ""
	}

	if depth == 1 {
		return "../index.html"
	}

	backURL := ""
	for i := 0; i < depth; i++ {
		backURL += "../"
	}
	backURL += "index.html"

	return backURL
}
