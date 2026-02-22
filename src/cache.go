package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/iashyam/gossg/src/parser"
)

// CachedFile represents the parsed data and hash of a single markdown file
type CachedFile struct {
	Hash        string             `json:"hash"`
	Frontmatter parser.Frontmatter `json:"frontmatter"`
	ContentHTML string             `json:"content_html"`
}

// Cache manages the state of all processed files
type Cache struct {
	cachePath string
	Files     map[string]CachedFile `json:"files"`
}

// NewCache initializes a new Cache instance
func NewCache(cachePath string) *Cache {
	return &Cache{
		cachePath: cachePath,
		Files:     make(map[string]CachedFile),
	}
}

// Load attempts to read the cache from disk
func (c *Cache) Load() error {
	data, err := os.ReadFile(c.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No cache file yet, which is fine
		}
		return fmt.Errorf("failed to read cache: %w", err)
	}

	if err := json.Unmarshal(data, &c.Files); err != nil {
		return fmt.Errorf("failed to parse cache: %w", err)
	}
	return nil
}

// Save writes the current cache state to disk
func (c *Cache) Save() error {
	data, err := json.MarshalIndent(c.Files, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	if err := os.WriteFile(c.cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}
	return nil
}

// ComputeHash calculates the SHA-256 hash of the given content
func ComputeHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}
