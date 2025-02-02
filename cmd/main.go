package main

import (
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"

	"net/http"
)

var temps *template.Template

func main() {

	mux := http.NewServeMux()

	var err error

	//HTML  -load templates once at startup
	temps, err = template.ParseGlob("views/*.html")
	if err != nil {
		log.Fatalf("Glob Parseing failed: %v", err)
	}

	//CSS
	// fs := http.FileServer(http.Dir("static"))
	// mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// images := http.FileServer(http.Dir("images"))
	// mux.Handle("/images/", http.StripPrefix("/images/", images))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	mux.HandleFunc("/", DynamicEntry)

	fmt.Println("server listing on port 8080.")

	log.Fatal(http.ListenAndServe(":8080", mux))

}

func DynamicEntry(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index"
	}

	templateFile := path + ".html"

	// ✅ Ignore requests with file extensions (like .css, .js, .ico)
	if filepath.Ext(path) != "" {
		http.NotFound(w, r)
		return
	}

	// Check if template exists
	if temps.Lookup(templateFile) == nil {
		temps.ExecuteTemplate(w, "404.html", nil)
		log.Printf("404 not found: %s", templateFile)
		return
	}

	// Render the requested template
	err := temps.ExecuteTemplate(w, templateFile, nil)
	if err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		log.Printf("Template execution failed: %v", err)
	}

}
