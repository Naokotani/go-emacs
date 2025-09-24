package main

import (
	"fmt"
	"os"
)

func (app *application) generatePosts() error {
	outputFile, err := os.Create(app.outputPath + "/index.html")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	ts, ok := app.templateCache["home.gotmpl"]

	if !ok {
		err := fmt.Errorf("Template does not exist in the template cache.")
		fmt.Printf("ERROR: %s", err)
	}

	fmt.Printf("year in index.go: %d", app.config.Date.Year())

	err = ts.ExecuteTemplate(outputFile, "base", app.config)
	if err != nil {
		return err
	}
	return nil
}
