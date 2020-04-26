package utils

import (
	"html/template"
	"net/http"
)

var templates *template.Template

//LoadTemplates Loads templates
func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

//ExecuteTemplate executes a template
func ExecuteTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}
