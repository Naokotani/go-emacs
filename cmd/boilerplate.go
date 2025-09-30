package main

import (
	"os"
)

func (app *application) buildOutputDirs() {
	app.infoLog.Printf("Creating output dir: %s\n", app.config.Output)
	if !fileExists(app.config.Output) {
		if err := os.Mkdir(app.config.Output, os.ModePerm); err != nil {
			app.errorLog.Fatalf("Failed to create ouput directory %s\n%s", app.config.Output, err)
		}
	}

	app.makeOutputDir("/static")
	app.makeOutputDir("/static/css")
	app.makeOutputDir("/static/js")
	app.makeOutputDir("/static/icons")
	app.makeOutputDir("/static/pdf")
	app.makeOutputDir("/posts")
	app.makeOutputDir("/posts/images")
	app.makeOutputDir("/tags")
	app.makeOutputDir("/images")
	app.makeOutputDir("/images/thumbs")
	app.makeOutputDir("/resume")
}

func (app *application) removeAllOutDirs() {
	app.removeOutputDir("/static")
	app.removeOutputDir("/posts")
	app.removeOutputDir("/pages")
	app.removeOutputDir("/resume")
	app.removeOutputDir("/tags")
	app.removeOutputDir("/images")
}

func (app *application) copySiteFiles() {
	copyDirectory(app.config.StaticDir+"/css", app.config.Output+"/static/css")
	copyDirectory(app.config.StaticDir+"/icons", app.config.Output+"/static/icons")
	copyDirectory(app.config.StaticDir+"/js", app.config.Output+"/static/js")
}

func (app *application) removeOutputDir(dir string) {
	dir = app.config.Output + dir
	if fileExists(dir) {
		err := os.RemoveAll(dir)
		if err != nil {
			app.errorLog.Fatalf("Failed to remove output directory: %s\n%s\n", dir, err)
		}
	}
}
