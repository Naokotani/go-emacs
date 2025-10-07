package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/naokotani/go-emacs/internal/logger"
)

type application struct {
	templateCache map[string]*template.Template
	config        Config
	infoLog       *log.Logger
	errorLog      *log.Logger
	warnLog       *log.Logger
}

func main() {
	logger := logger.NewLogger(os.Getenv("LOG_LEVEL"))

	app := &application{
		infoLog:  logger.InfoLog,
		errorLog: logger.ErrorLog,
		warnLog:  logger.WarnLog,
	}

	app.parseConfig()
	app.removeAllOutDirs()
	app.buildOutputDirs()
	app.getPostDirs()
	app.getResumeFiles()
	app.infoLog.Printf("Loading templates in %s\n", app.config.TemplateDir)
	templateCache, err := newTemplateCache(app.config.TemplateDir)
	if err != nil {
		logger.ErrorLog.Fatal(err)
		return
	}

	app.templateCache = templateCache

	app.copySiteFiles()
	css := app.generateCssVarsFile()
	app.generateViews(css)

	fmt.Println()
	fmt.Printf(`************************************

Blog Successfully generated at %s

************************************`, app.config.Output)
}
