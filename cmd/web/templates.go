package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/tomesm/virtd/pkg/forms"
	"github.com/tomesm/virtd/pkg/models"
)

type templateData struct {
	CurrentYear int
	Course      *models.Course
	Courses     []*models.Course
	Form        *forms.Form
	Flash       string
}

// Return nicely formatted string of time.Time object
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	// Get slice of all 'page' templates for the app
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		fname := filepath.Base(page)
		//Register template.FuncMap
		ts, err := template.New(fname).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Add all layout/templates
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		// Add all partials
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[fname] = ts
	}
	return cache, nil
}
