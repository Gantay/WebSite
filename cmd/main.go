package main

import (
	"html/template"
	"log"

	"net/http"
)

func main() {
	//css
	css := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", css))
	//html
	tmpl, err := template.ParseGlob("views/*.html")
	if err != nil {
		log.Fatalln("ParseFiles faile")
	}

	h1 := func(w http.ResponseWriter, _ *http.Request) {
		tmpl.Execute(w, nil)
	}

	http.HandleFunc("/hello", h1)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func helloHandleFunc(w http.ResponseWriter, r *http.Request) {
// 	tmp.Execute(w, nil)
// 	// fmt.Fprint(w, "Hello world")
// }
