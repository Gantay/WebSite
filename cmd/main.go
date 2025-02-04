package main

import (
	//"fmt"
	"golang.org/x/text/cases"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
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

	//HTTPS with self-signed
	// log.Println("Server running at https://localhost:8443")
	// log.Fatal(http.ListenAndServeTLS(":8443", "cert/cert.pem", "cert/key.pem", mux))
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
		Title:       strings.Title(strings.Replace(path, "-", " ", -1)),
		CurrentYear: time.Now().Year(),
		Path:        r.URL.Path,
	}

	// Check if template exists
	if templ.Lookup(templateFile) == nil {
		templ.ExecuteTemplate(w, "404.html", nil)
		log.Printf("404 not found: %s", templateFile)
		return
	}

	// Render the requested template
	err := templ.ExecuteTemplate(w, templateFile, nil)
	if err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		log.Printf("Template execution failed: %v", err)
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
