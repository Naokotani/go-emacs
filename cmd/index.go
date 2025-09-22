package main

import (
	"fmt"
	"os"
)

func (app *application) generateIndex(data Data) {
	path := fmt.Sprintf("%s/index.html", app.outputPath)
	outputFile, err := os.Create(path)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
	defer outputFile.Close() // Ensure the file is closed

	ts, ok := app.templateCache["home.gotmpl"]

	if !ok {
		err := fmt.Errorf("Template does not exist in the template cache.")
		fmt.Printf("ERROR: %s", err)
	}

	err = ts.Execute(outputFile, data)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
}
