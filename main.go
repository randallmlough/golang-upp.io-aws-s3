package main

import (
	"dropzone-s3/upload"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func rootHandler(resp http.ResponseWriter, req *http.Request) {

	render(resp, nil)
}

func render(resp http.ResponseWriter, data interface{}) {
	t := template.New("")
	t.Funcs(template.FuncMap{})
	temps, err := filepath.Glob("layouts/*.html")
	if err != nil {
		log.Fatal(err)
	}
	tmp, err := t.ParseFiles(temps...)
	if err != nil {
		log.Fatal(err)
	}
	if err := tmp.ExecuteTemplate(resp, "main", data); err != nil {
		http.Error(resp, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/sign", upload.PreSignRequest)
	log.Println("listening on port :3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
