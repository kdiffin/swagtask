package main

import (
	"html/template"
	"io"
	task_package "myapp/task"

	"myapp/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("../views/*.gohtml")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

// init code

func main() {
	e := echo.New()
	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	e.Static("/images", "../images")
	e.Static("/css", "../css")
	e.Static("/js", "../js")

	homePage := task_package.NewHomePage()

	router.Tasks(e, &homePage)
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", homePage)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
