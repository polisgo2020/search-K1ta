package templates

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type Renderer struct {
	Templates *template.Template
}

func Init() (*Renderer, error) {
	tmpl, err := template.ParseGlob("server/templates/*.html")
	return &Renderer{tmpl}, err
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return r.Templates.ExecuteTemplate(w, name, data)
}
