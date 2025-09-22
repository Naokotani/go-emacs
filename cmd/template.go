package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"
)

type post struct {
	title string
	date  time.Time
}

func humanDate(t time.Time) string {
	return t.Format("2 Jan 2006 at 15:04")
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.gotmpl")
	if err != nil {
		return nil, err
	}

	partials, err := filepath.Glob("./ui/html/partials/*.gotmpl")

	for _, partial := range partials {
		fmt.Printf("Partial: %s\n", partial)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		fmt.Printf("page: %s\n", page)

		ts, err := template.ParseFiles("./ui/html/base.gotmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.gotmpl")
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
