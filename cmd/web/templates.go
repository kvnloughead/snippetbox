package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/kvnloughead/snippetbox/internal/models"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// template.FuncMap struct provides a string keyed map of template functions.
// Must be registered with the template before calling ParseFiles.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Go templates only allow a single data argument, so we create a struct to
// store all necessary template data.
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
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

		// Before parsing base into a template set we call New() and Funcs() to
		// register the template functions.
		ts, err := template.New(name).Funcs(functions).ParseFiles(
			"./ui/html/base.tmpl",
		)
		if err != nil {
			return nil, err
		}

		// Call ParseGlob to add partials to the template set.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles to add the page template to the template set.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
