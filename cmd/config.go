package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Output       string
	Contact      Contact
	StylesConfig string
	Site         Site
	Resume       Resume
	Pages        Pages
}

type Site struct {
	Meta      string
	Title     string
	SubHeader string
	SkillTags string
	Posts     []Post
	Pages     []Page
	Tags      map[string][]Post
	FontUrl   string
}

type Resume struct {
	Dir  string
	Html string
	Pdf  string
}

type Pages struct {
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
	f := app.configPath
	if _, err := os.Stat(f); err != nil {
		return err
	}

	toml.DecodeFile(f, &app.config)

	postsDir := "./posts/"

	enteries, err := os.ReadDir(postsDir)
	var posts []Post
	if err != nil {
		return err
	}
	for _, e := range enteries {
		if e.IsDir() {
			posts = append(posts, Post{
				dir:      postsDir + e.Name() + "/",
				filename: e.Name() + ".html",
			})
		}
	}

	app.config.Site.Posts = posts

	return nil
}
