package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Output            string
	Contact           Contact
	StylesConfig      string
	staticDir         string
	templateDir       string
	Site              Site
	Resume            Resume
	Pages             Pages
	Posts             Posts
	GoEmacsPackageDir string
	goEmacsBlogDir    string
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
	app.config.goEmacsBlogDir, app.config.GoEmacsPackageDir = app.getGoEmacsDirs()

	toml.DecodeFile(filepath.Join(app.config.goEmacsBlogDir, "config.toml"), &app.config)

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

	app.config.templateDir = filepath.Join(app.config.GoEmacsPackageDir, "ui")
	app.config.staticDir = filepath.Join(app.config.GoEmacsPackageDir, "static")
}

func (app *application) getGoEmacsDirs() (string, string) {
	configPath := flag.String("d", "", "path to config.toml file")
	packageDir := flag.String("p", "", "path to go-emacs package")
	flag.Parse()

	if *configPath == "" {
		*configPath = app.getXDGDocumentsDir("")
	}

	app.infoLog.Printf("Config path set to %s\n", *configPath)

	if !fileExists(*configPath) {
		app.errorLog.Fatalf("Go emacs directory does not exist at %s\n", *configPath)
	}

	if *packageDir == "" && app.config.GoEmacsPackageDir != "" {
		app.infoLog.Printf("Package dir not set with -p. Falling back to config.toml path: %s\n",
			app.config.GoEmacsPackageDir)
		*packageDir = app.config.GoEmacsPackageDir
	} else if *packageDir == "" && app.config.GoEmacsPackageDir == "" {
		app.errorLog.Fatal("No emacs package directory set.\nMust set the path with the -p flag or in the config.toml file.")
	}

	return *configPath, *packageDir
}

func (app *application) getXDGDocumentsDir(dir string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		app.errorLog.Fatalf("Failed to set directory for %s", dir)
		return ""
	}
	return filepath.Join(home, "Documents/go-emacs", dir)
}
