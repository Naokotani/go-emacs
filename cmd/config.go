package main

import (
	"fmt"
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

	config := Config{}

	toml.DecodeFile(f, &config)

	app.config = config

	fmt.Printf("Site data title: %s\n", app.config.Contact.Name)

	return nil
}
