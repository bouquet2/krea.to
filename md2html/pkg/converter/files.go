package converter

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// CopyRobotsTxt copies and processes robots.txt template
func CopyRobotsTxt(outputDir string, siteURL string) error {
	robotsTemplatePath := "md2html/pkg/converter/templates/robots.txt"
	robotsContent, err := os.ReadFile(robotsTemplatePath)
	if err != nil {
		return fmt.Errorf("error reading robots.txt template: %v", err)
	}

	// Replace {{.URL}} with actual site URL
	robotsProcessed := strings.ReplaceAll(string(robotsContent), "{{.URL}}", siteURL)

	robotsOutputPath := filepath.Join(outputDir, "robots.txt")
	if err := os.WriteFile(robotsOutputPath, []byte(robotsProcessed), 0644); err != nil {
		return fmt.Errorf("error writing robots.txt: %v", err)
	}

	log.Debug().Str("robots_file", robotsOutputPath).Str("site_url", siteURL).Msg("robots.txt generated successfully")
	fmt.Printf("Generated robots.txt: %s\n", robotsOutputPath)
	return nil
}

// parseSettings parses a comma-separated list of settings into a map
// Example: "hide-topbar, fullscreen" -> map[string]bool{"hide-topbar": true, "fullscreen": true}
func parseSettings(settingsStr string) map[string]bool {
	settings := make(map[string]bool)
	if settingsStr == "" {
		return settings
	}

	parts := strings.Split(settingsStr, ",")
	for _, part := range parts {
		setting := strings.TrimSpace(strings.ToLower(part))
		if setting != "" {
			settings[setting] = true
		}
	}
	return settings
}

// ConvertLandingPage converts a markdown file to a landing page HTML
func ConvertLandingPage(mdFile string, config Config) error {
	logger := log.With().Str("file", mdFile).Str("type", "landing_page").Logger()
	logger.Debug().Msg("Starting conversion of landing page")

	// Read markdown file
	mdContent, err := os.ReadFile(mdFile)
	if err != nil {
		return fmt.Errorf("error reading markdown file: %v", err)
	}
	logger.Debug().Int("size_bytes", len(mdContent)).Msg("Landing page markdown file read successfully")

	// Extract metadata
	metadata, contentWithoutMeta := extractMetadata(mdContent)

	// Extract sections and links
	sections := extractLandingSections(contentWithoutMeta)
	links := extractLandingLinks(contentWithoutMeta)

	// Fetch recent blog posts from Git history
	recentPosts := fetchRecentPosts(filepath.Dir(mdFile), config, 5)

	// Load landing template
	tmpl, err := template.ParseFS(templateFS, "templates/landing.html")
	if err != nil {
		return fmt.Errorf("error parsing landing template: %v", err)
	}

	// Create output file (index.html at output root)
	outputFile := filepath.Join(config.OutputDir, "index.html")
	logger.Debug().Str("output_file", outputFile).Msg("Creating landing page HTML file")
	// Ensure the output directory exists
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %v", err)
	}
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()
	logger.Debug().Msg("Landing page output file created successfully")

	// Prepare template data
	title := metadata["Title"]
	if title == "" {
		title = "Home"
	}

	description := metadata["Description"]

	// Parse settings from metadata (comma-separated list)
	settings := parseSettings(metadata["Settings"])

	data := LandingData{
		Title:       title,
		Description: description,
		CSSPath:     config.CSSPath,
		JSPath:      config.JSPath,
		Sections:    sections,
		Links:       links,
		Settings:    settings,
		Posts:       recentPosts,
		DefaultTheme: func() string {
			if config.DefaultTheme != "" {
				return config.DefaultTheme
			}
			return "nord" // Fallback to nord if no theme specified
		}(),
		URL: config.SiteURL,
	}

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing landing template: %v", err)
	}

	logger.Debug().Str("output_file", outputFile).Msg("Landing page successfully generated")
	fmt.Printf("Generated landing page: %s\n", outputFile)
	return nil
}

