package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kreatoo/md2html/pkg/converter"
)

func main() {
	// Parse command line flags
	inputDir := flag.String("input", "", "Directory containing markdown files (required)")
	outputDir := flag.String("output", "output", "Directory to save HTML files")
	templateFile := flag.String("template", "", "HTML template file (optional)")
	cssPath := flag.String("css", "css/style-blog.css", "Path to CSS file (default: css/style-blog.css). Use '/path/to/file' for site-root-relative paths.")
	jsPath := flag.String("js", "js/script.js", "Path to JavaScript file (default: js/script.js). Use '/path/to/file' for site-root-relative paths.")
	siteTitle := flag.String("title", "Kreato's Website", "Site title")
	author := flag.String("author", "Kreato", "Default author name")
	createTemplate := flag.Bool("create-template", false, "Create a default template and exit")
	addList := flag.Bool("addlist", false, "Generate an index.html with a list of all blog posts")
	recursive := flag.Bool("recursive", false, "Process subdirectories recursively")
	generateRSS := flag.Bool("rss", false, "Generate RSS feed (feed.xml) for blog posts")
	siteURL := flag.String("site-url", "", "Site URL for RSS feed (e.g., https://krea.to)")
	flag.Parse()

	// Create default template if requested
	if *createTemplate {
		templatePath := "templates/default.html"
		if *templateFile != "" {
			templatePath = *templateFile
		}

		if err := converter.GenerateTemplateFile(templatePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating template: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Template created at %s\n", templatePath)
		os.Exit(0)
	}

	// Validate required fields
	if *inputDir == "" {
		fmt.Println("Error: input directory is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create default template if none provided
	if *templateFile == "" {
		tmpDir := filepath.Join(*outputDir, "tmp")
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating temp directory: %v\n", err)
			os.Exit(1)
		}

		defaultTemplate := filepath.Join(tmpDir, "default_template.html")
		if err := converter.GenerateTemplateFile(defaultTemplate); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating default template: %v\n", err)
			os.Exit(1)
		}

		*templateFile = defaultTemplate

		// Clean up temp directory at exit
		defer os.RemoveAll(tmpDir)
	}

	// Configure the converter
	config := converter.Config{
		TemplateFile:  *templateFile,
		OutputDir:     *outputDir,
		CSSPath:       *cssPath,
		JSPath:        *jsPath,
		SiteTitle:     *siteTitle,
		DefaultAuthor: *author,
		GenerateList:  *addList,
		Recursive:     *recursive,
		GenerateRSS:   *generateRSS,
		SiteURL:       *siteURL,
	}

	// Convert files
	if err := converter.ConvertDirectory(*inputDir, config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Conversion complete. Files saved to %s\n", *outputDir)
}
