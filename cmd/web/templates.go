package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/valentedev/template-server/internal/models"
	"github.com/valentedev/template-server/ui"
)

type templateData struct {
	CurrentYear     int
	Book            models.Book
	Books           []models.Book
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format(time.DateTime)
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// pages, err := filepath.Glob("./ui/html/pages/*.html")
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		// if err != nil {
		// 	return nil, err
		// }

		// ts, err = ts.ParseGlob("./ui/html/components/*.html")
		// if err != nil {
		// 	return nil, err
		// }

		// ts, err = ts.ParseFiles(page)
		// if err != nil {
		// 	return nil, err
		// }

		patterns := []string{
			"html/base.html",
			"html/components/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil

}