// ConvertFile converts a markdown file to HTML
func ConvertFile(mdFile string, config Config, inputRoot string) (map[string]string, error) {
	logger := log.With().Str("file", mdFile).Logger()
	logger.Debug().Msg("Starting conversion of markdown file")

	// Read markdown file
	mdContent, err := os.ReadFile(mdFile)
	if err != nil {
		return nil, fmt.Errorf("error reading markdown file: %v", err)
	}
	logger.Debug().Int("size_bytes", len(mdContent)).Msg("Markdown file read successfully")

	// Extract metadata
	metadata, contentWithoutMeta := extractMetadata(mdContent)
	logger.Debug().Str("title", metadata["Title"]).Str("author", metadata["Author"]).Msg("Metadata extracted")

	// Check if this is a landing page template
	if metadata["Template"] == "landing" {
		logger.Debug().Msg("File is landing page template, processing as landing page")
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
	logger.Debug().Str("output_file", outputFile).Msg("Creating output HTML file")
	// Ensure the output directory exists
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return nil, fmt.Errorf("error creating output directory: %v", err)
	}
	file, err := os.Create(outputFile)
	if err != nil {
		return nil, fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()
	logger.Debug().Msg("Output file created successfully")

	// Prepare template data
	title := metadata["Title"]
	if title == "" {
		title = displayTitle // Use the original filename with its hyphens
	}

	author := metadata["Author"]
	description := metadata["Description"]
	date := metadata["Date"]
	image := metadata["Image"]

	// Extract and normalize tags
	tags := normalizeTags(metadata["Tags"])

	// Get Git information for commit tracking
	var commitHash, commitDate, commitAuthor, commitURL string
	if config.ShowCommitInfo {
		gitInfo, err := GetFileGitInfo(mdFile)
		if err == nil {
			commitHash = gitInfo.CommitHash
			commitDate = gitInfo.LastModified.Format("2006-01-02 15:04:05")
			commitAuthor = gitInfo.Author
			commitURL = GetCommitURL(config.GitWebURL, commitHash)
			logger.Debug().
				Str("commit_hash", commitHash[:8]).
				Str("commit_date", commitDate).
				Str("commit_author", commitAuthor).
				Msg("Collected Git commit information")
		} else {
			logger.Debug().Err(err).Msg("Could not get Git information for file")
		}
	}

	// If no date or author in metadata, try to get from Git history
	if date == "" || author == "" {
		if commitHash != "" {
			if date == "" {
				date = commitDate[:10]  // Use just the date part
				metadata["Date"] = date // Store in metadata for later use
				logger.Debug().Str("git_date", date).Msg("Using Git commit date for post")
			}
			if author == "" {
				author = commitAuthor
				metadata["Author"] = author // Store in metadata for later use
				logger.Debug().Str("git_author", author).Msg("Using Git commit author for post")
			}
		}
	}

	// Fall back to default author if still empty
	if author == "" {
		author = config.DefaultAuthor
	}

	// Calculate the URL for the page (full URL for og:url and JSON-LD)
	var url string
	// Find the dist root (parent of blog, etc.)
	distRoot := config.OutputDir
	for {
		base := filepath.Base(distRoot)
		parent := filepath.Dir(distRoot)
		if base == "dist" || parent == distRoot || parent == "." || parent == "/" {
			break
		}
		distRoot = parent
	}
	// Calculate relative path from dist root
	relPath, err := filepath.Rel(distRoot, outputFile)
	if err != nil {
		relPath = fileNameWithoutExt + ".html"
	}
	relPath = strings.ReplaceAll(relPath, string(filepath.Separator), "/")
	// Build full URL with site URL
	if config.SiteURL != "" {
		url = strings.TrimSuffix(config.SiteURL, "/") + "/" + relPath
	} else {
		url = relPath
	}

	// Calculate reading time from the content
	readTime := calculateReadTime(contentWithoutMeta)

	data := PageData{
		Title:        title,
		Content:      sanitizedHTML,
		CSSPath:      config.CSSPath,
		JSPath:       config.JSPath,
		Author:       author,
		Description:  description,
		Date:         date,
		URL:          url,
		Image:        image,
		CommitHash:   commitHash,
		CommitDate:   commitDate,
		CommitAuthor: commitAuthor,
		CommitURL:    commitURL,
		DefaultTheme: func() string {
			if config.DefaultTheme != "" {
				return config.DefaultTheme
			}
			return "nord" // Fallback to nord if no theme specified
		}(),
		ReadTime:     readTime,
		Tags:         tags,
		RelatedPosts: nil, // Related posts will be calculated and added in a second pass
	}

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}

	logger.Debug().Str("output_file", outputFile).Msg("Markdown file successfully converted to HTML")
	return metadata, nil
}

