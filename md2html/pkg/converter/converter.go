package converter

import (
	"bytes"
	"embed"
	"encoding/xml"
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
	GenerateRSS   bool
	SiteURL       string
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
	FullURL     string // Full URL for RSS feed
	Author      string // Author for RSS feed
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

// LandingSection represents a terminal section in the landing page
type LandingSection struct {
	Command string        // e.g., "cat about.txt"
	Content template.HTML // Rendered markdown content
}

// LandingLink represents a link in the landing page "ls" output
type LandingLink struct {
	Name     string // Display name, e.g., "blog.md"
	URL      string // Link URL
	External bool   // Whether to open in new tab
}

// LandingData represents data for the landing page template
type LandingData struct {
	Title       string
	Description string
	CSSPath     string
	JSPath      string
	Sections    []LandingSection
	Links       []LandingLink
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
			adjustedConfig.CSSPath = "../../css/style-blog.css"
			adjustedConfig.JSPath = "../../js/script.js"
			return adjustedConfig
		} else if len(blogComponents) > 2 {
			// For nested blog structure (blog/X/Y), count the components to determine proper path
			nestingDepth := len(blogComponents) - 1 // blog/Nim/SubDir would be depth 2
			prefix := ""
			for i := 0; i < nestingDepth; i++ {
				prefix += "../"
			}
			adjustedConfig.CSSPath = prefix + "css/style-blog.css"
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

// extractLandingSections parses markdown content into landing page sections
// Sections are delimited by <!-- Section: command --> comments
func extractLandingSections(content []byte) []LandingSection {
	var sections []LandingSection
	contentStr := string(content)

	// Regular expression to match section delimiters
	// Format: <!-- Section: command -->
	sectionRe := regexp.MustCompile(`<!--\s*Section:\s*(.+?)\s*-->`)

	// Find all section markers
	matches := sectionRe.FindAllStringSubmatchIndex(contentStr, -1)

	if len(matches) == 0 {
		return sections
	}

	for i, match := range matches {
		// match[0] and match[1] are the full match boundaries
		// match[2] and match[3] are the command capture group boundaries
		command := contentStr[match[2]:match[3]]

		// Determine the content boundaries
		contentStart := match[1] // End of the current marker
		var contentEnd int
		if i+1 < len(matches) {
			contentEnd = matches[i+1][0] // Start of the next marker
		} else {
			// For the last section, find <!-- Links --> or end of content
			linksIdx := strings.Index(contentStr[contentStart:], "<!-- Links -->")
			if linksIdx != -1 {
				contentEnd = contentStart + linksIdx
			} else {
				contentEnd = len(contentStr)
			}
		}

		// Extract and trim the section content
		sectionContent := strings.TrimSpace(contentStr[contentStart:contentEnd])

		if sectionContent != "" {
			// Parse the markdown content to HTML
			htmlContent := parseMarkdown([]byte(sectionContent))

			sections = append(sections, LandingSection{
				Command: command,
				Content: htmlContent,
			})
		}
	}

	return sections
}

// extractLandingLinks parses the <!-- Links --> section from markdown content
// Format: <!-- Links --> followed by markdown links like - [name](url)
func extractLandingLinks(content []byte) []LandingLink {
	var links []LandingLink
	contentStr := string(content)

	// Find the <!-- Links --> marker
	linksIdx := strings.Index(contentStr, "<!-- Links -->")
	if linksIdx == -1 {
		return links
	}

	// Get content after the marker
	linksContent := contentStr[linksIdx+len("<!-- Links -->"):]

	// Regular expression to match markdown links: [name](url)
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

	matches := linkRe.FindAllStringSubmatch(linksContent, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			name := match[1]
			url := match[2]

			// Determine if external (starts with http:// or https://)
			external := strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")

			links = append(links, LandingLink{
				Name:     name,
				URL:      url,
				External: external,
			})
		}
	}

	return links
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
		// Ensure the output directory exists
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %v", err)
		}
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

// ==========================================================================
// RSS Feed Generation
// ==========================================================================

// RSSFeed represents an RSS 2.0 feed structure
type RSSFeed struct {
	XMLName   xml.Name   `xml:"rss"`
	Version   string     `xml:"version,attr"`
	XMLNSAtom string     `xml:"xmlns:atom,attr,omitempty"`
	Channel   RSSChannel `xml:"channel"`
}

// AtomLink represents an atom:link element
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

// RSSChannel represents the channel element in RSS
type RSSChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	AtomLink      *AtomLink `xml:"atom:link,omitempty"`
	Items         []RSSItem `xml:"item"`
}

