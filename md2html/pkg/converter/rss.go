package converter

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

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
