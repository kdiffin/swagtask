package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"swagtask/database"
	"swagtask/pages"
	"swagtask/router"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

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

func main() {
	e := echo.New()
	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	e.Static("/images", "../images")
	e.Static("/css", "../css")
	e.Static("/js", "../js")

	dbpool := database.DatabaseInit()
	// close db pool when the funciton main ends
	defer dbpool.Close()

	// todo: reimplement
	router.Tasks(e, dbpool)
	// router.Tags(e, dbpool)
	e.GET("/", func(c echo.Context) error {
		page := pages.IndexPage{
			Tasks: []database.TaskWithTags{},
		}
		tagsWithTasks, err := database.GetAllTasksWithTags(dbpool)
		if err != nil {
			return c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		page.Tasks = tagsWithTasks
		return c.Render(200, "index", page)
	})

	e.Logger.Fatal(e.Start(":4000"))
}
