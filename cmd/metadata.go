package main

import (
	"github.com/BurntSushi/toml"
	"strings"
	"time"
)

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

func (app *application) logPostdata(field, data, filename string) {
	if data == "" {
		app.warnLog.Printf("No %s data for %s", field, filename)
	} else {
		app.infoLog.Printf("%s %s: %s", filename, field, data)
	}
}

func getPageMetadata(app *application, page Page, file string) Page {
	app.infoLog.Printf("Reading page metadata for: %s\n", page.dirName)
	if !fileExists(file) {
		app.errorLog.Fatalf("Metadata file does not exist in %s", page.dirName)
		return page
	}

	toml.DecodeFile(file, &page)

	app.infoLog.Printf("Page metadata for %s loaded. Title: %s", page.dirName, page.Title)

	return page
}
