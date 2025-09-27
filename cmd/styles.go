package main

import (
	"bufio"
	"os"

	"github.com/BurntSushi/toml"
)

type Css struct {
	Font  Font
	Dark  Dark
	Light Light
}

type Font struct {
	FontUrl         string
	HeadingFont     string
	BodyFont        string
	BaseFontSize    string
	H1              string
	H2              string
	H3              string
	H4              string
	H5              string
	SmallText       string
	ContentDivWidth string
}
type Dark struct {
	Background  string
	TextColor   string
	LineColor   string
	HeaderText  string
	CardBg      string
	CardText    string
	AnchorColor string
	NavColor    string
	NavText     string
	NavHover    string
	HeaderColor string
	CodeBlocks  string
	Code        string
	BorderColor string
	LinkColor   string
	LinkHover   string
	FooterColor string
}
type Light struct {
	Background  string
	TextColor   string
	LineColor   string
	HeaderText  string
	CardBg      string
	CardText    string
	AnchorColor string
	NavColor    string
	NavText     string
	NavHover    string
	HeaderColor string
	CodeBlocks  string
	Code        string
	BorderColor string
	LinkHover   string
	LinkColor   string
	FooterColor string
}

func (app *application) generateCssVarsFile() Css {
	output, err := os.Create(app.config.Output + "/static/css/vars.css")
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer output.Close()
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	css, err := parseStylesConfig(app.config.StylesConfig)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	lines := []string{
		":root {",
		buildVarString("headingFont", css.Font.HeadingFont),
		buildVarString("bodyFont", css.Font.BodyFont),
		buildVarString("baseFontSize", css.Font.BaseFontSize),
		buildVarString("h1", css.Font.H1),
		buildVarString("h2", css.Font.H2),
		buildVarString("h3", css.Font.H3),
		buildVarString("h4", css.Font.H4),
		buildVarString("h5", css.Font.H5),
		buildVarString("smallText", css.Font.SmallText),
		buildVarString("contentDivWidth", css.Font.ContentDivWidth),
		"}",
		"body.dark {",
		buildVarString("background", css.Dark.Background),
		buildVarString("textColor", css.Dark.TextColor),
		buildVarString("lineColor", css.Dark.LineColor),
		buildVarString("headerText", css.Dark.HeaderColor),
		buildVarString("cardBg", css.Dark.CardBg),
		buildVarString("cardText", css.Dark.CardText),
		buildVarString("anchorColor", css.Dark.AnchorColor),
		buildVarString("navColor", css.Dark.NavColor),
		buildVarString("navText", css.Dark.NavText),
		buildVarString("navHover", css.Dark.NavHover),
		buildVarString("headerColor", css.Dark.HeaderColor),
		buildVarString("codeBlocks", css.Dark.CodeBlocks),
		buildVarString("code", css.Dark.Code),
		buildVarString("borderColor", css.Dark.BorderColor),
		buildVarString("linkColor", css.Dark.LinkColor),
		buildVarString("linkHover", css.Dark.LinkHover),
		buildVarString("footerColor", css.Dark.FooterColor),
		"}",
		"body.light {",
		buildVarString("background", css.Light.Background),
		buildVarString("textColor", css.Light.TextColor),
		buildVarString("lineColor", css.Light.LineColor),
		buildVarString("headerText", css.Light.HeaderColor),
		buildVarString("cardBg", css.Light.CardBg),
		buildVarString("cardText", css.Light.CardText),
		buildVarString("anchorColor", css.Light.AnchorColor),
		buildVarString("navColor", css.Light.NavColor),
		buildVarString("navText", css.Light.NavText),
		buildVarString("navHover", css.Light.NavHover),
		buildVarString("headerColor", css.Light.HeaderColor),
		buildVarString("codeBlocks", css.Light.CodeBlocks),
		buildVarString("code", css.Light.Code),
		buildVarString("borderColor", css.Light.BorderColor),
		buildVarString("linkColor", css.Light.LinkColor),
		buildVarString("linkHover", css.Light.LinkHover),
		buildVarString("footerColor", css.Light.FooterColor),
		"}",
	}

	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			app.errorLog.Fatal(err)
		}
	}

	return css
}

func buildVarString(variable, value string) string {
	return " --" + variable + ": " + value + ";"
}

func parseStylesConfig(configPath string) (Css, error) {
	css := Css{}
	f := configPath

	if _, err := os.Stat(f); err != nil {
		return css, err
	}

	toml.DecodeFile(f, &css)
	return css, nil
}
