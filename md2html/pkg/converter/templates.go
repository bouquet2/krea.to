package converter

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

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
