package main

import (
	//"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// PageData holds data to be passed to templates
type PageData struct {
	Title       string
	Content     string
	Tag         string
	CurrentYear int
	Path        string
}

var templ *template.Template

func main() {

	mux := http.NewServeMux()

	// Initialize templates with functions
	var err error
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"title": cases.Title,
	}

	templ, err = template.New("").Funcs(funcMap).ParseGlob("views/*.html")
	if err != nil {
		log.Fatalf("Template parsing failed: %v", err)
	}

	//static
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", DynamicEntry)

	log.Println("Server live at http://localhost:8080")
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())

}

func DynamicEntry(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index"
	}

	templateFile := path + ".html"

	// Ignore requests for non-HTML files
	if ext := filepath.Ext(path); ext != "" && ext != ".html" {
		http.NotFound(w, r)
		return
	}

	// Prepare page data
	data := PageData{
		Title:       cases.Title(language.English).String(strings.Replace(path, "-", " ", -1)),
		CurrentYear: time.Now().Year(),
		Path:        r.URL.Path,
	}

	// Check template existence and render
	if templ.Lookup(templateFile) == nil {
		w.WriteHeader(http.StatusNotFound)
		if err := templ.ExecuteTemplate(w, "404.html", data); err != nil {
			log.Printf("Error rendering 404 page: %v", err)
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		log.Printf("404 not found: %s", templateFile)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render template
	if err := templ.ExecuteTemplate(w, templateFile, data); err != nil {
		log.Printf("Template execution failed for %s: %v", templateFile, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}

// Middleware for request logging
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

// Middleware for adding cache headers to static files
func addCacheHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		next.ServeHTTP(w, r)
	})
}
