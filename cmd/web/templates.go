package main

import (
	"html/template"
	"path/filepath"

	"github.com/kvnloughead/snippetbox/internal/models"
)

// Go templates only allow a single data argument, so we create a struct to
// store all necessary template data.
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a map to serve as a cache.
	cache := map[string]*template.Template{}

	// Gather a slice of all `pages` templates.
	pages, err := filepath.Glob("ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	/* Grab the basename of each page, generate the corresponding template, and add it to the cache. */
	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"ui/html/base.tmpl",
			"ui/html/partials/nav.tmpl",
			page,
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}
