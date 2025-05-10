package converter

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

//go:embed templates/*
var templateFS embed.FS

// Config holds configuration for the HTML generation
type Config struct {
	TemplateFile  string
	OutputDir     string
	CSSPath       string
	JSPath        string
	SiteTitle     string
	DefaultAuthor string
	GenerateList  bool
	Recursive     bool
}

// PageData represents the data to be passed to the HTML template
type PageData struct {
	Title       string
	Content     template.HTML
	CSSPath     string
	JSPath      string
	Author      string
	Description string
	Date        string
	URL         string
	Image       string
}

// BlogPost represents a blog post entry for the index page
type BlogPost struct {
	Title       string
	Link        string
	Description string
	Date        string
}

// Directory represents a subdirectory for the index page
type Directory struct {
	Name string
	Link string
}

// IndexData represents data for the directory index template
type IndexData struct {
	Title          string
	CSSPath        string
	JSPath         string
	Posts          []BlogPost
	Directories    []Directory
	BackURL        string
	HasDirectories bool
	URL            string
	Image          string
}

// ==========================================================================
// Path Utilities
// ==========================================================================

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
	// Debug information
	fmt.Printf("Calculating path depth - Input dir: %s, Input root: %s\n", inputDir, inputRoot)

	// Get absolute paths to ensure accurate comparison
	absInputRoot, err := filepath.Abs(inputRoot)
	if err != nil {
		return 0, fmt.Errorf("error getting absolute path for root: %v", err)
	}

	absCurrentDir, err := filepath.Abs(inputDir)
	if err != nil {
		return 0, fmt.Errorf("error getting absolute path for current dir: %v", err)
	}

	fmt.Printf("Absolute paths - Root: %s, Current: %s\n", absInputRoot, absCurrentDir)

	// Compare the paths to determine depth
	if absInputRoot == absCurrentDir {
		fmt.Printf("Input directory is the root directory (depth 0)\n")
		return 0, nil
	}

	// Get the relative path from root to current directory
	relPath, err := filepath.Rel(absInputRoot, absCurrentDir)
	if err != nil {
		// Handle special case: on some platforms, paths may not be directly relatable
		fmt.Printf("Warning: could not get relative path: %v\n", err)

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
			fmt.Printf("Fallback depth calculation: %d\n", depth)
			return depth, nil
		}

		return 0, fmt.Errorf("error calculating relative path: %v", err)
	}

	if relPath == "." {
		fmt.Printf("Relative path is '.' (depth 0)\n")
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
	fmt.Printf("Path components: %v (depth %d)\n", nonEmptyComponents, depth)

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
		fmt.Printf("Windows path check - components: %v (depth %d)\n", nonEmptyComponents, depth)
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
					fmt.Printf("Blog subfolder detected, forcing depth to 1\n")
					return 1, nil
				} else {
					// Nested folder within blog structure
					fmt.Printf("Nested blog subfolder detected (level %d)\n", nonEmptySubSegments)
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

	// Debug information
	fmt.Printf("Path adjustment - depth: %d, output dir: %s\n", depth, outputDir)
	fmt.Printf("Original paths - CSS: %s, JS: %s\n", config.CSSPath, config.JSPath)

	// Special case for blog/category structure - explicit detection
	outComponents := strings.Split(outputDir, string(filepath.Separator))
	if len(outComponents) >= 2 && outComponents[0] == "blog" {
		// For first-level blog categories (blog/X), ALWAYS use ../../ prefix regardless of depth calculation
		if len(outComponents) == 2 {
			adjustedConfig.CSSPath = "../../css/style-blog.css"
			adjustedConfig.JSPath = "../../js/script.js"
			fmt.Printf("Blog/category structure detected, using ../../ prefix\n")
			return adjustedConfig
		} else {
			// For nested blog structure (blog/X/Y), count the components to determine proper path
			nestingDepth := len(outComponents) - 1 // blog/Nim/SubDir would be depth 2
			prefix := ""
			for i := 0; i < nestingDepth; i++ {
				prefix += "../"
			}
			adjustedConfig.CSSPath = prefix + "css/style-blog.css"
			adjustedConfig.JSPath = prefix + "js/script.js"
			fmt.Printf("Nested blog structure detected with depth %d, using prefix: %s\n", nestingDepth, prefix)
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

	// Add an additional check for known CSS and JS paths
	if !strings.HasPrefix(adjustedConfig.CSSPath, "http") && !strings.Contains(adjustedConfig.CSSPath, "css/") {
		// If CSS path doesn't contain "css/" folder, it might be misconfigured
		// Try to fix by appending to path
		if depth > 0 {
			prefix := ""
			for i := 0; i < depth; i++ {
				prefix += "../"
			}
			adjustedConfig.CSSPath = prefix + "css/style-blog.css"
		} else {
			adjustedConfig.CSSPath = "css/style-blog.css"
		}
	}

	if !strings.HasPrefix(adjustedConfig.JSPath, "http") && !strings.Contains(adjustedConfig.JSPath, "js/") {
		// If JS path doesn't contain "js/" folder, it might be misconfigured
		// Try to fix by appending to path
		if depth > 0 {
			prefix := ""
			for i := 0; i < depth; i++ {
				prefix += "../"
			}
			adjustedConfig.JSPath = prefix + "js/script.js"
		} else {
			adjustedConfig.JSPath = "js/script.js"
		}
	}

	fmt.Printf("Adjusted paths - CSS: %s, JS: %s\n", adjustedConfig.CSSPath, adjustedConfig.JSPath)
	return adjustedConfig
}

// ==========================================================================
// Markdown Parsing
// ==========================================================================

// extractMetadata parses the top of the markdown file for metadata in HTML comments
func extractMetadata(mdContent []byte) (map[string]string, []byte) {
	content := string(mdContent)
	metadata := make(map[string]string)

	// Check if the file starts with an HTML comment
	if strings.HasPrefix(content, "<!--") {
		endIndex := strings.Index(content, "-->")
		if endIndex != -1 {
			metaContent := content[4:endIndex] // Trim the <!-- and -->
			lines := strings.Split(metaContent, "\n")

			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					metadata[key] = value
				}
			}

			// Return the content without the metadata comment
			return metadata, []byte(content[endIndex+3:])
		}
	}

	return metadata, mdContent
}

