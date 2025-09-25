package main

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Page struct {
	Title     string
	Content   template.HTML
	SubHeader string
	SkillTags string
	Post      Post
	Site      Site
	Meta      string
	Contact   Contact
	Date      time.Time
	Resume    Resume
	About     About
}

type Post struct {
	dir        string
	filename   string
	Title      string
	Summary    string
	Slug       string
	TagString  string
	Thumb      string
	Tags       []string
	DateString string
	Date       time.Time
}

func (app *application) generatePages() error {
	err := generatePosts(app)
	if err != nil {
		return err
	}
	err = generateAbout(app)
	if err != nil {
		return err
	}
	err = generateResume(app)
	if err != nil {
		return err
	}
	err = generateIndex(app)
	if err != nil {
		return err
	}
	return nil
}

func generateIndex(app *application) error {
	page := Page{
		Title:     app.config.Site.Title,
		SubHeader: app.config.Site.SubHeader,
		Contact:   app.config.Contact,
		Date:      time.Now(),
		Site:      app.config.Site,
		Resume:    app.config.Resume,
		About:     app.config.About,
	}

	output, err := os.Create(app.config.Output + "index.html")
	if err != nil {
		return err
	}
	defer output.Close()

	ts, ok := app.pagesTemplateCache["index.gotmpl"]
	if !ok {
		err := fmt.Errorf("Template index.html does not exist in the cache")
		return err
	}

	err = ts.ExecuteTemplate(output, "base", page)
	return nil
}

func generateResume(app *application) error {
	resume := app.config.Resume.Html
	if resume == "" {
		fmt.Println("Resume string empty, skipping page. Add to config.toml to generate resume.")
		return nil
	}
	fmt.Printf("Generating resume from %s\n", resume)

	page := Page{
		Title:     app.config.Site.Title,
		SubHeader: app.config.Site.SubHeader,
		Contact:   app.config.Contact,
		Date:      time.Now(),
		Site:      app.config.Site,
		Resume:    app.config.Resume,
		About:     app.config.About,
	}

	output, err := os.Create(app.config.Output + "resume.html")
	if err != nil {
		return err
	}
	defer output.Close()

	html, err := os.ReadFile(resume)
	if err != nil {
		return err
	}
	page.Content = template.HTML(html)

	ts, ok := app.pagesTemplateCache["resume.gotmpl"]
	if !ok {
		err := fmt.Errorf("Template resume.gotmpl does not exist in the cache")
		return err
	}

	err = ts.ExecuteTemplate(output, "base", page)
	return nil
}

func generateAbout(app *application) error {
	about := app.config.About.Html
	if about == "" {
		fmt.Println("About page intro string empty, skipping page. Add to config.toml to generate about page.")
		return nil
	}
	fmt.Printf("Generating about from %s\n", about)

	page := Page{
		Title:     app.config.Site.Title,
		SubHeader: app.config.Site.SubHeader,
		Contact:   app.config.Contact,
		Date:      time.Now(),
		Site:      app.config.Site,
		Resume:    app.config.Resume,
		About:     app.config.About,
	}

	output, err := os.Create(app.config.Output + "about.html")
	if err != nil {
		return err
	}
	defer output.Close()

	html, err := os.ReadFile(about)
	if err != nil {
		return err
	}
	page.Content = template.HTML(html)
	fmt.Printf("%s\n", page.Content)

	ts, ok := app.pagesTemplateCache["about.gotmpl"]
	if !ok {
		err := fmt.Errorf("Template resume.gotmpl does not exist in the cache")
		return err
	}

	err = ts.ExecuteTemplate(output, "base", page)
	return nil
}

func generatePosts(app *application) error {
	_, err := os.Stat(app.config.Output + "posts")

	if errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(app.config.Output+"posts", os.ModePerm); err != nil {
			return err
		}
	}

	var posts []Post
	for _, post := range app.config.Site.Posts {
		post, err = getPostMetadata(post)
		if err != nil {
			return err
		}

		page := Page{
			Title:   post.Title,
			Contact: app.config.Contact,
			Date:    time.Now(),
			Post:    post,
			Resume:  app.config.Resume,
			About:   app.config.About,
		}

		output, err := os.Create(app.config.Output + "posts/" + post.filename)
		if err != nil {
			return err
		}
		defer output.Close()

		html, err := os.ReadFile(post.dir + post.filename)
		if err != nil {
			return err
		}
		page.Content = template.HTML(html)

		ts, ok := app.pagesTemplateCache["post.gotmpl"]
		if !ok {
			err := fmt.Errorf("Template post.gotmpl does not exist in the cache")
			return err
		}

		err = ts.ExecuteTemplate(output, "base", page)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}

	app.config.Site.Posts = posts

	return nil
}

func getPostMetadata(post Post) (Post, error) {
	f := post.dir + "metadata.toml"
	if _, err := os.Stat(f); err != nil {
		return post, err
	}

	toml.DecodeFile(f, &post)
	fmt.Printf("post title: %s\n", post.Title)

	post.Tags = strings.Split(post.TagString, " ")
	post.TagString = strings.ReplaceAll(post.TagString, " ", " | ")
	post.Slug = "/posts/" + post.filename
	post.Thumb = "/images/" + strings.Split(post.filename, ".")[0] + ".png"
	layout := "Mon, 02 Jan 2006 15:04:05-07:00"
	t, err := time.Parse(layout, post.DateString)
	if err == nil {
		post.Date = t
	} else {
		fmt.Printf("Could not parse time for %s with time string: %s\n", post.filename, post.DateString)
	}

	return post, nil
}
