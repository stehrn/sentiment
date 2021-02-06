package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

// Analyse takes the the JSON encoded "message" field in the body
// of the request and analyzes the sentiment using the Cloud Natural Language API
func Analyse(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		log.Printf("json.NewDecoder: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if d.Message == "" {
		log.Printf("Message was empty")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sentiment, err := scoreSentiment(d.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Sentiment is: %s", sentiment.String())
	fmt.Fprint(w, html.EscapeString(sentiment.String()))
}
