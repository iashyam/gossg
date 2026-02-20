package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedFm   Frontmatter
		expectedBody string
		expectErr    bool
	}{
		{
			name: "Valid Frontmatter",
			input: `---
title: "My First Post"
date: "2023-10-01"
tags: ["go", "ssg"]
---
# Hello World
This is the body.`,
			expectedFm: Frontmatter{
				Title: "My First Post",
				Date:  "2023-10-01",
				Tags:  []string{"go", "ssg"},
			},
			expectedBody: "# Hello World\nThis is the body.",
			expectErr:    false,
		},
		{
			name: "No Frontmatter",
			input: `# Hello World
This is the body.`,
			expectedFm:   Frontmatter{},
			expectedBody: "# Hello World\nThis is the body.",
			expectErr:    false,
		},
		{
			name: "Malformed Frontmatter (no closing)",
			input: `---
title: "My First Post"
# Hello World`,
			expectedFm:   Frontmatter{},
			expectedBody: "---\ntitle: \"My First Post\"\n# Hello World",
			expectErr:    false,
		},
		{
			name:         "Empty File",
			input:        ``,
			expectedFm:   Frontmatter{},
			expectedBody: "",
			expectErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body, err := ExtractFrontmatter(tt.input)

			if (err != nil) != tt.expectErr {
				t.Errorf("ExtractFrontmatter() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !reflect.DeepEqual(fm, tt.expectedFm) {
				t.Errorf("ExtractFrontmatter() fm = %v, want %v", fm, tt.expectedFm)
			}

			// normalize line endings for comparison
			body = strings.ReplaceAll(body, "\r\n", "\n")
			tt.expectedBody = strings.ReplaceAll(tt.expectedBody, "\r\n", "\n")

			if body != tt.expectedBody {
				t.Errorf("ExtractFrontmatter() body = %v, want %v", body, tt.expectedBody)
			}
		})
	}
}
