package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"

	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
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

// handlerDisplayEvent for post event
func (s server) handlerDisplayEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}

		var ms libmetier.MessageSocial
		err = json.Unmarshal(body, &ms)
		if err != nil {
			log.Printf("parse error")
			log.Println(err)
		} else {
			var message pubsub.PubsubMessage
			b, err := json.Marshal(ms)
			if err != nil {
				log.Println(err)
			}
			message.Data = []byte(b)

			s.publishmessage(&message)
			log.Printf("POST done")
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