// renderCodeWithSyntaxHighlighting renders a code block with syntax highlighting
func renderCodeWithSyntaxHighlighting(w *bytes.Buffer, lang string, code []byte) {
	// Get language lexer
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Use the "nord" style or any other style you prefer
	style := styles.Get("nord")
	if style == nil {
		style = styles.Fallback
	}

	// Create a formatter with line numbers
	formatter := chromahtml.New(
		chromahtml.WithClasses(true),
		chromahtml.WithLineNumbers(true),
	)

	// Tokenize the code
	iterator, err := lexer.Tokenise(nil, string(code))
	if err != nil {
		// Fallback to plain text if tokenization fails
		w.WriteString("<pre><code>")
		w.WriteString(string(code))
		w.WriteString("</code></pre>")
		return
	}

	// Format the tokens
	err = formatter.Format(w, style, iterator)
	if err != nil {
		// Fallback to plain text if formatting fails
		w.WriteString("<pre><code>")
		w.WriteString(string(code))
		w.WriteString("</code></pre>")
	}
}

// parseMarkdown converts markdown content to HTML
func parseMarkdown(content []byte) template.HTML {
	// Parse markdown
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(content)

	// Set up custom renderer to handle code blocks with syntax highlighting
	htmlFlags := mdhtml.CommonFlags | mdhtml.HrefTargetBlank

	// Custom hook to handle code blocks with syntax highlighting
	renderHook := func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if cb, ok := node.(*ast.CodeBlock); ok && entering {
			var buf bytes.Buffer
			lang := string(cb.Info)
			renderCodeWithSyntaxHighlighting(&buf, lang, cb.Literal)
			w.Write(buf.Bytes())
			return ast.GoToNext, true
		}
		return ast.GoToNext, false
	}

	opts := mdhtml.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: renderHook,
	}

	renderer := mdhtml.NewRenderer(opts)
	htmlContent := markdown.Render(doc, renderer)

	// Sanitize HTML (remove potentially unsafe content)
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("class").OnElements("pre", "code", "span", "div")
	policy.AllowAttrs("style").OnElements("pre", "code", "span", "div")
	sanitizedHTML := policy.SanitizeBytes(htmlContent)

	return template.HTML(sanitizedHTML)
}

