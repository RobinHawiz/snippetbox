package main

import "github.com/robinhawiz/snippetbox/internal/models"

type templateData struct {
	Snippet models.Snippet
	Snippets []models.Snippet
}