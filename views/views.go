package views

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func Render(resp http.ResponseWriter, layout string, data interface{}) {
	t := template.New("")
	t.Funcs(template.FuncMap{})
	templates, err := filepath.Glob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	templates = append(templates, layout)
	tmp, err := t.ParseFiles(templates...)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmp.ExecuteTemplate(resp, "main", data); err != nil {
		http.Error(resp, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
