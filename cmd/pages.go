package main

import (
	"errors"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/naokotani/go-emacs/internal/images"
	_ "github.com/naokotani/go-emacs/internal/images"
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
	Css       Css
	Tags      map[string][]Post
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

func (app *application) generatePages(css Css) {
	page := Page{
		Title:     app.config.Site.Title,
		SubHeader: app.config.Site.SubHeader,
		Contact:   app.config.Contact,
		Date:      time.Now(),
		Site:      app.config.Site,
		Resume:    app.config.Resume,
		About:     app.config.About,
		Css:       css,
	}

	generatePosts(app, page)
	generateAbout(app, page)
	generateResume(app, page)
	page.Tags = app.config.Site.Tags
	generateIndex(app, page)

	for tag, posts := range page.Tags {
		page.Site.Posts = posts
		generateTagHome(app, page, tag)
	}
}

func generateIndex(app *application, page Page) {
	page.Site.Posts = app.config.Site.Posts

	output, err := os.Create(app.config.Output + "index.html")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer output.Close()

	ts, ok := app.templateCache["index.gotmpl"]
	if !ok {
		app.errorLog.Fatal("Template index.html does not exist in the cache")
	}

	err = ts.ExecuteTemplate(output, "base", page)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func generateTagHome(app *application, page Page, tag string) {
	app.makeOutputDir("tags")

	output, err := os.Create(app.config.Output + "tags/" + tag + ".html")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer output.Close()

	ts, ok := app.templateCache["index.gotmpl"]
	if !ok {
		app.errorLog.Fatal("Template index.html does not exist in the cache")
	}

	err = ts.ExecuteTemplate(output, "base", page)

	for k := range page.Tags {
		app.infoLog.Printf("Generating tag pages: %s ", k)
	}
}

func generateResume(app *application, page Page) {
	resume := app.config.Resume.Html
	if resume == "" {
		app.warnLog.Printf("Resume string empty, skipping page. Add to config.toml to generate resume.\n")
		return
	}
	app.infoLog.Printf("Generating resume from %s\n", resume)

	output, err := os.Create(app.config.Output + "resume.html")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer output.Close()

	html, err := os.ReadFile(resume)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	page.Content = template.HTML(html)

	ts, ok := app.templateCache["resume.gotmpl"]
	if !ok {
		app.errorLog.Fatal("Template resume.gotmpl does not exist in the cache")
	}

	err = ts.ExecuteTemplate(output, "base", page)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func generateAbout(app *application, page Page) {
	about := app.config.About.Html
	if about == "" {
		app.warnLog.Printf("About page intro string empty, skipping page. Add to config.toml to generate about page.\n")
		return
	}
	app.infoLog.Printf("Generating about from %s\n", about)

	output, err := os.Create(app.config.Output + "about.html")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer output.Close()

	html, err := os.ReadFile(about)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	page.Content = template.HTML(html)

	ts, ok := app.templateCache["about.gotmpl"]
	if !ok {
		app.errorLog.Fatal("Template resume.gotmpl does not exist in the cache")
	}

	err = ts.ExecuteTemplate(output, "base", page)
}

func generatePosts(app *application, page Page) {
	app.makeOutputDir("posts")
	app.makeOutputDir("images/thumbs")

	var posts []Post
	for _, post := range app.config.Site.Posts {
		post := getPostMetadata(app, post)
		post = app.createThumb(post)

		page.Title = post.Title
		page.Post = post

		output, err := os.Create(app.config.Output + "posts/" + post.filename)
		if err != nil {
			app.errorLog.Fatal(err)
		}
		defer output.Close()

		html, err := os.ReadFile(post.dir + post.filename)
		if err != nil {
			app.errorLog.Fatal(err)
		}
		page.Content = template.HTML(html)

		postFile := "post.gotmpl"
		ts, ok := app.templateCache[postFile]
		if !ok {
			app.errorLog.Fatalf("Template %s does not exist in the cache", postFile)
		}

		err = ts.ExecuteTemplate(output, "base", page)
		if err != nil {
			app.errorLog.Fatal(err)
		}
		posts = append(posts, post)
	}

	app.config.Site.Tags = make(map[string][]Post)
	app.config.Site.Posts = posts
	for _, p := range posts {
		for _, t := range p.Tags {
			if t != "" {
				app.config.Site.Tags[t] = append(app.config.Site.Tags[t], p)
			}
		}
	}
}

func getPostMetadata(app *application, post Post) Post {
	f := post.dir + "metadata.toml"
	if !fileExists(f) {
		app.errorLog.Fatalf("Metadata file does not exist in %s", post.dir)
		return post
	}

	toml.DecodeFile(f, &post)

	post.Tags = strings.Split(post.TagString, " ")
	post.TagString = strings.ReplaceAll(post.TagString, " ", " | ")
	post.Slug = "/posts/" + post.filename
	layout := "Mon, 02 Jan 2006 15:04:05-07:00"
	t, err := time.Parse(layout, post.DateString)
	post.Thumb = ""

	app.infoLog.Printf("Reading post for file: %s\n", post.filename)
	app.logPostdata("title", post.Title, post.filename)
	app.logPostdata("image", post.Title, post.filename)
	app.logPostdata("summary", post.Title, post.filename)

	if err == nil {
		post.Date = t
		app.logPostdata("date", post.Title, post.filename)
	} else {
		app.warnLog.Printf("Could not parse time for %s with time string: %s\n", post.filename, post.DateString)
	}

	return post
}

func (app *application) createThumb(post Post) Post {
	output := app.config.Output + "images/thumbs/" + strings.Split(post.filename, ".")[0] + ".png"
	switch {
	case fileExists(post.dir + "thumb.png"):
		if images.IsThumbTooWide(post.dir+"thumb.png", 200) {
			images.ResizePng(post.dir+"thumb.png", output, 0, 200)
		} else {
			err := CopyFile(post.dir+"thumb.png", output)
			if err != nil {
				app.errorLog.Fatalf("Failed to copy image %s\n%s\n", post.filename, err)
			}
		}
		//TODO test
	case fileExists(post.dir + "thumb.jpg"):
		images.ResizeJpegToPng(post.dir+"thumb.jpg", output, 0, 200)
		app.infoLog.Printf("Created thumb %s", output)
	case fileExists(post.dir + "thumb.jpeg"):
		images.ResizeJpegToPng(post.dir+"thumb.jpeg", output, 0, 200)
		app.infoLog.Printf("Created thumb %s", output)
	default:
		app.errorLog.Printf("Failed to open image for %s", post.filename)
		return post
	}
	app.infoLog.Printf("Created thumb %s", output)
	post.Thumb = "/images/thumbs/" + strings.Split(post.filename, ".")[0] + ".png"
	return post
}

func (app *application) logPostdata(field, data, filename string) {
	if data == "" {
		app.warnLog.Printf("No %s data for %s", field, filename)
	} else {
		app.infoLog.Printf("%s %s: %s", filename, field, data)
	}
}

func (app *application) makeOutputDir(dir string) {
	_, err := os.Stat(app.config.Output + dir)

	if errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(app.config.Output+dir, os.ModePerm); err != nil {
			app.errorLog.Fatalf("Failed to create uput dir %s\n%s", dir, err)
		}
	}
}
