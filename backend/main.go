package main

import (
	"html/template"
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
    tmpl *template.Template
}

func newTemplate() *Template {
    return &Template{
        tmpl: template.Must(template.ParseGlob("../views/*.html")),
    }
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.tmpl.ExecuteTemplate(w, name, data)
}
 
type Contact struct {
    Name string
    Email string
    Id int
}
type Contacts = []Contact
func newContact(name string, email string,id int) Contact {
    return Contact{
        Name: name,
        Email: email,
        Id: id,
    }
}


type PageData struct {
    Contacts Contacts
}
func newPage() PageData {
    return PageData{
        Contacts: []Contact{
            newContact("Oglan", "apar@gmail.com", 1),
            newContact("Celil", "ok@gmail.com", 2),
        },
    }
}
func (p *PageData) hasEmail(email string) bool {
    for _, contact := range p.Contacts {
        if contact.Email == email {
            return true
        }
    }

    return false
}

func main() {

    e := echo.New()

    e.Renderer = newTemplate()
    e.Use(middleware.Logger())
    e.Static("/images", "../images")
    e.Static("/css", "../css")
    e.Static("/js", "../js")


    page := newPage()
    id := 3

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", page)
    });
    
    e.POST("/create-contact", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")

        contact := Contact{
            Name: name,
            Email: email,
            Id: id,
        }
        if (page.hasEmail(email)) {
            blockData := struct {ErrorText string}{ErrorText: "STOP ERRORING"}
            return c.Render(422, "form-error", blockData)
        } else {
            c.Render(200, "form-success","")
        }

        page.Contacts = append([]Contact{newContact(name,email,id)}, page.Contacts...)
        id++

        return c.Render(200, "contact", contact)
    });

    e.DELETE("/contacts/:id",  func(c echo.Context) error {
        time.Sleep(2 * time.Second)
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)

        if err != nil {
            return c.String(400, "Id must be an integer")
        }

        deleted := false
        for i, contact := range page.Contacts {
            if contact.Id == id {
                page.Contacts = append(page.Contacts[:i], page.Contacts[i+1:]...)
                deleted = true
                break
            }
        }
        if !deleted {
            return c.String(400, "Contact not found")
        }

        return c.NoContent(200)
    });

    e.Logger.Fatal(e.Start(":42069"))
}