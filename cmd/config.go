package main

import (
	"flag"
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
	goEmacsDir   string
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
	app.config.goEmacsDir = app.getGoEmacsDir()

	toml.DecodeFile(filepath.Join(app.getXDGDocumentsDir("config.toml")), &app.config)

	setDefaultLocations(app)

	return nil
}

func setDefaultLocations(app *application) {
	if app.config.Output == "" {
		app.config.Output = app.getXDGDocumentsDir("blog")
		app.infoLog.Printf("Output directory not set. setting to XDG default: %s\n", app.config.Output)
	} else {
		app.infoLog.Printf("Output directory set to %s\n", app.config.Output)
	}

	if app.config.Posts.Dir == "" {
		app.config.Posts.Dir = app.getXDGDocumentsDir("posts")
		app.infoLog.Printf("Posts directory not set. setting to XDG default: %s\n", app.config.Posts.Dir)
	} else {
		app.infoLog.Printf("Posts directory set to %s\n", app.config.Posts.Dir)
	}

	if app.config.Pages.Dir == "" {
		app.config.Pages.Dir = app.getXDGDocumentsDir("pages")
		app.infoLog.Printf("Pages directory not set. setting to XDG default: %s\n", app.config.Pages.Dir)
	} else {
		app.infoLog.Printf("Pages directory set to %s\n", app.config.Pages.Dir)
	}

	if app.config.Resume.Dir == "" {
		app.config.Resume.Dir = app.getXDGDocumentsDir("resume")
		app.infoLog.Printf("Resume directory not set. setting to XDG default: %s\n", app.config.Resume.Dir)
	} else {
		app.infoLog.Printf("Resume directory set to %s\n", app.config.Resume.Dir)
	}

	if app.config.StylesConfig == "" {
		app.config.StylesConfig = app.getXDGDocumentsDir("styles.toml")
		app.infoLog.Printf("styles.toml not set, setting to XDG default: %s\n", app.config.StylesConfig)
	} else {
		app.infoLog.Printf("styles.toml set to %s\n", app.config.StylesConfig)
	}

	if app.config.TemplateDir == "" {
		app.config.TemplateDir = filepath.Join(app.config.goEmacsDir, "ui")
		app.infoLog.Printf("Template directory not set, setting to XDG default: %s\n", app.config.TemplateDir)
	} else {
		app.infoLog.Printf("Template directory set to %s\n", app.config.TemplateDir)
	}

	if app.config.StaticDir == "" {
		app.config.StaticDir = filepath.Join(app.config.goEmacsDir, "static")
		app.infoLog.Printf("StaticDir not set, setting to XDG default: %s\n", app.config.StaticDir)
	} else {
		app.infoLog.Printf("StaticDir set to %s\n", app.config.StaticDir)

	}

}

func getDefaultGoEmacsDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "emacs", "elpa", "cache", "go-emacs"), nil
}

func (app *application) getGoEmacsDir() string {
	configPath := flag.String("d", "", "path to config.toml file")
	flag.Parse()
	if *configPath == "" {
		var err error
		*configPath, err = getDefaultGoEmacsDir()
		if err != nil {
			app.errorLog.Fatalf("Failed to get config path")
		}
	}
	app.infoLog.Printf("Config path set to %s\n", *configPath)

	if !fileExists(*configPath) {
		app.errorLog.Fatalf("Go emacs directory does not exist at %s\n", *configPath)
	}

	return *configPath
}

func (app *application) getXDGDocumentsDir(dir string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		app.errorLog.Fatalf("Failed to set directory for %s", dir)
		return ""
	}
	return filepath.Join(home, "Documents/go-emacs", dir)
}
