package main

import (
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/brensch/recreation"
)

func HandleScrapeAvailabilities(client *recreation.Obfuscator, fs *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DoAvailabilitiesSync(r.Context(), client, fs, time.Now())
	}
}

func HandleScrapeGrounds(client *recreation.Obfuscator, fs *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DoGroundsSync(r.Context(), client, fs)
	}
}
