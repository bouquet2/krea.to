package converter

import (
	"strings"
)

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
