package main

import (
	"fmt"
	"html/template"
	"os"
	"time"
)

type application struct {
	outputPath    string
	configPath    string
	templateCache map[string]*template.Template
	config        Config
}

func main() {
	fmt.Println("Hello, World!")
	outputPath := os.Getenv("OUTPUT_PATH")
	configPath := os.Getenv("CONFIG_PATH")
	fmt.Printf("output path: %s\n", outputPath)
	templateCache, err := newTemplateCache()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	config := Config{}
	config.Date = time.Now()
	fmt.Printf("Year: %d", config.Date.Year())

	app := &application{
		templateCache: templateCache,
		outputPath:    outputPath,
		configPath:    configPath,
		config:        config,
	}

	app.parseConfig()

	err = app.generateIndex()
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}

	fmt.Println("Goodbye, World!")
}
