package main

import (
	"fmt"
	"html/template"
	"os"
	"time"
)

type application struct {
	outputPath         string
	configPath         string
	pagesTemplateCache map[string]*template.Template
	config             Config
}

func main() {
	outputPath := os.Getenv("OUTPUT_PATH")
	configPath := os.Getenv("CONFIG_PATH")
	pagesTemplateCache, err := newPagesTemplateCache("./ui/html/pages/*.gotmpl")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	config := Config{}
	config.Date = time.Now()
	fmt.Printf("Year: %d", config.Date.Year())

	app := &application{
		pagesTemplateCache: pagesTemplateCache,
		outputPath:         outputPath,
		configPath:         configPath,
		config:             config,
	}

	app.parseConfig()

	err = app.generatePages()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
}
