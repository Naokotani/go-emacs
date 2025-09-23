package main

import (
	"fmt"
	"os"
)

func (app *application) generateIndex() error {
	path := fmt.Sprintf("%s/index.html", app.outputPath)
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outputFile.Close() // Ensure the file is closed

	ts, ok := app.templateCache["home.gotmpl"]

	if !ok {
		err := fmt.Errorf("Template does not exist in the template cache.")
		fmt.Printf("ERROR: %s", err)
	}

	err = ts.Execute(outputFile, app.config)
	if err != nil {
		return err
	}
	return nil
}
