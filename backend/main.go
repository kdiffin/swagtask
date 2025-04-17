package main

import (
	"html/template"
	"io"

	"myapp/database"
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

	database := database.NewDatabase()
	router.Tasks(e, &database)
	router.Tags(e, &database)
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", database)
	})

	e.Logger.Fatal(e.Start(":4000"))
}
