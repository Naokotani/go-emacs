package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	xmlTemplate "text/template"
	"time"
)

type post struct {
	title string
	date  time.Time
}

func humanDate(t time.Time) string {
	return t.Format("2 Jan 2006 at 15:04 MST")
}

func rssDate(t time.Time) string {
	local := t.In(time.Now().Location())
	return local.Format(time.RFC1123Z)
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(path string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(path, "templates/*.gotmpl"))
	if err != nil {
		return nil, fmt.Errorf("Failed to scan templates %s:", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		if name == "rss.gotmpl" {
			break
		}

		ts, err := template.New(name).Funcs(functions).ParseFiles(filepath.Join(path, "base.gotmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, fmt.Errorf("Cant read templates in %s: %s", page, err)
		}

		ts, err = ts.ParseGlob(filepath.Join(path, "partials/*.gotmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (app *application) getRssTemplate() *xmlTemplate.Template {
	path := app.config.TemplateDir + "/templates/rss.gotmpl"

	funcs := xmlTemplate.FuncMap{
		"rssDate": rssDate,
	}

	ts, err := xmlTemplate.New("rss").Funcs(funcs).ParseFiles(path)
	if err != nil {
		app.errorLog.Fatalf("Failed to parse rss template %s", path)
	}
	return ts
}
