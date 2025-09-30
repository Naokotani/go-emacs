package main

import (
	"github.com/naokotani/go-emacs/internal/logger"
	"html/template"
	"log"
)

type application struct {
	templateCache map[string]*template.Template
	config        Config
	infoLog       *log.Logger
	errorLog      *log.Logger
	warnLog       *log.Logger
}

func main() {
	logger := logger.NewLogger()

	app := &application{
		infoLog:  logger.InfoLog,
		errorLog: logger.ErrorLog,
		warnLog:  logger.WarnLog,
	}

	app.parseConfig()
	app.getPostDirs()
	app.getResumeFiles()
	app.infoLog.Printf("Loading templates in %s\n", app.config.TemplateDir)
	templateCache, err := newTemplateCache(app.config.TemplateDir)
	if err != nil {
		logger.ErrorLog.Fatal(err)
		return
	}

	app.templateCache = templateCache

	app.buildOutputDirs()
	app.copySiteFiles()
	css := app.generateCssVarsFile()
	app.generateViews(css)
}
