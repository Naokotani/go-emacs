package main

import (
	"fmt"
	"html/template"
	"os"
)

type application struct {
	configPath         string
	pagesTemplateCache map[string]*template.Template
	config             Config
}

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	pagesTemplateCache, err := newTemplateCache("./ui/html/pages/*.gotmpl")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	app := &application{
		pagesTemplateCache: pagesTemplateCache,
		configPath:         configPath,
		config:             Config{},
	}

	app.parseConfig()

	css, err := app.generateCssVarsFile()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}

	err = app.generatePages(css)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}

}
