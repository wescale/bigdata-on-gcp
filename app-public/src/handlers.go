package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

// handler http (index /)
func handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("WTF dude, error parsing your template.")

	}
	t.Execute(w, os.Getenv("TOPIC_NAME"))
	log.Println("Finished HTTP request at", r.URL.Path)
}