// ==========================================================================
// Template Handling
// ==========================================================================

// loadTemplate loads either a custom template from the file system or the embedded default template
func loadTemplate(templatePath string) (*template.Template, error) {
	if templatePath != "" {
		// Use custom template from filesystem
		return template.ParseFiles(templatePath)
	}

	// Use embedded template
	return template.ParseFS(templateFS, "templates/blog.html")
}

// generateIndex creates an index.html file in the output directory
func generateIndex(outputDir string, blogPosts []BlogPost, subdirectories []Directory, folderName string, cssPath string, jsPath string, backURL string) error {
	// Sort directories alphabetically
	sort.Slice(subdirectories, func(i, j int) bool {
		return subdirectories[i].Name < subdirectories[j].Name
	})

	// Sort blog posts by date (newest first) if we have any
	if len(blogPosts) > 0 {
		sort.Slice(blogPosts, func(i, j int) bool {
			dateI, _ := time.Parse("2006-01-02", blogPosts[i].Date)
			dateJ, _ := time.Parse("2006-01-02", blogPosts[j].Date)
			return dateI.After(dateJ)
		})
	}

	// Only generate index.html if we have posts or directories to display
	if len(blogPosts) > 0 || len(subdirectories) > 0 {
		// Load index template
		indexTmpl, err := template.ParseFS(templateFS, "templates/index.html")
		if err != nil {
			return fmt.Errorf("error parsing index template: %v", err)
		}

		// Create index.html
		indexFile := filepath.Join(outputDir, "index.html")
		file, err := os.Create(indexFile)
		if err != nil {
			return fmt.Errorf("error creating index file: %v", err)
		}
		defer file.Close()

		// Execute template
		data := IndexData{
			Title:          folderName,
			CSSPath:        cssPath,
			JSPath:         jsPath,
			Posts:          blogPosts,
			Directories:    subdirectories,
			BackURL:        backURL,
			HasDirectories: len(subdirectories) > 0,
			URL:            backURL,
			Image:          "", // Default to empty string for index pages
		}

		if err := indexTmpl.Execute(file, data); err != nil {
			return fmt.Errorf("error executing index template: %v", err)
		}
	}

	return nil
}

// getFolderTitle determines the title for a directory index
func getFolderTitle(inputDir string, defaultTitle string) string {
	folderName := filepath.Base(inputDir)
	if folderName == "." || folderName == "" {
		// If at root, use the default site title
		folderName = defaultTitle
	}

	// Check for a custom title file
	customTitleFile := filepath.Join(inputDir, "index_title.txt")
	if _, err := os.Stat(customTitleFile); err == nil {
		// File exists, read its content for the title
		titleContent, err := os.ReadFile(customTitleFile)
		if err == nil && len(titleContent) > 0 {
			// Use the custom title
			folderName = strings.TrimSpace(string(titleContent))
		}
	}

	return folderName
}

// calculateBackURL determines the URL for going back to the parent directory
func calculateBackURL(depth int) string {
	fmt.Printf("Calculating back URL for depth: %d\n", depth)

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

	fmt.Printf("Generated back URL: %s\n", backURL)
	return backURL
}

// ==========================================================================
// Primary Functions
// ==========================================================================

// GenerateTemplateFile creates a default template
func GenerateTemplateFile(path string) error {
	// Read the default template from embedded files
	templateContent, err := templateFS.ReadFile("templates/blog.html")
	if err != nil {
		return fmt.Errorf("error reading default template: %v", err)
	}

	// Write the template to the specified path
	return os.WriteFile(path, templateContent, 0644)
}