// RSSItem represents an item in the RSS feed
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Author      string `xml:"author,omitempty"`
	GUID        string `xml:"guid"`
}

// formatAuthorAsEmail formats an author name as an RSS 2.0 compliant email address
// RSS 2.0 requires author field to be in format "email@example.com (Name)" or "email@example.com"
func formatAuthorAsEmail(authorName string, siteURL string) string {
	if authorName == "" {
		return ""
	}

	// Extract domain from site URL
	domain := "example.com"
	if siteURL != "" {
		// Remove protocol
		url := strings.TrimPrefix(siteURL, "http://")
		url = strings.TrimPrefix(url, "https://")
		// Remove path
		if idx := strings.Index(url, "/"); idx != -1 {
			url = url[:idx]
		}
		if url != "" {
			domain = url
		}
	}

	// Create email from author name (lowercase, replace spaces with dots)
	emailName := strings.ToLower(strings.ReplaceAll(authorName, " ", "."))
	email := fmt.Sprintf("%s@%s", emailName, domain)

	// Return in RSS 2.0 format: "email@domain.com (Name)"
	return fmt.Sprintf("%s (%s)", email, authorName)
}

// generateRSSFeed creates an RSS 2.0 feed from blog posts
func generateRSSFeed(blogPosts []BlogPost, config Config, outputDir string) error {
	if !config.GenerateRSS || len(blogPosts) == 0 {
		return nil
	}

	// Sort posts by date (newest first)
	sortedPosts := make([]BlogPost, len(blogPosts))
	copy(sortedPosts, blogPosts)
	sort.Slice(sortedPosts, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", sortedPosts[i].Date)
		dateJ, errJ := time.Parse("2006-01-02", sortedPosts[j].Date)
		if errI != nil || errJ != nil {
			// If dates can't be parsed, maintain original order
			return false
		}
		return dateI.After(dateJ)
	})

	// Build RSS items
	var items []RSSItem
	for _, post := range sortedPosts {
		// Parse date for RSS format
		var pubDate string
		if post.Date != "" {
			if t, err := time.Parse("2006-01-02", post.Date); err == nil {
				pubDate = t.Format(time.RFC1123Z)
			}
		} else {
			pubDate = time.Now().Format(time.RFC1123Z)
		}

		// Build full URL
		fullURL := post.FullURL
		if fullURL == "" {
			// Fallback to relative link if FullURL not set
			fullURL = post.Link
		}

		// Make it absolute if site URL is provided
		if config.SiteURL != "" && !strings.HasPrefix(fullURL, "http") {
			// Ensure site URL doesn't end with / and link doesn't start with /
			siteURL := strings.TrimSuffix(config.SiteURL, "/")
			if strings.HasPrefix(fullURL, "/") {
				fullURL = siteURL + fullURL
			} else {
				fullURL = siteURL + "/" + fullURL
			}
		}

		// Use FullURL as GUID if available, otherwise use link
		guid := fullURL
		if guid == "" {
			guid = post.Link
		}

		item := RSSItem{
			Title:       post.Title,
			Link:        fullURL,
			Description: post.Description,
			PubDate:     pubDate,
			GUID:        guid,
		}

		// Add author if available (format as email for RSS 2.0 compliance)
		authorName := post.Author
		if authorName == "" {
			authorName = config.DefaultAuthor
		}
		if authorName != "" {
			item.Author = formatAuthorAsEmail(authorName, config.SiteURL)
		}

		items = append(items, item)
	}

	// Calculate feed URL for atom:link self reference
	var atomLink *AtomLink
	if config.SiteURL != "" {
		// Determine the relative path of the feed from site root
		// For blog root: /blog/feed.xml
		// For blog/category: /blog/category/feed.xml
		feedURL := config.SiteURL

		// Find "blog" in the output path and extract everything after it
		outComponents := strings.Split(outputDir, string(filepath.Separator))
		blogIdx := -1
		for i, comp := range outComponents {
			if comp == "blog" {
				blogIdx = i
				break
			}
		}

		if blogIdx >= 0 {
			// Get path components after "blog"
			afterBlog := outComponents[blogIdx+1:]
			if len(afterBlog) == 0 {
				// At blog root
				feedURL = strings.TrimSuffix(feedURL, "/") + "/blog/feed.xml"
			} else {
				// In a blog subdirectory
				subPath := strings.Join(afterBlog, "/")
				feedURL = strings.TrimSuffix(feedURL, "/") + "/blog/" + subPath + "/feed.xml"
			}
		} else {
			// Fallback: just append feed.xml to the output dir
			feedURL = strings.TrimSuffix(feedURL, "/") + "/" + strings.ReplaceAll(outputDir, string(filepath.Separator), "/") + "/feed.xml"
		}

		atomLink = &AtomLink{
			Href: feedURL,
			Rel:  "self",
			Type: "application/rss+xml",
		}
	}

	// Build RSS feed
	xmlnsAtom := ""
	if atomLink != nil {
		xmlnsAtom = "http://www.w3.org/2005/Atom"
	}
	feed := RSSFeed{
		Version:   "2.0",
		XMLNSAtom: xmlnsAtom,
		Channel: RSSChannel{
			Title:         config.SiteTitle,
			Link:          config.SiteURL,
			Description:   fmt.Sprintf("Blog posts from %s", config.SiteTitle),
			Language:      "en-us",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			AtomLink:      atomLink,
			Items:         items,
		},
	}

	// Generate XML
	xmlData, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling RSS feed: %v", err)
	}

	// Add XML header
	xmlOutput := []byte(xml.Header + string(xmlData))

	// Write RSS feed file
	rssFile := filepath.Join(outputDir, "feed.xml")
	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}
	if err := os.WriteFile(rssFile, xmlOutput, 0644); err != nil {
		return fmt.Errorf("error writing RSS feed file: %v", err)
	}

	return nil
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

