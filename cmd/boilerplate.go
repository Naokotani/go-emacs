package main

func (app *application) buildOutputDirs() {
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
}

func (app *application) copySiteFiles() {
	copyDirectory("./static/css", app.config.Output+"static/css")
	copyDirectory("./static/icons", app.config.Output+"static/icons")
	copyDirectory("./static/js", app.config.Output+"static/js")
}