// ConvertFile converts a markdown file to HTML
func ConvertFile(mdFile string, config Config, inputRoot string) (map[string]string, error) {
	// Read markdown file
	mdContent, err := os.ReadFile(mdFile)
	if err != nil {
		return nil, fmt.Errorf("error reading markdown file: %v", err)
	}

	// Extract metadata
	metadata, contentWithoutMeta := extractMetadata(mdContent)

	// Calculate the directory of the current file
	fileDir := filepath.Dir(mdFile)

	// Process wiki links
	contentWithLinks := processWikiLinks(contentWithoutMeta, fileDir, inputRoot, config.OutputDir)

	// Parse markdown to HTML
	sanitizedHTML := parseMarkdown(contentWithLinks)

	// Load template
	tmpl, err := loadTemplate(config.TemplateFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}

	// Get filename without extension
	baseName := filepath.Base(mdFile)
	fileNameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	// Keep original filename for display (preserving hyphens)
	displayTitle := fileNameWithoutExt
	// Replace spaces with hyphens only for the filename
	fileNameWithoutExt = strings.ReplaceAll(fileNameWithoutExt, " ", "-")

	// Create output file
	outputFile := filepath.Join(config.OutputDir, fileNameWithoutExt+".html")
	file, err := os.Create(outputFile)
	if err != nil {
		return nil, fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	// Prepare template data
	title := metadata["Title"]
	if title == "" {
		title = displayTitle // Use the original filename with its hyphens
	}

	author := metadata["Author"]
	if author == "" {
		author = config.DefaultAuthor
	}

	description := metadata["Description"]
	date := metadata["Date"]
	image := metadata["Image"]

	// Calculate the URL for the page
	var url string
	if config.OutputDir != "" {
		// If we have an output directory, use that as the base for the URL
		relPath, err := filepath.Rel(config.OutputDir, outputFile)
		if err != nil {
			// If we can't get a relative path, use the filename
			url = fileNameWithoutExt + ".html"
		} else {
			url = strings.ReplaceAll(relPath, string(filepath.Separator), "/")
		}
	} else {
		// If no output directory specified, just use the filename
		url = fileNameWithoutExt + ".html"
	}

	data := PageData{
		Title:       title,
		Content:     sanitizedHTML,
		CSSPath:     config.CSSPath,
		JSPath:      config.JSPath,
		Author:      author,
		Description: description,
		Date:        date,
		URL:         url,
		Image:       image,
	}

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}

	return metadata, nil
}

// ConvertDirectory converts all markdown files in a directory
func ConvertDirectory(inputDir string, config Config) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}

	// Get the top-level input directory (for absolute path calculations)
	inputRoot := inputDir
	cwd, cwdErr := os.Getwd()

	if config.Recursive && cwdErr == nil {
		// When recursive is enabled, use the current working directory as root
		inputRoot = cwd
	} else if config.Recursive {
		// If we can't get CWD, try to use absolute path
		absInputDir, err := filepath.Abs(inputDir)
		if err == nil {
			inputRoot = filepath.Dir(absInputDir)
		}
	}

	// Calculate directory depth
	currentDepth, err := calculatePathDepth(inputDir, inputRoot)
	if err != nil {
		// If we can't determine depth, default to 0
		currentDepth = 0
	}

	// Special handling for blog paths
	outComponents := strings.Split(config.OutputDir, string(filepath.Separator))
	if len(outComponents) >= 2 && outComponents[0] == "blog" {
		// Force blog structure to use the correct depth
		if len(outComponents) == 2 {
			// First level category (blog/X)
			fmt.Printf("Blog category detected: %s\n", config.OutputDir)

			// Use explicit depth of 2 levels for CSS/JS paths (../../)
			currentConfig := config
			currentConfig.CSSPath = "../../css/style-blog.css"
			currentConfig.JSPath = "../../js/script.js"

			// Process files with blog-specific config
			return processFiles(inputDir, inputRoot, currentConfig, 1) // 1 for blog category depth
		} else {
			// Nested blog structure (blog/X/Y/...)
			nestingDepth := len(outComponents) - 1 // Count depth starting from blog
			fmt.Printf("Nested blog structure detected: %s (depth: %d)\n",
				config.OutputDir, nestingDepth)

			// Calculate proper path prefix
			prefix := ""
			for i := 0; i < nestingDepth; i++ {
				prefix += "../"
			}

			// Apply blog-specific paths with proper nesting
			currentConfig := config
			currentConfig.CSSPath = prefix + "css/style-blog.css"
			currentConfig.JSPath = prefix + "js/script.js"

			return processFiles(inputDir, inputRoot, currentConfig, nestingDepth)
		}
	}

	// For non-blog folders, use standard path adjustment
	currentConfig := adjustPaths(config, currentDepth, config.OutputDir)
	return processFiles(inputDir, inputRoot, currentConfig, currentDepth)
}

