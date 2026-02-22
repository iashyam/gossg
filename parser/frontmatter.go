package parser

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// Frontmatter represents the metadata at the top of a Markdown file.
type Frontmatter struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Tags        []string `yaml:"tags"`
	Image       string   `yaml:"image"`
	Link        string   `yaml:"link"`
	Description string   `yaml:"description"`
}

// ExtractFrontmatter separates the YAML frontmatter from the Markdown content.
func ExtractFrontmatter(content string) (Frontmatter, string, error) {
	var fm Frontmatter
	content = strings.TrimSpace(content)

	// Check if the file starts with the frontmatter delimiter
	if !strings.HasPrefix(content, "---") {
		return fm, content, nil // No frontmatter
	}

	// Find the end of the frontmatter block
	endIdx := strings.Index(content[3:], "---")
	if endIdx == -1 {
		return fm, content, nil // Malformed frontmatter, treat as no frontmatter
	}

	// Extract the YAML block (offset by 3 to account for the first "---")
	yamlBlock := content[3 : endIdx+3]

	// Parse the YAML
	err := yaml.Unmarshal([]byte(yamlBlock), &fm)
	if err != nil {
		return fm, content, err
	}

	// Extract the rest of the Markdown content
	// We offset by 6 (3 for first "---" + 3 for second "---")
	markdownContent := content[endIdx+6:]

	// Trim leading whitespace/newlines from the actual markdown
	return fm, strings.TrimSpace(markdownContent), nil
}
