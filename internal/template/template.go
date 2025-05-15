package template

import (
	"html/template"
	"io"
	"swagtask/internal/utils"
)

// Template struct holds parsed templates
type Template struct {
	tmpl *template.Template
}

// NewTemplate parses all templates and returns a Template
func NewTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("../views/**.gohtml")), // fixed your glob pattern too
	}
}

func (t *Template) Render(w io.Writer, name string, data any) error {
	err := t.tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		utils.LogError("Error rendering template:", err)
	}
	return nil
}
