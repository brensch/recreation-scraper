package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/brensch/recreation"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

var (
	log *zap.Logger
)

func main() {
	flag.Parse()

	logConfig := zap.NewProductionConfig()
	logConfig.Level.SetLevel(zap.DebugLevel)
	// this ensures google logs pick things up properly
	logConfig.EncoderConfig.MessageKey = "message"
	logConfig.EncoderConfig.LevelKey = "severity"
	logConfig.EncoderConfig.TimeKey = "time"
	// logConfig.Encoding = "console"

	// init logger
	log, err := logConfig.Build()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	ctx := context.Background()
	rec := recreation.InitServer(ctx, log.With(zap.String("module", "recreation_api")), 750*time.Millisecond)

	// init firestore
	fs, err := InitFirestore(ctx)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	// init collection for ground summary
	err = InitGroundsSummary(ctx, fs)
	if err != nil {
		panic(err)
	}

	log.Info("starting server")
	http.HandleFunc("/avails", HandleScrapeAvailabilities(log, rec, fs))
	http.HandleFunc("/grounds", HandleScrapeGrounds(log, rec, fs))

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Info("using default port", zap.String("port", port))
	}

	// Start HTTP server.
	log.Info("listening", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("server had error", zap.Error(err))
	}

}

func InitFirestore(ctx context.Context) (*firestore.Client, error) {

	var app *firebase.App

	// check if creds exist
	_, err := os.Stat("creds.json")
	if err != nil {
		// use default creds if not (means we're in the cloud)
		app, err = firebase.NewApp(context.Background(), nil)
	} else {
		// use file if creds.json does exist
		opt := option.WithCredentialsFile("creds.json")
		conf := &firebase.Config{ProjectID: "campr-app"}
		app, err = firebase.NewApp(context.Background(), conf, opt)
	}
	if err != nil {
		return nil, err
	}

	return app.Firestore(ctx)

}
