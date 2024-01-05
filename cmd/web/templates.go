package main

import "github.com/kvnloughead/snippetbox/internal/models"

// Go templates only allow a single data argument, so we create a struct to
// store all necessary template data.
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
