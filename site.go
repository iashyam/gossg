package main

import (
	"fmt"
	"gossg/parser"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
)

// Post represents a blog post
type Post struct {
	parser.Frontmatter
	ContentHTML template.HTML
	Slug        string
}

// Page represents a standalone page (like about, contact)
type Page struct {
	parser.Frontmatter
	ContentHTML template.HTML
	Slug        string
}

// Site holds all the content needed to generate the static site
type Site struct {
	Posts []Post
	Pages []Page
	Tags  map[string][]Post
	Cache *Cache
}

func NewSite() *Site {
	return &Site{
		Posts: []Post{},
		Pages: []Page{},
		Tags:  make(map[string][]Post),
		Cache: NewCache(".gossg_cache.json"),
	}
}

func (s *Site) LoadContent(contentDir string) error {
	// Load cache from disk
	if err := s.Cache.Load(); err != nil {
		fmt.Printf("Warning: failed to load cache: %v\n", err)
	}
	// 1. Load Posts
	postsDir := filepath.Join(contentDir, "posts")
	if err := s.loadPosts(postsDir); err != nil {
		return fmt.Errorf("error loading posts: %w", err)
	}

	// 2. Load Pages
	pagesDir := filepath.Join(contentDir, "pages")
	if err := s.loadPages(pagesDir); err != nil {
		return fmt.Errorf("error loading pages: %w", err)
	}

	// Save cache back to disk
	if err := s.Cache.Save(); err != nil {
		fmt.Printf("Warning: failed to save cache: %v\n", err)
	}

	return nil
}

func (s *Site) loadPosts(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist, that's fine
		}
		return err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		path := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hash := ComputeHash(content)
		slug := strings.TrimSuffix(file.Name(), ".md")

		var post Post
		if cachedFile, hit := s.Cache.Files[path]; hit && cachedFile.Hash == hash {
			// Cache Hit: file hasn't changed, skip Lexing and Parsing
			fmt.Printf("Cache hit: %s\n", path)
			post = Post{
				Frontmatter: cachedFile.Frontmatter,
				ContentHTML: template.HTML(cachedFile.ContentHTML),
				Slug:        slug,
			}
		} else {
			// Cache Miss: extract, parse, and update cache
			fmt.Printf("Cache miss: parsing %s...\n", path)
			fm, textContent, err := parser.ExtractFrontmatter(string(content))
			if err != nil {
				fmt.Printf("Warning: failed to parse frontmatter for %s: %v\n", path, err)
				continue
			}

			// Parse Markdown to HTML
			var buf strings.Builder
			md := goldmark.New(
				goldmark.WithExtensions(mathjax.MathJax),
			)
			if err := md.Convert([]byte(textContent), &buf); err != nil {
				fmt.Printf("Warning: failed to convert markdown for %s: %v\n", path, err)
				continue
			}
			htmlContent := buf.String()

			s.Cache.Files[path] = CachedFile{
				Hash:        hash,
				Frontmatter: fm,
				ContentHTML: htmlContent,
			}

			post = Post{
				Frontmatter: fm,
				ContentHTML: template.HTML(htmlContent),
				Slug:        slug,
			}
		}

		s.Posts = append(s.Posts, post)

		// Populate Tags map
		for _, tag := range post.Frontmatter.Tags {
			s.Tags[tag] = append(s.Tags[tag], post)
		}
	}

	// Sort posts by date descending (newest first)
	sort.Slice(s.Posts, func(i, j int) bool {
		return s.Posts[i].Date > s.Posts[j].Date
	})

	return nil
}

func (s *Site) loadPages(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		path := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hash := ComputeHash(content)
		slug := strings.TrimSuffix(file.Name(), ".md")

		var page Page
		if cachedFile, hit := s.Cache.Files[path]; hit && cachedFile.Hash == hash {
			// Cache Hit: file hasn't changed, skip Lexing and Parsing
			fmt.Printf("Cache hit: %s\n", path)
			page = Page{
				Frontmatter: cachedFile.Frontmatter,
				ContentHTML: template.HTML(cachedFile.ContentHTML),
				Slug:        slug,
			}
		} else {
			// Cache Miss: extract, parse, and update cache
			fmt.Printf("Cache miss: parsing %s...\n", path)
			fm, textContent, err := parser.ExtractFrontmatter(string(content))
			if err != nil {
				fmt.Printf("Warning: failed to parse frontmatter for %s: %v\n", path, err)
				continue
			}

			// Parse Markdown to HTML
			var buf strings.Builder
			md := goldmark.New(
				goldmark.WithExtensions(mathjax.MathJax),
			)
			if err := md.Convert([]byte(textContent), &buf); err != nil {
				fmt.Printf("Warning: failed to convert markdown for %s: %v\n", path, err)
				continue
			}
			htmlContent := buf.String()

			s.Cache.Files[path] = CachedFile{
				Hash:        hash,
				Frontmatter: fm,
				ContentHTML: htmlContent,
			}

			page = Page{
				Frontmatter: fm,
				ContentHTML: template.HTML(htmlContent),
				Slug:        slug,
			}
		}

		s.Pages = append(s.Pages, page)
	}

	return nil
}
