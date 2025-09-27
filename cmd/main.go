package main

import (
	"github.com/naokotani/go-emacs/internal/logger"
	"html/template"
	"log"
	"os"
)

type application struct {
	configPath    string
	templateCache map[string]*template.Template
	config        Config
	infoLog       *log.Logger
	errorLog      *log.Logger
	warnLog       *log.Logger
}

func main() {
	logger := logger.NewLogger()

	app := &application{
		config:   Config{},
		infoLog:  logger.InfoLog,
		errorLog: logger.ErrorLog,
		warnLog:  logger.WarnLog,
	}

	app.configPath = os.Getenv("CONFIG_PATH")
	templateCache, err := newTemplateCache("./ui/html/templates/*.gotmpl")
	if err != nil {
		logger.ErrorLog.Fatal(err)
		return
	}
	app.templateCache = templateCache

	app.parseConfig()
	css := app.generateCssVarsFile()
	app.generatePages(css)
}
