package main

import (
	"html/template"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Site    Site
	Page    Page
	Contact Contact
	Date    time.Time
}

type Site struct {
	Meta string
}

type Page struct {
	Title string
	Post  template.HTML
}

type Contact struct {
	Name     string
	Git      string
	Email    string
	LinkedIn string
}

func (app *application) parseConfig() error {
	f := app.configPath
	if _, err := os.Stat(f); err != nil {
		return err
	}

	toml.DecodeFile(f, &app.config)

	return nil
}
