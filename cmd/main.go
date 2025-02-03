package main

import (
	//"fmt"
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

	temps, err = template.ParseGlob("views/*.html")
	if err != nil {
		log.Fatalf("Glob Parseing failed: %v", err)
	}

	//static
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", DynamicEntry)

	log.Println("Server live at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

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

	//Ignore requests with file extensions (like .css, .js, .ico)
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