// ConvertLandingPage converts a markdown file to a landing page HTML
func ConvertLandingPage(mdFile string, config Config) error {
	// Read markdown file
	mdContent, err := os.ReadFile(mdFile)
	if err != nil {
		return fmt.Errorf("error reading markdown file: %v", err)
	}

	// Extract metadata
	metadata, contentWithoutMeta := extractMetadata(mdContent)

	// Extract sections and links
	sections := extractLandingSections(contentWithoutMeta)
	links := extractLandingLinks(contentWithoutMeta)

	// Load landing template
	tmpl, err := template.ParseFS(templateFS, "templates/landing.html")
	if err != nil {
		return fmt.Errorf("error parsing landing template: %v", err)
	}

	// Create output file (index.html at output root)
	outputFile := filepath.Join(config.OutputDir, "index.html")
	// Ensure the output directory exists
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	// Prepare template data
	title := metadata["Title"]
	if title == "" {
		title = "Home"
	}

	description := metadata["Description"]

	data := LandingData{
		Title:       title,
		Description: description,
		CSSPath:     config.CSSPath,
		JSPath:      config.JSPath,
		Sections:    sections,
		Links:       links,
	}

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing landing template: %v", err)
	}

	fmt.Printf("Generated landing page: %s\n", outputFile)
	return nil
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

	// Check if this is a landing page template
	if metadata["Template"] == "landing" {
		err := ConvertLandingPage(mdFile, config)
		if err != nil {
			return nil, err
		}
		return metadata, nil
	}

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
	// Ensure the output directory exists
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return nil, fmt.Errorf("error creating output directory: %v", err)
	}
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

	// Check for landing page (index.md with Template: landing) at input root
	landingFile := filepath.Join(inputDir, "index.md")
	if _, err := os.Stat(landingFile); err == nil {
		// Read and check metadata
		mdContent, err := os.ReadFile(landingFile)
		if err == nil {
			metadata, _ := extractMetadata(mdContent)
			if metadata["Template"] == "landing" {
				// Process landing page with site root CSS
				landingConfig := config
				landingConfig.CSSPath = "css/style.css"
				landingConfig.JSPath = "js/script.js"

				if err := ConvertLandingPage(landingFile, landingConfig); err != nil {
					return fmt.Errorf("error converting landing page: %v", err)
				}
			}
		}
	}

	// Calculate directory depth
	currentDepth, err := calculatePathDepth(inputDir, inputRoot)
	if err != nil {
		// If we can't determine depth, default to 0
		currentDepth = 0
	}

	// Special handling for blog paths
	// Find blog in the output path components (may be prefixed with dist/ or similar)
	outComponents := strings.Split(config.OutputDir, string(filepath.Separator))
	blogIndex := -1
	for i, comp := range outComponents {
		if comp == "blog" {
			blogIndex = i
			break
		}
	}

	// Calculate blog-relative components (everything after "blog")
	var blogComponents []string
	if blogIndex >= 0 {
		blogComponents = outComponents[blogIndex:]
	}
	isBlogRoot := len(blogComponents) == 1 && blogIndex >= 0

	if blogIndex >= 0 {
		// Force blog structure to use the correct depth
		if len(blogComponents) == 2 {
			// First level category (blog/X)
			// Use explicit depth of 2 levels for CSS/JS paths (../../)
			currentConfig := config
			currentConfig.CSSPath = "../../css/style-blog.css"
			currentConfig.JSPath = "../../js/script.js"
			// Enable recursive processing for RSS to collect all posts from this category
			if config.GenerateRSS {
				currentConfig.Recursive = true
			}

			// Process files with blog-specific config
			blogPosts, err := processFiles(inputDir, inputRoot, currentConfig, 1) // 1 for blog category depth
			if err != nil {
				return err
			}

			// Generate RSS feed for this category if enabled
			if config.GenerateRSS {
				// Use category name for the feed title
				categoryName := blogComponents[1]
				categoryConfig := config
				categoryConfig.SiteTitle = fmt.Sprintf("%s - %s", config.SiteTitle, categoryName)
				if err := generateRSSFeed(blogPosts, categoryConfig, config.OutputDir); err != nil {
					return fmt.Errorf("error generating RSS feed for category %s: %v", categoryName, err)
				}
			}

			return nil
		} else if isBlogRoot {
			// Blog root directory
			// Use explicit depth of 1 level for CSS/JS paths (../)
			currentConfig := config
			currentConfig.CSSPath = "../css/style-blog.css"
			currentConfig.JSPath = "../js/script.js"
			// Enable recursive processing for RSS to collect all posts from subdirectories
			if config.GenerateRSS {
				currentConfig.Recursive = true
			}

			// Process files with blog-specific config
			blogPosts, err := processFiles(inputDir, inputRoot, currentConfig, 0)
			if err != nil {
				return err
			}

			// Generate RSS feed at blog root if enabled
			if config.GenerateRSS {
				if err := generateRSSFeed(blogPosts, config, config.OutputDir); err != nil {
					return fmt.Errorf("error generating RSS feed: %v", err)
				}
			}

			return nil
		} else {
			// Nested blog structure (blog/X/Y/...)
			nestingDepth := len(blogComponents) - 1 // Count depth starting from blog

			// Calculate proper path prefix
			prefix := ""
			for i := 0; i < nestingDepth; i++ {
				prefix += "../"
			}

			// Apply blog-specific paths with proper nesting
			currentConfig := config
			currentConfig.CSSPath = prefix + "css/style-blog.css"
			currentConfig.JSPath = prefix + "js/script.js"
			// Enable recursive processing for RSS to collect all posts from this subdirectory
			if config.GenerateRSS {
				currentConfig.Recursive = true
			}

			blogPosts, err := processFiles(inputDir, inputRoot, currentConfig, nestingDepth)
			if err != nil {
				return err
			}

			// Generate RSS feed for this nested directory if enabled
			if config.GenerateRSS {
				// Use directory path for the feed title
				dirName := blogComponents[len(blogComponents)-1]
				categoryConfig := config
				categoryConfig.SiteTitle = fmt.Sprintf("%s - %s", config.SiteTitle, dirName)
				if err := generateRSSFeed(blogPosts, categoryConfig, config.OutputDir); err != nil {
					return fmt.Errorf("error generating RSS feed for directory %s: %v", dirName, err)
				}
			}

			return nil
		}
	}

	// For non-blog folders, use standard path adjustment
	currentConfig := adjustPaths(config, currentDepth, config.OutputDir)
	_, err = processFiles(inputDir, inputRoot, currentConfig, currentDepth)
	return err
}

