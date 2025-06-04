package main

import (
	"html/template"
	"net/http"
)

type Post struct {
	Title   string
	Content string
}

var posts = []Post{}

var tmpl = template.Must(template.ParseFiles("./templates/layout.html", "./templates/post.html"))

func blogHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "layout", posts)
}

func main() {
	http.HandleFunc("/", blogHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}
