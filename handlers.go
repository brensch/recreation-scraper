package main

import (
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/brensch/recreation"
	"go.uber.org/zap"
)

func HandleScrapeAvailabilities(log *zap.Logger, rec *recreation.Server, fs *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DoAvailabilitiesSync(r.Context(), log, rec, fs, time.Now())
	}
}

func HandleScrapeGrounds(log *zap.Logger, rec *recreation.Server, fs *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DoGroundsSync(r.Context(), rec, fs)
	}
}
