package main

import (
	"fmt"
	"html/template"
	"os"
)

type application struct {
	outputPath    string
	templateCache map[string]*template.Template
}

type Data struct {
	title   string
	content string
}

func main() {
	fmt.Println("Hello, World!")
	outputPath := os.Getenv("OUTPUT_PATH")
	fmt.Printf("output path: %s\n", outputPath)
	templateCache, err := newTemplateCache()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	data := Data{
		title: "foo",
	}

	app := &application{
		templateCache: templateCache,
		outputPath:    outputPath,
	}

	app.generateIndex(data)

	fmt.Println("Goodbye, World!")
}