// ConvertDirectory converts all markdown files in a directory
func ConvertDirectory(inputDir string, config Config) error {
	logger := log.With().Str("input_dir", inputDir).Str("output_dir", config.OutputDir).Logger()
	logger.Debug().Msg("Starting directory conversion")
	logger.Debug().Bool("recursive", config.Recursive).Bool("generate_rss", config.GenerateRSS).Msg("Conversion config")
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
				// Use configured paths without modification at root level
				landingConfig.CSSPath = config.CSSPath
				landingConfig.JSPath = config.JSPath

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
			currentConfig.CSSPath = prependPathPrefix(config.CSSPath, 2)
			currentConfig.JSPath = prependPathPrefix(config.JSPath, 2)
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
			currentConfig.CSSPath = prependPathPrefix(config.CSSPath, 1)
			currentConfig.JSPath = prependPathPrefix(config.JSPath, 1)
			// Enable recursive processing to collect all posts from subdirectories for RSS and tags
			currentConfig.Recursive = true

			// Process files with blog-specific config
			blogPosts, err := processFiles(inputDir, inputRoot, currentConfig, 0)
			if err != nil {
				return err
			}

			// Generate tag pages at blog root
			if len(blogPosts) > 0 {
				tagMap := collectTags(blogPosts)

				// Generate individual tag pages
				for tagName, tagInfo := range tagMap {
					if err := generateTagPage(tagName, tagInfo, config, config.OutputDir); err != nil {
						fmt.Fprintf(os.Stderr, "Warning: error generating tag page for '%s': %v\n", tagName, err)
					}
				}

				// Generate tags index page
				if err := generateTagsIndex(tagMap, config, config.OutputDir); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: error generating tags index: %v\n", err)
				}
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

			// Apply blog-specific paths with proper nesting
			currentConfig := config
			currentConfig.CSSPath = prependPathPrefix(config.CSSPath, nestingDepth)
			currentConfig.JSPath = prependPathPrefix(config.JSPath, nestingDepth)
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
	logger := log.With().Str("dir", inputDir).Int("depth", depth).Logger()
	logger.Debug().Msg("Starting to process files in directory")

	files, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, fmt.Errorf("error reading input directory: %v", err)
	}
	logger.Debug().Int("file_count", len(files)).Msg("Files in directory")

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
			logger.Debug().Str("dir", file.Name()).Msg("Processing subdirectory recursively")
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
				subConfig.CSSPath = prependPathPrefix(config.CSSPath, blogDepth)
				subConfig.JSPath = prependPathPrefix(config.JSPath, blogDepth)
			} else {
				// For non-blog structure, use standard depth calculation
				subConfig = adjustPaths(subConfig, depth+1, subOutputDir)
			}

			// Process the subdirectory recursively and collect its blog posts
			subPosts, err := processFiles(filePath, inputRoot, subConfig, depth+1)
			if err != nil {
				return nil, err
			}

			// Prepend subdirectory name to each post's Link for correct relative paths
			for i := range subPosts {
				subPosts[i].Link = file.Name() + "/" + subPosts[i].Link
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
						logger.Debug().Str("file", file.Name()).Msg("Skipping landing page (already processed)")
						continue
					}
				}
			}

			logger.Debug().Str("file", file.Name()).Msg("Processing markdown file")
			metadata, err := ConvertFile(filePath, config, inputRoot)
			if err != nil {
				return nil, err
			}
			logger.Debug().Str("file", file.Name()).Msg("Markdown file processed successfully")

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

			// Read file content for search indexing
			mdContent, _ := os.ReadFile(filePath)
			_, contentWithoutMeta := extractMetadata(mdContent)
			plainContent := extractPlainText(contentWithoutMeta)

			// Get commit information if enabled
			var commitHash, commitDate, commitAuthor, commitURL string
			if config.ShowCommitInfo {
				gitInfo, err := GetFileGitInfo(filePath)
				if err == nil {
					commitHash = gitInfo.CommitHash
					commitDate = gitInfo.LastModified.Format("2006-01-02 15:04:05")
					commitAuthor = gitInfo.Author
					commitURL = GetCommitURL(config.GitWebURL, commitHash)
				}
			}

			// Calculate reading time
			readTime := calculateReadTime(contentWithoutMeta)

			// Extract tags
			postTags := normalizeTags(metadata["Tags"])

			blogPosts = append(blogPosts, BlogPost{
				Title:        title,
				Link:         fileNameWithoutExt + ".html",
				Description:  metadata["Description"],
				Date:         metadata["Date"],
				FullURL:      relPath,
				Author:       author,
				Content:      plainContent,
				FilePath:     filePath,
				CommitHash:   commitHash,
				CommitDate:   commitDate,
				CommitAuthor: commitAuthor,
				CommitURL:    commitURL,
				ReadTime:     readTime,
				Tags:         postTags,
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
			config,
		); err != nil {
			return nil, err
		}
	}

	// Generate tag pages if this is the blog root directory
	outComponents := strings.Split(config.OutputDir, string(filepath.Separator))
	blogIdx := -1
	for i, comp := range outComponents {
		if comp == "blog" {
			blogIdx = i
			break
		}
	}
	if blogIdx >= 0 {
		blogComponents := outComponents[blogIdx:]
		isBlogRoot := len(blogComponents) == 1

		if isBlogRoot && len(blogPosts) > 0 {
			tagMap := collectTags(blogPosts)

			// Get root output directory (dist/) by going up from blog directory
			rootOutputDir := filepath.Join(outComponents[:blogIdx]...)
			if rootOutputDir == "" {
				rootOutputDir = "."
			}

			// Generate individual tag pages
			for tagName, tagInfo := range tagMap {
				if err := generateTagPage(tagName, tagInfo, config, rootOutputDir); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: error generating tag page for '%s': %v\n", tagName, err)
				}
			}

			// Generate tags index page
			if err := generateTagsIndex(tagMap, config, rootOutputDir); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: error generating tags index: %v\n", err)
			}
		}
	}

	return blogPosts, nil
}

