package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/valentedev/template-server/internal/models"
)

type templateData struct {
	CurrentYear int
	Book        models.Book
	Books       []models.Book
	Form        any
}

func humanDate(t time.Time) string {
	return t.Format(time.DateTime)
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Path to all templates under "./ui/html/pages/"
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// files := []string{
		// 	"./ui/html/base.html",
		// 	"./ui/html/components/nav.html",
		// 	page,
		// }

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// ts, err = template.ParseFiles("./ui/html/base.html")
		// if err != nil {
		// 	return nil, err
		// }

		ts, err = ts.ParseGlob("./ui/html/components/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil

}