// processFiles processes markdown files in a directory and handles subdirectories
// Returns collected blog posts for RSS feed generation
func processFiles(inputDir string, inputRoot string, config Config, depth int) ([]BlogPost, error) {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, fmt.Errorf("error reading input directory: %v", err)
	}

	var blogPosts []BlogPost
	hasIndexMd := false
	hasLandingPage := false

	// First pass: Check if there's an index.md file (and if it's a landing page)
	for _, file := range files {
		if !file.IsDir() && (file.Name() == "index.md" || file.Name() == "index.markdown") {
			hasIndexMd = true
			// Check if it's a landing page template
			filePath := filepath.Join(inputDir, file.Name())
			mdContent, err := os.ReadFile(filePath)
			if err == nil {
				metadata, _ := extractMetadata(mdContent)
				if metadata["Template"] == "landing" {
					hasLandingPage = true
				}
			}
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
				return nil, fmt.Errorf("error calculating relative path: %v", err)
			}

			subOutputDir := filepath.Join(config.OutputDir, relPath)

			// Create a new config for the subdirectory
			subConfig := config
			subConfig.OutputDir = subOutputDir

			// Adjust CSS/JS paths for the subdirectory depth
			// Find blog in the output path components (may be prefixed with dist/ or similar)
			subOutComponents := strings.Split(subOutputDir, string(filepath.Separator))
			subBlogIndex := -1
			for i, comp := range subOutComponents {
				if comp == "blog" {
					subBlogIndex = i
					break
				}
			}

			if subBlogIndex >= 0 {
				// For blog structure, calculate depth from site root
				// blog/ = 1 level up, blog/X/ = 2 levels up, blog/X/Y/ = 3 levels up, etc.
				blogDepth := len(subOutComponents) - subBlogIndex
				prefix := ""
				for i := 0; i < blogDepth; i++ {
					prefix += "../"
				}
				subConfig.CSSPath = prefix + "css/style-blog.css"
				subConfig.JSPath = prefix + "js/script.js"
			} else {
				// For non-blog structure, use standard depth calculation
				subConfig = adjustPaths(subConfig, depth+1, subOutputDir)
			}

			// Process the subdirectory recursively and collect its blog posts
			subPosts, err := processFiles(filePath, inputRoot, subConfig, depth+1)
			if err != nil {
				return nil, err
			}

			// Generate RSS feed for blog folders if enabled
			if config.GenerateRSS && len(subPosts) > 0 {
				// Check if this is a blog folder (blog or blog/X)
				// Find blog index in path components
				blogIdx := -1
				for i, comp := range subOutComponents {
					if comp == "blog" {
						blogIdx = i
						break
					}
				}
				if blogIdx >= 0 {
					blogRelComponents := subOutComponents[blogIdx:]
					if len(blogRelComponents) == 1 {
						// This is the blog root folder (blog/)
						// Generate RSS feed with all posts from all categories
						if err := generateRSSFeed(subPosts, config, subOutputDir); err != nil {
							// Log error but don't fail the entire process
							fmt.Fprintf(os.Stderr, "Warning: error generating RSS feed for blog: %v\n", err)
						}
					} else if len(blogRelComponents) == 2 {
						// This is a first-level category folder (blog/X)
						categoryName := blogRelComponents[1]
						categoryConfig := config
						categoryConfig.SiteTitle = fmt.Sprintf("%s - %s", config.SiteTitle, categoryName)
						if err := generateRSSFeed(subPosts, categoryConfig, subOutputDir); err != nil {
							// Log error but don't fail the entire process
							fmt.Fprintf(os.Stderr, "Warning: error generating RSS feed for category %s: %v\n", categoryName, err)
						}
					}
				}
			}

			// Add subdirectory posts to our collection
			blogPosts = append(blogPosts, subPosts...)
		} else if !file.IsDir() && (strings.HasSuffix(file.Name(), ".md") || strings.HasSuffix(file.Name(), ".markdown")) {
			// Skip landing page files (they're processed separately in ConvertDirectory)
			if file.Name() == "index.md" || file.Name() == "index.markdown" {
				// Check if it's a landing page template
				mdContent, err := os.ReadFile(filePath)
				if err == nil {
					metadata, _ := extractMetadata(mdContent)
					if metadata["Template"] == "landing" {
						// Skip - already processed by ConvertDirectory
						continue
					}
				}
			}

			metadata, err := ConvertFile(filePath, config, inputRoot)
			if err != nil {
				return nil, err
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

			// Build full URL for RSS (relative to site root)
			// Calculate path relative to blog root directory, then prepend "blog/"
			outputFile := filepath.Join(config.OutputDir, fileNameWithoutExt+".html")
			// Find blog root (directory named "blog")
			blogRoot := config.OutputDir
			foundBlogRoot := false
			for {
				base := filepath.Base(blogRoot)
				if base == "blog" {
					foundBlogRoot = true
					break
				}
				parent := filepath.Dir(blogRoot)
				if parent == blogRoot || parent == "." || parent == "/" {
					break
				}
				blogRoot = parent
			}

			var relPath string
			if foundBlogRoot {
				// Calculate relative path from blog root
				relPath, err = filepath.Rel(blogRoot, outputFile)
				if err != nil {
					relPath = fileNameWithoutExt + ".html"
				}
				// Convert to web path (forward slashes) and prepend "blog/"
				relPath = strings.ReplaceAll(relPath, string(filepath.Separator), "/")
				relPath = "blog/" + relPath
			} else {
				// Fallback: use output directory structure
				relPath = strings.ReplaceAll(outputFile, string(filepath.Separator), "/")
			}

			// Get author from metadata or use default
			author := metadata["Author"]
			if author == "" {
				author = config.DefaultAuthor
			}

			blogPosts = append(blogPosts, BlogPost{
				Title:       title,
				Link:        fileNameWithoutExt + ".html",
				Description: metadata["Description"],
				Date:        metadata["Date"],
				FullURL:     relPath,
				Author:      author,
			})
		}
	}

	// Generate index.html either if requested via config or if there's no index.md
	// But skip if there's a landing page (it's handled separately)
	if !hasLandingPage && (config.GenerateList || !hasIndexMd) {
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

		// Calculate appropriate back URL based on blog structure
		// Find blog in the output path to calculate correct relative path
		outComponents := strings.Split(config.OutputDir, string(filepath.Separator))
		blogIdx := -1
		for i, comp := range outComponents {
			if comp == "blog" {
				blogIdx = i
				break
			}
		}

		var backURL string
		if blogIdx >= 0 {
			// Calculate back URL relative to blog structure
			blogDepth := len(outComponents) - blogIdx - 1 // depth within blog (0 for blog/, 1 for blog/X/, etc.)
			if blogDepth == 0 {
				// At blog root, go back to site root
				backURL = "../index.html"
			} else {
				// In a blog category, go back to blog root
				backURL = "../index.html"
			}
		} else {
			// Not in blog, use standard calculation
			backURL = calculateBackURL(depth)
		}

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
			return nil, err
		}
	}

	return blogPosts, nil
}
