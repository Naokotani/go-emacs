package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Output       string
	Contact      Contact
	StaticDir    string
	StylesConfig string
	TemplateDir  string
	Site         Site
	Resume       Resume
	Pages        Pages
	Posts        Posts
}

type Site struct {
	Meta      string
	Title     string
	Url       string
	Rss       bool
	Cards     bool
	SubHeader string
	SkillTags string
	Posts     []Post
	Pages     []Page
	Tags      map[string][]Post
	FontUrl   string
}

type Resume struct {
	IsResume bool
	Dir      string
	Pdf      string
}

type Pages struct {
	Dir string
}

type Posts struct {
	Dir string
}

type Contact struct {
	Name     string
	Git      string
	Email    string
	LinkedIn string
	Pfp      string
}

func (app *application) parseConfig() error {
	f := getConfigLocation(app, "CONFIG_PATH")
	app.infoLog.Printf("Config path set to %s\n", f)
	if _, err := os.Stat(f); err != nil {
		return err
	}

	toml.DecodeFile(f, &app.config)

	setDefaultLocations(app)

	return nil
}

func setDefaultLocations(app *application) {
	if app.config.Output == "" {
		app.config.Output = filepath.Join(getXDGGoEmacsDir(), "www")
		app.infoLog.Printf("Output directory not set. setting to XDG default: %s\n", app.config.Output)
	} else {
		app.infoLog.Printf("Output directory set to %s\n", app.config.Output)
	}

	if app.config.Posts.Dir == "" {
		app.config.Posts.Dir = filepath.Join(getXDGGoEmacsDir(), "posts")
		app.infoLog.Printf("Posts directory not set. setting to XDG default: %s\n", app.config.Posts.Dir)
	} else {
		app.infoLog.Printf("Posts directory set to %s\n", app.config.Posts.Dir)
	}

	if app.config.Pages.Dir == "" {
		app.config.Pages.Dir = filepath.Join(getXDGGoEmacsDir(), "pages")
		app.infoLog.Printf("Pages directory not set. setting to XDG default: %s\n", app.config.Pages.Dir)
	} else {
		app.infoLog.Printf("Pages directory set to %s\n", app.config.Pages.Dir)
	}

	if app.config.Resume.Dir == "" {
		app.config.Resume.Dir = filepath.Join(getXDGGoEmacsDir(), "resume")
		app.infoLog.Printf("Resume directory not set. setting to XDG default: %s\n", app.config.Resume.Dir)
	} else {
		app.infoLog.Printf("Resume directory set to %s\n", app.config.Resume.Dir)
	}

	if app.config.TemplateDir == "" {
		app.config.TemplateDir = filepath.Join(getXDGGoEmacsDir(), "ui")
		app.infoLog.Printf("styles.toml not set, setting to XDG default: %s\n", app.config.StylesConfig)
	} else {
		app.infoLog.Printf("styles.toml set to %s\n", app.config.StylesConfig)
	}

	if app.config.StaticDir == "" {
		app.config.StaticDir = filepath.Join(getXDGGoEmacsDir(), "static")
		app.infoLog.Printf("StaticDir not set, setting to XDG default: %s\n", app.config.StaticDir)
	} else {
		app.infoLog.Printf("StaticDir set to %s\n", app.config.StaticDir)

	}

}

func getConfigLocation(app *application, envVar string) string {
	if v := os.Getenv(envVar); v != "" {
		app.infoLog.Printf("Conig env var set, setting to %s", v)
		return v
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "~/Documents/go-emacs/config.toml"
	}

	configPath := filepath.Join(home, "Documents", "go-emacs", "config.toml")
	app.infoLog.Printf("Config env var not provided, setting to XDG default: %s\n", configPath)
	return configPath
}

func getXDGGoEmacsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/Documents/go-emacs"
	}
	return filepath.Join(home, "Documents/go-emacs")
}
