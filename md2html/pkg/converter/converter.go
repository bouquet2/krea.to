package converter

import (
	"embed"
	"html/template"
)

//go:embed templates/*
var templateFS embed.FS

// Config holds configuration for the HTML generation
type Config struct {
	TemplateFile   string
	OutputDir      string
	CSSPath        string
	JSPath         string
	SiteTitle      string
	DefaultAuthor  string
	GenerateList   bool
	Recursive      bool
	GenerateRSS    bool
	SiteURL        string
	GitWebURL      string // Base URL for git web interface (e.g., https://github.com/user/repo/commit/)
	ShowCommitInfo bool   // Whether to display commit information in templates
	DefaultTheme   string // Default theme for the website
}

// PageData represents the data to be passed to the HTML template
type PageData struct {
	Title        string
	Content      template.HTML
	CSSPath      string
	JSPath       string
	Author       string
	Description  string
	Date         string
	URL          string
	Image        string
	CommitHash   string // Latest commit hash
	CommitDate   string // Commit date in readable format
	CommitAuthor string // Commit author
	CommitURL    string // URL to commit in git web interface
	DefaultTheme string // Default theme for the website
	ThemeCSSPath string // Path to the theme CSS file (if DefaultTheme is not nord)
}

// BlogPost represents a blog post entry for the index page
type BlogPost struct {
	Title        string
	Link         string
	Description  string
	Date         string
	FullURL      string // Full URL for RSS feed
	Author       string // Author for RSS feed
	Content      string // Plain text content for search indexing
	FilePath     string // Original file path for Git history lookup
	CommitHash   string // Latest commit hash
	CommitDate   string // Commit date in readable format
	CommitAuthor string // Commit author
	CommitURL    string // URL to commit in git web interface
	ReadTime     string // Estimated reading time (e.g., "~4 min read")
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
	DefaultTheme   string
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
	Title        string
	Description  string
	CSSPath      string
	JSPath       string
	Sections     []LandingSection
	Links        []LandingLink
	Settings     map[string]bool // Template settings from metadata
	Posts        []BlogPost      // Most recent blog posts from Git history
	DefaultTheme string          // Default theme for the website
	URL          string          // Full URL for og:url and canonical link
}

// ==========================================================================
// Primary Functions
// ==========================================================================