// fetchRecentPosts finds and returns the most recent blog posts sorted by Git commit date
func fetchRecentPosts(rootDir string, config Config, limit int) []BlogPost {
	logger := log.With().Str("root_dir", rootDir).Int("limit", limit).Logger()
	logger.Debug().Msg("Fetching recent posts from Git history")

	// Look for blog directory relative to root
	blogDir := filepath.Join(rootDir, "blog")
	if _, err := os.Stat(blogDir); os.IsNotExist(err) {
		logger.Debug().Msg("Blog directory not found, skipping recent posts")
		return nil
	}

	// Open Git repository
	gitRepo, err := OpenGitRepository(rootDir)
	if err != nil {
		logger.Debug().Err(err).Msg("Could not open Git repository")
		return nil
	}

	// Collect all markdown files from blog directory
	var mdFiles []string
	err = filepath.Walk(blogDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".markdown")) {
			// Skip index files
			if info.Name() != "index.md" && info.Name() != "index.markdown" {
				mdFiles = append(mdFiles, path)
			}
		}
		return nil
	})
	if err != nil {
		logger.Debug().Err(err).Msg("Error walking blog directory")
		return nil
	}

	logger.Debug().Int("files_found", len(mdFiles)).Msg("Found markdown files in blog directory")

	// Build posts with Git dates
	type postWithDate struct {
		post    BlogPost
		gitDate time.Time
	}
	var postsWithDates []postWithDate

	for _, filePath := range mdFiles {
		// Read file to get metadata
		mdContent, err := os.ReadFile(filePath)
		if err != nil {
			logger.Debug().Str("file", filePath).Err(err).Msg("Could not read markdown file")
			continue
		}

		metadata, _ := extractMetadata(mdContent)

		// Skip landing page templates
		if metadata["Template"] == "landing" {
			continue
		}

		// Get Git commit date
		gitInfo, err := gitRepo.GetFileLastModified(filePath)
		var gitDate time.Time
		if err != nil {
			logger.Debug().Str("file", filePath).Err(err).Msg("Could not get Git date, using metadata date")
			// Fall back to metadata date
			gitDate, _ = time.Parse("2006-01-02", metadata["Date"])
		} else {
			gitDate = gitInfo.LastModified
		}

		// Build the link path relative to site root
		relPath, err := filepath.Rel(rootDir, filePath)
		if err != nil {
			relPath = filePath
		}

		// Convert to HTML link
		fileNameWithoutExt := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
		fileNameWithoutExt = strings.ReplaceAll(fileNameWithoutExt, " ", "-")

		// Build the full link path (blog/Category/post.html)
		dirPath := filepath.Dir(relPath)
		link := filepath.Join(dirPath, fileNameWithoutExt+".html")
		link = strings.ReplaceAll(link, string(filepath.Separator), "/")

		// Get title from metadata or filename
		title := metadata["Title"]
		if title == "" {
			title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
		}

		// Get commit information if enabled
		var commitHash, commitDate, commitAuthor, commitURL string
		if config.ShowCommitInfo {
			if gitInfo != nil {
				commitHash = gitInfo.CommitHash
				commitDate = gitInfo.LastModified.Format("2006-01-02 15:04:05")
				commitAuthor = gitInfo.Author
				commitURL = GetCommitURL(config.GitWebURL, commitHash)
			}
		}

		// Calculate reading time
		_, contentWithoutMeta := extractMetadata(mdContent)
		readTime := calculateReadTime(contentWithoutMeta)

		// Extract tags
		postTags := normalizeTags(metadata["Tags"])

		post := BlogPost{
			Title:        title,
			Link:         link,
			Description:  metadata["Description"],
			Date:         gitDate.Format("2006-01-02"),
			FilePath:     filePath,
			CommitHash:   commitHash,
			CommitDate:   commitDate,
			CommitAuthor: commitAuthor,
			CommitURL:    commitURL,
			ReadTime:     readTime,
			Tags:         postTags,
		}

		postsWithDates = append(postsWithDates, postWithDate{post: post, gitDate: gitDate})
	}

	// Sort by Git date (newest first)
	sort.Slice(postsWithDates, func(i, j int) bool {
		return postsWithDates[i].gitDate.After(postsWithDates[j].gitDate)
	})

	// Limit the results
	if len(postsWithDates) > limit {
		postsWithDates = postsWithDates[:limit]
	}

	// Extract just the posts
	var posts []BlogPost
	for _, pwd := range postsWithDates {
		posts = append(posts, pwd.post)
	}

	logger.Debug().Int("posts_returned", len(posts)).Msg("Recent posts fetched successfully")
	return posts
}
