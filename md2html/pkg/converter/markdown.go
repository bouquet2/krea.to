package converter

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"regexp"
	"strings"
	"unicode"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

// Average reading speed in words per minute
const wordsPerMinute = 200

// calculateReadTime estimates reading time based on word count
// Returns a string like "~4 min read"
func calculateReadTime(content []byte) string {
	// Extract plain text from markdown
	plainText := extractPlainText(content)

	// Count words by splitting on whitespace
	words := strings.Fields(plainText)
	wordCount := len(words)

	// Calculate minutes, minimum 1 minute
	minutes := wordCount / wordsPerMinute
	if minutes < 1 {
		minutes = 1
	}

	return fmt.Sprintf("~%d min read", minutes)
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
			if lang == "" {
				lang = "plaintext"
			}
			// Wrap in a div with data-lang attribute for reliable JS detection
			w.Write([]byte(fmt.Sprintf(`<div class="code-wrapper" data-lang="%s">`, lang)))
			renderCodeWithSyntaxHighlighting(&buf, lang, cb.Literal)
			w.Write(buf.Bytes())
			w.Write([]byte(`</div>`))
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
	policy.AllowAttrs("data-lang").OnElements("div")
	sanitizedHTML := policy.SanitizeBytes(htmlContent)

	return template.HTML(sanitizedHTML)
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

// extractPlainText converts markdown content to plain text for search indexing
func extractPlainText(content []byte) string {
	text := string(content)

	// Remove code blocks (``` ... ```)
	codeBlockRe := regexp.MustCompile("(?s)```[^`]*```")
	text = codeBlockRe.ReplaceAllString(text, " ")

	// Remove inline code (`...`)
	inlineCodeRe := regexp.MustCompile("`[^`]+`")
	text = inlineCodeRe.ReplaceAllString(text, " ")

	// Remove images ![alt](url)
	imageRe := regexp.MustCompile(`!\[[^\]]*\]\([^)]+\)`)
	text = imageRe.ReplaceAllString(text, " ")

	// Convert links [text](url) to just text
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	text = linkRe.ReplaceAllString(text, "$1")

	// Remove HTML comments
	commentRe := regexp.MustCompile(`<!--[\s\S]*?-->`)
	text = commentRe.ReplaceAllString(text, " ")

	// Remove HTML tags
	htmlRe := regexp.MustCompile(`<[^>]+>`)
	text = htmlRe.ReplaceAllString(text, " ")

	// Remove markdown formatting characters
	text = strings.ReplaceAll(text, "#", " ")
	text = strings.ReplaceAll(text, "*", "")
	text = strings.ReplaceAll(text, "_", "")
	text = strings.ReplaceAll(text, "~", "")
	text = strings.ReplaceAll(text, ">", " ")
	text = strings.ReplaceAll(text, "|", " ")

	// Normalize whitespace
	spaceRe := regexp.MustCompile(`\s+`)
	text = spaceRe.ReplaceAllString(text, " ")

	// Trim and limit length for reasonable index size
	text = strings.TrimSpace(text)

	// Remove non-printable characters
	text = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			return r
		}
		return -1
	}, text)

	return text
}
