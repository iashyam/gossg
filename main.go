package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// 1. Initialize Site
	site := NewSite()

	// 2. Load Content
	fmt.Println("Loading content...")
	if err := site.LoadContent("content"); err != nil {
		fmt.Printf("Error loading content: %v\n", err)
		return
	}

	// 3. Setup output directory
	if err := os.RemoveAll("public"); err != nil {
		fmt.Printf("Error clearing public dir: %v\n", err)
		return
	}
	if err := os.MkdirAll("public/posts", 0755); err != nil {
		fmt.Printf("Error creating public dir: %v\n", err)
		return
	}
	if err := os.MkdirAll("public/tags", 0755); err != nil {
		fmt.Printf("Error creating public tags dir: %v\n", err)
		return
	}

	// 4. Copy Static Assets
	fmt.Println("Copying assets...")
	if err := copyDir("content/assets", "public/assets"); err != nil {
		fmt.Printf("Warning: failed to copy assets: %v\n", err)
	}

	// 5. Load Templates
	// We parse the base template alongside each specific template
	postTmpl := template.Must(template.ParseFiles("templates/base.html", "templates/post.html"))
	indexTmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	listTmpl := template.Must(template.ParseFiles("templates/base.html", "templates/list.html"))
	tagsTmpl := template.Must(template.ParseFiles("templates/base.html", "templates/tags.html"))
	projTmpl := template.Must(template.ParseFiles("templates/base.html", "templates/projects.html"))

	// 5. Generate Pages
	for _, page := range site.Pages {
		generateFile(filepath.Join("public", page.Slug+".html"), postTmpl, page)
	}

	// 6. Generate Posts
	for _, post := range site.Posts {
		generateFile(filepath.Join("public", "posts", post.Slug+".html"), postTmpl, post)
	}

	// 7. Generate Home Page (Index) with Pagination
	postsPerPage := 6
	totalPosts := len(site.Posts)
	totalPages := (totalPosts + postsPerPage - 1) / postsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		startIdx := (pageNum - 1) * postsPerPage
		endIdx := startIdx + postsPerPage
		if endIdx > totalPosts {
			endIdx = totalPosts
		}

		pagePosts := site.Posts[startIdx:endIdx]

		data := map[string]interface{}{
			"Title":       "Home",
			"Posts":       pagePosts,
			"CurrentPage": pageNum,
			"TotalPages":  totalPages,
			"PrevPage":    pageNum - 1,
			"NextPage":    pageNum + 1,
		}

		var outputPath string
		if pageNum == 1 {
			outputPath = "public/index.html"
		} else {
			os.MkdirAll(fmt.Sprintf("public/page/%d", pageNum), 0755)
			outputPath = fmt.Sprintf("public/page/%d/index.html", pageNum)
		}

		generateFile(outputPath, indexTmpl, data)
	}

	// 8. Generate Timeline
	generateFile("public/timeline.html", listTmpl, map[string]interface{}{
		"Title": "Timeline",
		"Posts": site.Posts,
	})

	// 8. Generate Tags Index
	generateFile("public/tags.html", tagsTmpl, map[string]interface{}{
		"Title": "All Tags",
		"Tags":  site.Tags,
	})

	// 9. Generate Projects Page
	generateFile("public/projects.html", projTmpl, map[string]interface{}{
		"Title":    "Projects",
		"Projects": site.Projects,
	})

	// 10. Generate Individual Tag Pages
	for tag, posts := range site.Tags {
		generateFile(filepath.Join("public", "tags", tag+".html"), listTmpl, map[string]interface{}{
			"Title": "Tag: " + tag,
			"Posts": posts,
		})
	}

	fmt.Println("Site generation complete! Check the 'public' directory.")
}

func generateFile(outputPath string, tmpl *template.Template, data interface{}) {
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Failed to create file %s: %v\n", outputPath, err)
		return
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		fmt.Printf("Failed to execute template for %s: %v\n", outputPath, err)
	}
}

// copyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist.
func copyDir(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Source directory doesn't exist, ignore
		}
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}
