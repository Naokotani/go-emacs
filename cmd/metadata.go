package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func getPostMetadata(app *application, post Post) Post {
	f := filepath.Join(post.dir, "metadata.toml")
	if !fileExists(f) {
		app.errorLog.Fatalf("Metadata file does not exist in %s", post.dir)
		return post
	}

	toml.DecodeFile(f, &post)

	post.Tags = strings.Split(post.TagString, " ")
	post.TagString = strings.ReplaceAll(post.TagString, " ", " | ")
	post.Slug = "/posts/" + post.filename
	layout := "[2006-01-02 Mon 15:04]"
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

func (app *application) logPostdata(field, data, filename string) {
	if data == "" {
		app.warnLog.Printf("No %s data for %s", field, filename)
	} else {
		app.infoLog.Printf("%s %s: %s", filename, field, data)
	}
}

func getPageMetadata(app *application, page Page, file string) Page {
	app.infoLog.Printf("Reading page metadata for %s in %s\n", page.dirName, page.dst)
	if !fileExists(file) {
		app.errorLog.Fatalf("Metadata file does not exist in %s", page.dirName)
		return page
	}

	toml.DecodeFile(file, &page)

	app.infoLog.Printf("Page metadata for %s loaded. Title: %s", page.dirName, page.Title)

	return page
}

func (app *application) getPostDirs() error {
	postsDir := app.config.Posts.Dir

	enteries, err := os.ReadDir(postsDir)
	var posts []Post
	if err != nil {
		return err
	}
	for _, e := range enteries {
		if e.IsDir() {
			posts = append(posts, Post{
				dir:      filepath.Join(postsDir, e.Name()) + "/",
				filename: e.Name() + ".html",
			})
		}
	}

	app.config.Site.Posts = posts
	return nil
}

func (app *application) getResumeFiles() {
	resumeDir := app.config.Resume.Dir
	if !fileExists(filepath.Join(resumeDir, "resume.html")) {
		app.warnLog.Printf("No resume file found in %s. Add 'resume.html' to the resume directory to generate\n", app.config.Resume.Dir)
		app.config.Resume.IsResume = false
		return
	} else {
		app.config.Resume.IsResume = true
	}

	if fileExists(filepath.Join(resumeDir, app.config.Resume.Pdf)) {
		resumeSrc := filepath.Join(resumeDir, app.config.Resume.Pdf)
		resumeDst := filepath.Join(app.config.Output, "resume", app.config.Resume.Pdf)
		err := CopyFile(resumeSrc, resumeDst)
		if err != nil {
			app.errorLog.Fatalf("Failed to copy resume file from %s to %s\n", resumeSrc, resumeDst)
		}
	} else if app.config.Resume.Pdf == "" {
		app.infoLog.Print("No resume pdf file set. Skipping resume pdf")
	} else {
		app.warnLog.Printf("Resuem pdf set to %s, but file does not exist in %s",
			app.config.Resume.Pdf, app.config.Resume.Dir)
	}
}
