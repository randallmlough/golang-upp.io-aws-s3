package main

import (
	"dropzone-s3/images"
	"dropzone-s3/views"
	"log"
	"net/http"
)

func rootHandler(resp http.ResponseWriter, req *http.Request) {

	views.Render(resp, "templates/upload.html", nil)
}
func formHandler(resp http.ResponseWriter, req *http.Request) {

	views.Render(resp, "templates/form.html", nil)
}
func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/sign", images.PreSignRequest)
	http.HandleFunc("/images", images.GetImages)
	http.HandleFunc("/image", images.GetSignRequest)
	http.HandleFunc("/delete", images.DeleteImage)
	http.HandleFunc("/buckets", images.ListBuckets)
	log.Println("listening on port :3000...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
