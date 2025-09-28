package main

import (
	"html/template"
	"path/filepath"
	"time"
)

type post struct {
	title string
	date  time.Time
}

func humanDate(t time.Time) string {
	return t.Format("2 Jan 2006 at 15:04 MST")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(path string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/base.gotmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/partials/*.gotmpl")
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
