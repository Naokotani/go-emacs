package main

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

func (app *application) generatePages() error {
	for key, ts := range app.pagesTemplateCache {
		filename := strings.Split(key, ".")[0] + ".html"
		path := app.outputPath + "/" + filename
		output, err := os.Create(path)
		if err != nil {
			return err
		}
		defer output.Close()
		fmt.Printf("template: %s\n", filename)
		err = ts.ExecuteTemplate(output, "base", app.config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *application) generatePosts() error {
	postFiles, err := filepath.Glob("./ui/html/pages/posts/html/*.html")
	if err != nil {
		return err
	}
	_, err = os.Stat(app.outputPath + "posts")

	if errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(app.outputPath+"/posts", os.ModePerm); err != nil {
			return err
		}
	}

	for _, post := range postFiles {
		name := filepath.Base(post)
		fmt.Printf("Post name: %s\n", name)
		path := app.outputPath + "/posts/" + name
		output, err := os.Create(path)
		if err != nil {
			return err
		}
		defer output.Close()

		html, err := os.ReadFile(post)
		if err != nil {
			return err
		}
		app.config.Page.Post = template.HTML(html)

		ts, err := getPostTemplate()
		if err != nil {
			return err
		}

		err = ts.ExecuteTemplate(output, "base", app.config)
	}

	return nil
}

func getPostTemplate() (*template.Template, error) {

	ts, err := template.ParseFiles("./ui/html/base.gotmpl")
	if err != nil {
		return nil, err
	}

	ts, err = ts.ParseFiles("./ui/html/templates/post.gotmpl")
	if err != nil {
		return nil, err
	}

	ts, err = ts.ParseGlob("./ui/html/partials/*.gotmpl")
	if err != nil {
		return nil, err
	}

	return ts, nil
}
