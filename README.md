# goSSG

A blazing fast, minimalist Static Site Generator written in Go.

goSSG is designed to transform markdown files into a beautiful, static blog/portfolio out of the box. It features automatic pagination, tag generation, a project showcase, built-in MathJax support, and incremental caching for ultra-fast builds.

## Features

- **Blazing Fast**: Uses Go and an incremental caching system to only parse markdown files that have changed.
- **Markdown & Math**: Supports standard Markdown along with full MathJax rendering (`$$` and `\\(`) out of the box.
- **Pre-configured Layouts**: Comes with a fully responsive, beautiful default theme (using Tailwind CSS).
- **Zero Dependencies via Embedding**: The default HTML templates are directly embedded into the Go binary. You only need the single executable to generate a full site.

## Installation

You can install goSSG directly via Go:

```bash
go install github.com/iashyam/gossg@latest
```

## Setup & Usage

To use goSSG, you'll need to set up a specific directory structure. Create a new directory for your website and recreate the following structure:

```
my-website/
├── config.yaml
└── content/
    ├── assets/       # Put your images here
    ├── pages/        # Standalone pages like about.md
    ├── posts/        # Your blog posts (e.g., 2024-01-01-hello.md)
    └── projects/     # Your project pages
```

### Configuration (`config.yaml`)

Create a `config.yaml` file at the root of your project to specify your base URL:

```yaml
baseURL: "https://yourdomain.com"
```
*(If you are deploying to a subpath like GitHub Pages, use `"https://username.github.io/repo"`)*

### Creating Content

Write your content in Markdown files. Every markdown file must include YAML frontmatter at the top:

```yaml
---
title: "My First Post"
date: "2024-01-01"
tags: ["programming", "welcome"]
image: "/assets/cover.jpg" # Optional
description: "A short description." # Useful for projects
---
# Main Content
Hello world!
```

### Building the Site

Once your content is ready, simply run the `gossg` command from the root of your project structure (where `config.yaml` is located):

```bash
gossg
```

This will parse your content and generate a static website inside a new `public/` directory. You can then host this `public/` directory on GitHub Pages, Vercel, Netlify, or any static hosting platform.

## Customizing Templates

Currently, **dynamic custom template support is not implemented**. The HTML templates required to build the site are fully embedded inside the goSSG binary using Go's `//go:embed` functionality.

If you would like to customize the colors, structure, or HTML of the website, you must fork or clone this repository, modify the templates, and build the binary yourself.

### How to customize:

1. Clone the repository:
   ```bash
   git clone https://github.com/iashyam/gossg.git
   cd gossg
   ```
2. Edit the HTML and Tailwind CSS structure found inside the `src/templates/` directory (e.g., modify `src/templates/base.html`).
3. Rebuild and install your custom version of the binary locally:
   ```bash
   go install .
   ```
4. Now, running `gossg` on your machine will use your customized embedded templates!

## Publishing to GitHub Pages

goSSG is built perfectly to be automated via GitHub Actions for seamless deployment to GitHub Pages.

1. **Configure your URL**: Ensure your `config.yaml` uses your GitHub Pages URL:
   ```yaml
   baseURL: "https://username.github.io/repo"
   ```
2. **Setup the Workflow**: Create a `.github/workflows/deploy.yml` file in your repository:
   ```yaml
   name: Deploy goSSG to GitHub Pages

   on:
     push:
       branches: ["main"]

   permissions:
     contents: write

   jobs:
     build-and-deploy:
       runs-on: ubuntu-latest
       steps:
         - name: Checkout Code
           uses: actions/checkout@v4

         - name: Setup Go
           uses: actions/setup-go@v4
           with:
             go-version: '1.21'

         - name: Build gossg
           run: go build -o gossg
             
         - name: Generate Static Site
           run: ./gossg
             
         - name: Deploy to GitHub Pages
           uses: peaceiris/actions-gh-pages@v3
           with:
             github_token: ${{ secrets.GITHUB_TOKEN }}
             publish_dir: ./public
   ```
3. **Configure Pages Setting**: In your repository Settings -> Pages, set the source to deploy from the `gh-pages` branch. Every time you push to `main`, GitHub Actions will build your goSSG site and deploy it automatically!
