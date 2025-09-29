package main

import (
	"errors"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/naokotani/go-emacs/internal/images"
)

type View struct {
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
	Css       Css
	Pages     []Page
	Tags      map[string][]Post
}

type Page struct {
	Title   string
	Slug    string
	src     string
	dst     string
	dirName string
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

func (app *application) generateViews(css Css) {
	pages := app.getPagesData()

	for _, page := range pages {
		app.infoLog.Printf("Loaded page data for %s\n", page.Title)
	}

	page := View{
		Title:     app.config.Site.Title,
		SubHeader: app.config.Site.SubHeader,
		Contact:   app.config.Contact,
		Date:      time.Now(),
		Site:      app.config.Site,
		Resume:    app.config.Resume,
		Pages:     pages,
		Css:       css,
	}

	generatePosts(app, page)
	generatePages(app, page)
	generateResume(app, page)
	page.Tags = app.config.Site.Tags
	generateIndex(app, page)

	for tag, posts := range page.Tags {
		page.Site.Posts = posts
		app.infoLog.Printf("Generating tag page for: %s ", tag)
		generateTagHome(app, page, tag)
	}
}

func generateIndex(app *application, view View) {
	view.Site.Posts = app.config.Site.Posts

	output, err := os.Create(app.config.Output + "index.html")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer output.Close()

	ts, ok := app.templateCache["index.gotmpl"]
	if !ok {
		app.errorLog.Fatal("Template index.html does not exist in the cache")
	}

	err = ts.ExecuteTemplate(output, "base", view)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func generateTagHome(app *application, view View, tag string) {
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

	err = ts.ExecuteTemplate(output, "base", view)
}

func generateResume(app *application, view View) {
	resume := app.config.Resume.Html
	if resume == "" {
		app.warnLog.Printf("Resume string empty, skipping page. Add to config.toml to generate resume.\n")
		return
	}
	app.infoLog.Printf("Generating resume from %s\n", resume)

	app.makeOutputDir("resume")
	output, err := os.Create(app.config.Output + "resume/resume.html")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer output.Close()

	html, err := os.ReadFile(resume)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	view.Content = template.HTML(html)

	ts, ok := app.templateCache["resume.gotmpl"]
	if !ok {
		app.errorLog.Fatal("Template resume.gotmpl does not exist in the cache")
	}

	err = ts.ExecuteTemplate(output, "base", view)
	if err != nil {
		app.errorLog.Fatal(err)
	}
}

func (app *application) getPagesData() []Page {
	dirs, err := os.ReadDir(app.config.Pages.Dir)

	if err != nil {
		app.errorLog.Fatalf("Failed to read pages directory: %s\n", app.config.Pages.Dir)
	}

	var pages []Page

	for _, dir := range dirs {
		page := Page{
			Title:   dir.Name(),
			Slug:    "/" + dir.Name() + "/" + dir.Name() + ".html",
			src:     app.config.Pages.Dir + dir.Name() + "/" + dir.Name(),
			dst:     app.config.Output + dir.Name() + "/" + dir.Name(),
			dirName: dir.Name(),
		}
		page = getPageMetadata(app, page, app.config.Pages.Dir+dir.Name()+"/metadata.toml")
		pages = append(pages, page)
	}
	return pages
}

func generatePages(app *application, view View) {
	for _, page := range view.Pages {
		app.makeOutputDir(page.dirName)
		writePage(app, view, page.dst+".html", page.src+".html")
		err := copyDirectory(app.config.Pages.Dir+page.dirName+"/images/", app.config.Output+page.dirName+"/images/")
		if err != nil {
			app.errorLog.Fatalf("Failed to copy images for page: %s\n%s\n", page.dirName, err)
		}
	}

}

func writePage(app *application, view View, dst, srcFile string) {
	output, err := os.Create(dst)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer output.Close()

	html, err := os.ReadFile(srcFile)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	view.Content = template.HTML(html)

	ts, ok := app.templateCache["page.gotmpl"]
	if !ok {
		app.errorLog.Fatal("Template page.gotmpl does not exist in the cache")
	}

	err = ts.ExecuteTemplate(output, "base", view)
	app.infoLog.Printf("Page %s written", srcFile)

}

func generatePosts(app *application, view View) {
	posts := writePostHtml(app, view)

	app.config.Site.Tags = make(map[string][]Post)
	app.config.Site.Posts = posts
	var emptyTags []string
	app.config.Site.Tags, emptyTags = buildTagMap(posts)
	app.warnLog.Printf("Posts with empty tag string: %s", emptyTags)

	for _, post := range posts {
		err := copyDirectory(post.dir+"/images", app.config.Output+"posts/images/")
		if err != nil {
			app.errorLog.Fatalf("Posts images failed to copy for: %s\n%s\n", post.filename, err)
		}
	}
}

func buildTagMap(posts []Post) (map[string][]Post, []string) {
	var emptyTags []string
	pMap := make(map[string][]Post)
	for _, p := range posts {
		for _, t := range p.Tags {
			if t != "" {
				pMap[t] = append(pMap[t], p)
			} else {
				emptyTags = append(emptyTags, p.filename)
			}
		}
	}
	return pMap, emptyTags
}

func writePostHtml(app *application, view View) []Post {
	var posts []Post
	for _, post := range app.config.Site.Posts {
		post := getPostMetadata(app, post)
		post = app.createThumb(post)

		view.Title = post.Title
		view.Post = post

		output, err := os.Create(app.config.Output + "posts/" + post.filename)
		if err != nil {
			app.errorLog.Fatal(err)
		}
		defer output.Close()

		html, err := os.ReadFile(post.dir + post.filename)
		if err != nil {
			app.errorLog.Fatal(err)
		}
		view.Content = template.HTML(html)

		postFile := "post.gotmpl"
		ts, ok := app.templateCache[postFile]
		if !ok {
			app.errorLog.Fatalf("Template %s does not exist in the cache", postFile)
		}

		err = ts.ExecuteTemplate(output, "base", view)
		if err != nil {
			app.errorLog.Fatal(err)
		}
		posts = append(posts, post)
	}
	return posts
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

func (app *application) makeOutputDir(dir string) {
	_, err := os.Stat(app.config.Output + dir)

	if errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(app.config.Output+dir, os.ModePerm); err != nil {
			app.errorLog.Fatalf("Failed to create uput dir %s\n%s", dir, err)
		}
	}
}
