package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
)

type Post struct {
	Title       string
	Description string
	Date        time.Time
	Data        string
	Slug        string
}

var posts []Post

func blogHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).ParseFiles(
		"./templates/base.html",
		"./templates/home.html",
		"./templates/posts.html",
	))

	tmpl.ExecuteTemplate(w, "home", posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).ParseFiles(
		"./templates/base.html",
		"./templates/post.html",
	))

	slug := strings.TrimPrefix(r.URL.Path, "/post/")
	for _, post := range posts {
		if post.Slug == slug {
			fmt.Println("DEBUG:", post.Title, post.Data)
			tmpl.ExecuteTemplate(w, "post", post)
			return
		}
	}
	http.NotFound(w, r)
}

func loadContent() *[]Post {
	files, err := os.ReadDir("content")
	if err != nil {
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		content, err := os.ReadFile("content/" + file.Name())
		if err != nil {
			return nil
		}

		raw := string(content)

		body := raw

		var title string
		var description string
		var date time.Time

		if len(raw) > 0 && raw[:3] == "---" {
			parts := strings.SplitN(raw, "---", 3)
			if len(parts) >= 3 {
				frontMatter := parts[1]
				body = parts[2]
				for _, line := range strings.Split(frontMatter, "\n") {
					if strings.HasPrefix(line, "title:") {
						title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
					} else if strings.HasPrefix(line, "description:") {
						description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
					} else if strings.HasPrefix(line, "date:") {
						dateStr := strings.TrimSpace(strings.TrimPrefix(line, "date:"))
						date, err = time.Parse(time.RFC3339, dateStr)
						if err != nil {
							return nil // Handle error appropriately
						}
					}
				}
			}
		}

		mdContent := markdown.ToHTML([]byte(body), nil, nil)
		html := bluemonday.UGCPolicy().SanitizeBytes(mdContent)
		slug := strings.TrimSuffix(file.Name(), ".md")
		post := Post{
			Title:       title,
			Description: description,
			Date:        date,
			Data:        string(html),
			Slug:        slug,
		}
		posts = append(posts, post)
	}

	return &posts
}

func main() {
	postsList := loadContent()
	if postsList == nil {
		http.Error(http.ResponseWriter(nil), "Failed to load content", http.StatusInternalServerError)
		return
	}
	http.HandleFunc("/", blogHandler)
	http.HandleFunc("/post/", postHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}
