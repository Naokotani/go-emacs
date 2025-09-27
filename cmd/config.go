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
	About        About
}

type Site struct {
	Meta      string
	Title     string
	SubHeader string
	SkillTags string
	Posts     []Post
	Tags      map[string][]Post
	FontUrl   string
}

type Resume struct {
	Html string
	Pdf  string
}

type About struct {
	Html  string
	Image string
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

	postsDir := "./ui/html/pages/posts/"

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
