package main

import (
	"fmt"
	"os"
	"strings"
)

func (app *application) generatePages() error {
	// outputFile, err := os.Create(app.outputPath + "/index.html")
	// if err != nil {
	// 	return err
	// }
	// defer outputFile.Close()

	for key, ts := range app.templateCache {
		filename := strings.Split(key, ".")[0] + ".html"
		path := app.outputPath + "/" + filename
		output, err := os.Create(path)
		if err != nil {
			return err
		}
		defer output.Close()
		fmt.Printf("template: %s\n", filename)
		err = ts.ExecuteTemplate(output, "base", app.config)
		if err != nil {
			return err
		}
	}

	// ts, ok := app.templateCache["home.gotmpl"]
	//
	// if !ok {
	// 	err := fmt.Errorf("Template does not exist in the template cache.")
	// 	fmt.Printf("ERROR: %s", err)
	// }
	//
	// fmt.Printf("year in index.go: %d", app.config.Date.Year())
	//
	// err = ts.ExecuteTemplate(outputFile, "base", app.config)
	// if err != nil {
	// 	return err
	// }
	return nil
}
