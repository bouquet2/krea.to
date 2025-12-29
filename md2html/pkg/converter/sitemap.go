package converter

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// URL represents a URL entry in sitemap.xml
type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

// URLSet is the root element of sitemap.xml
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// GenerateSitemap creates a sitemap.xml file for the website
func GenerateSitemap(outputDir string, siteURL string, inputRoot string) error {
	logger := log.With().Str("output_dir", outputDir).Str("site_url", siteURL).Logger()
	logger.Debug().Msg("Starting sitemap generation")

	// Validate site URL
	if siteURL == "" {
		return fmt.Errorf("site URL is required for sitemap generation")
	}

	// Open Git repository for last modification dates
	gitRepo, err := OpenGitRepository(inputRoot)
	if err != nil {
		logger.Debug().Err(err).Msg("Could not open Git repository, sitemap will not have lastmod dates")
		gitRepo = nil
	}

	// Collect all HTML files
	var urls []URL
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-HTML files
		if info.IsDir() {
			return nil
		}

		// Only process HTML files
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".html") {
			return nil
		}

		// Skip search-index.json (not a page)
		if strings.HasSuffix(info.Name(), "search-index.json") {
			return nil
		}

		// Calculate relative path from output directory
		relPath, err := filepath.Rel(outputDir, path)
		if err != nil {
			logger.Debug().Str("file", path).Err(err).Msg("Could not calculate relative path")
			return nil
		}

		// Convert to web path (forward slashes)
		webPath := strings.ReplaceAll(relPath, string(filepath.Separator), "/")

		// Build full URL
		fullURL := siteURL
		if !strings.HasSuffix(siteURL, "/") {
			fullURL += "/"
		}
		fullURL += webPath

		// Get last modification date from Git if available
		var lastMod string
		if gitRepo != nil {
			// Try to find the corresponding markdown file
			mdFile := strings.TrimSuffix(path, ".html") + ".md"
			if _, err := os.Stat(mdFile); err == nil {
				gitInfo, err := gitRepo.GetFileLastModified(mdFile)
				if err == nil {
					lastMod = gitInfo.LastModified.Format(time.RFC3339)
				}
			}
		}

		// Determine change frequency and priority based on file type
		changeFreq := "monthly"
		priority := 0.5

		if strings.HasSuffix(webPath, "index.html") {
			// Index pages are more important
			if webPath == "index.html" {
				// Homepage
				priority = 1.0
				changeFreq = "weekly"
			} else {
				// Category index pages
				priority = 0.8
				changeFreq = "weekly"
			}
		} else if strings.Contains(webPath, "blog/") {
			// Blog posts
			priority = 0.7
			changeFreq = "monthly"
		}

		urls = append(urls, URL{
			Loc:        fullURL,
			LastMod:    lastMod,
			ChangeFreq: changeFreq,
			Priority:   priority,
		})

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking output directory: %v", err)
	}

	// Create URLSet
	urlSet := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(urlSet, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling sitemap: %v", err)
	}

	// Add XML header
	xmlHeader := []byte(xml.Header)
	xmlData = append(xmlHeader, xmlData...)

	// Write sitemap.xml file
	sitemapPath := filepath.Join(outputDir, "sitemap.xml")
	if err := os.WriteFile(sitemapPath, xmlData, 0644); err != nil {
		return fmt.Errorf("error writing sitemap file: %v", err)
	}

	logger.Debug().Str("sitemap_file", sitemapPath).Int("url_count", len(urls)).Msg("Sitemap generated successfully")
	fmt.Printf("Generated sitemap: %s (%d URLs)\n", sitemapPath, len(urls))

	return nil
}