// processFiles processes markdown files in a directory and handles subdirectories
func processFiles(inputDir string, inputRoot string, config Config, depth int) error {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("error reading input directory: %v", err)
	}

	var blogPosts []BlogPost
	hasIndexMd := false

	// First pass: Check if there's an index.md file
	for _, file := range files {
		if !file.IsDir() && (file.Name() == "index.md" || file.Name() == "index.markdown") {
			hasIndexMd = true
			break
		}
	}

	// Second pass: Process files and collect blog posts
	for _, file := range files {
		filePath := filepath.Join(inputDir, file.Name())

		if file.IsDir() && config.Recursive {
			// For subdirectories, create corresponding output directory
			relPath, err := filepath.Rel(inputDir, filePath)
			if err != nil {
				return fmt.Errorf("error calculating relative path: %v", err)
			}

			subOutputDir := filepath.Join(config.OutputDir, relPath)

			// Create a new config for the subdirectory
			subConfig := config
			subConfig.OutputDir = subOutputDir

			// Process the subdirectory recursively
			if err := ConvertDirectory(filePath, subConfig); err != nil {
				return err
			}
		} else if !file.IsDir() && (strings.HasSuffix(file.Name(), ".md") || strings.HasSuffix(file.Name(), ".markdown")) {
			metadata, err := ConvertFile(filePath, config, inputRoot)
			if err != nil {
				return err
			}

			// Skip index.md for the blog post list
			if file.Name() == "index.md" || file.Name() == "index.markdown" {
				continue
			}

			// Add to blog posts list
			fileNameWithoutExt := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			// Keep original filename for display (preserving hyphens)
			displayTitle := fileNameWithoutExt
			// Replace spaces with hyphens only for the filename
			fileNameWithoutExt = strings.ReplaceAll(fileNameWithoutExt, " ", "-")
			title := metadata["Title"]
			if title == "" {
				title = displayTitle // Use the original filename with its hyphens
			}

			blogPosts = append(blogPosts, BlogPost{
				Title:       title,
				Link:        fileNameWithoutExt + ".html",
				Description: metadata["Description"],
				Date:        metadata["Date"],
			})
		}
	}

	// Generate index.html either if requested via config or if there's no index.md
	if config.GenerateList || !hasIndexMd {
		// Create a list of subdirectories
		var subdirectories []Directory
		for _, file := range files {
			if file.IsDir() && file.Name() != "." && file.Name() != ".." {
				// Skip hidden directories
				if strings.HasPrefix(file.Name(), ".") {
					continue
				}
				subdirectories = append(subdirectories, Directory{
					Name: file.Name(),
					Link: fmt.Sprintf("%s/index.html", file.Name()),
				})
			}
		}

		// Get folder name for the title
		folderName := getFolderTitle(inputDir, config.SiteTitle)

		// Calculate appropriate back URL based on depth
		backURL := calculateBackURL(depth)

		// Generate the index file
		if err := generateIndex(
			config.OutputDir,
			blogPosts,
			subdirectories,
			folderName,
			config.CSSPath,
			config.JSPath,
			backURL,
		); err != nil {
			return err
		}
	}

	return nil
}
