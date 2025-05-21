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
	t := template.Must(template.ParseGlob("./web/views/components/**.gohtml"))
	template.Must(t.ParseGlob("./web/views/**.gohtml"))
	return &Template{
		tmpl: t,
	}
}

func (t *Template) Render(w io.Writer, name string, data any) error {
	err := t.tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		utils.LogError("Error rendering template:", err)
	}
	return nil
}
