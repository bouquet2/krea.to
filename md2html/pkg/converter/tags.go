package converter

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

// collectTags aggregates all tags from blog posts and builds a tag map
func collectTags(posts []BlogPost) map[string]*TagInfo {
	tagMap := make(map[string]*TagInfo)

	for _, post := range posts {
		for _, tag := range post.Tags {
			if tagMap[tag] == nil {
				tagMap[tag] = &TagInfo{
					Name:  tag,
					Count: 0,
					Posts: []BlogPost{},
				}
			}
			tagMap[tag].Count++
			tagMap[tag].Posts = append(tagMap[tag].Posts, post)
		}
	}

	return tagMap
}

// getRelatedPosts finds posts that share tags with the current post
func getRelatedPosts(currentPost BlogPost, allPosts []BlogPost, limit int) []BlogPost {
	if len(currentPost.Tags) == 0 {
		return nil
	}

	type scoredPost struct {
		post  BlogPost
		score int
	}

	var scored []scoredPost

	// Create a map of current post tags for quick lookup
	currentTags := make(map[string]bool)
	for _, tag := range currentPost.Tags {
		currentTags[tag] = true
	}

	// Score each post based on shared tags
	for _, post := range allPosts {
		// Skip the current post itself
		if post.Link == currentPost.Link {
			continue
		}

		score := 0
		for _, tag := range post.Tags {
			if currentTags[tag] {
				score++
			}
		}

		if score > 0 {
			scored = append(scored, scoredPost{post: post, score: score})
		}
	}

	// Sort by score (highest first), then by date
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].score != scored[j].score {
			return scored[i].score > scored[j].score
		}
		return scored[i].post.Date > scored[j].post.Date
	})

	// Limit results
	if len(scored) > limit {
		scored = scored[:limit]
	}

	// Extract just the posts
	var related []BlogPost
	for _, sp := range scored {
		related = append(related, sp.post)
	}

	return related
}

// generateTagPage generates an HTML page for a specific tag
func generateTagPage(tagName string, tagInfo *TagInfo, config Config, outputBaseDir string) error {
	logger := log.With().Str("tag", tagName).Logger()
	logger.Debug().Msg("Generating tag page")

	// Load tag template
	tmpl, err := template.ParseFS(templateFS, "templates/tag.html")
	if err != nil {
		return fmt.Errorf("error parsing tag template: %v", err)
	}

	// Create tags directory inside blog directory
	tagsDir := filepath.Join(outputBaseDir, "blog", "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return fmt.Errorf("error creating tags directory: %v", err)
	}

	// Create output file
	outputFile := filepath.Join(tagsDir, tagName+".html")
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating tag page file: %v", err)
	}
	defer file.Close()

	// Sort posts by date (newest first)
	sortedPosts := make([]BlogPost, len(tagInfo.Posts))
	copy(sortedPosts, tagInfo.Posts)
	sort.Slice(sortedPosts, func(i, j int) bool {
		return sortedPosts[i].Date > sortedPosts[j].Date
	})

	// Calculate depth for CSS/JS paths
	// tags/ is two levels deep from root (blog/tags/)
	cssPath := prependPathPrefix(config.CSSPath, 2)
	jsPath := prependPathPrefix(config.JSPath, 2)

	// Build URL for the tag page
	url := ""
	if config.SiteURL != "" {
		url = strings.TrimSuffix(config.SiteURL, "/") + "/blog/tags/" + tagName + ".html"
	}

	data := TagPageData{
		Title:        fmt.Sprintf("Posts tagged with: %s", tagName),
		TagName:      tagName,
		CSSPath:      cssPath,
		JSPath:       jsPath,
		Posts:        sortedPosts,
		BackURL:      "index.html",
		URL:          url,
		DefaultTheme: config.DefaultTheme,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing tag template: %v", err)
	}

	logger.Debug().Str("output_file", outputFile).Int("post_count", len(sortedPosts)).Msg("Tag page generated")
	fmt.Printf("Generated tag page: %s (%d posts)\n", outputFile, len(sortedPosts))
	return nil
}

// generateTagsIndex generates the tags overview page
func generateTagsIndex(tagMap map[string]*TagInfo, config Config, outputBaseDir string) error {
	logger := log.With().Str("type", "tags_index").Logger()
	logger.Debug().Msg("Generating tags index page")

	// Load tags index template
	tmpl, err := template.ParseFS(templateFS, "templates/tags-index.html")
	if err != nil {
		return fmt.Errorf("error parsing tags index template: %v", err)
	}

	// Create tags directory inside blog directory
	tagsDir := filepath.Join(outputBaseDir, "blog", "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return fmt.Errorf("error creating tags directory: %v", err)
	}

	// Create output file
	outputFile := filepath.Join(tagsDir, "index.html")
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating tags index file: %v", err)
	}
	defer file.Close()

	// Convert map to sorted slice
	var tags []TagInfo
	for _, tagInfo := range tagMap {
		tags = append(tags, *tagInfo)
	}

	// Sort tags alphabetically
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	// Calculate depth for CSS/JS paths
	cssPath := prependPathPrefix(config.CSSPath, 2)
	jsPath := prependPathPrefix(config.JSPath, 2)

	// Build URL
	url := ""
	if config.SiteURL != "" {
		url = strings.TrimSuffix(config.SiteURL, "/") + "/blog/tags/index.html"
	}

	data := TagsIndexData{
		Title:        "All Tags",
		CSSPath:      cssPath,
		JSPath:       jsPath,
		Tags:         tags,
		BackURL:      "../index.html",
		URL:          url,
		DefaultTheme: config.DefaultTheme,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing tags index template: %v", err)
	}

	logger.Debug().Str("output_file", outputFile).Int("tag_count", len(tags)).Msg("Tags index generated")
	fmt.Printf("Generated tags index: %s (%d tags)\n", outputFile, len(tags))
	return nil
}

// getPopularTags returns the top N tags by post count
func getPopularTags(tagMap map[string]*TagInfo, limit int) []TagInfo {
	var tags []TagInfo
	for _, tagInfo := range tagMap {
		tags = append(tags, *tagInfo)
	}

	// Sort by count (descending), then alphabetically
	sort.Slice(tags, func(i, j int) bool {
		if tags[i].Count != tags[j].Count {
			return tags[i].Count > tags[j].Count
		}
		return tags[i].Name < tags[j].Name
	})

	// Limit results
	if len(tags) > limit {
		tags = tags[:limit]
	}

	return tags
}
