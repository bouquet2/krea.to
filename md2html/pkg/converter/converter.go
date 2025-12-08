package converter

import (
	"embed"
	"html/template"
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
	Content     string // Plain text content for search indexing
	FilePath    string // Original file path for Git history lookup
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
	Settings    map[string]bool // Template settings from metadata
	Posts       []BlogPost      // Most recent blog posts from Git history
}

// ==========================================================================
// Primary Functions
// ==========================================================================
