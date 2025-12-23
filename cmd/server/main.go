package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/matthew-jp2525/yt-summary-server/internal/config"
	"github.com/matthew-jp2525/yt-summary-server/internal/httpapi"
	"github.com/matthew-jp2525/yt-summary-server/internal/logger"
	"github.com/matthew-jp2525/yt-summary-server/internal/summarizer"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	logger.Init(cfg.Debug)

	httpapi.SetConfig(&cfg)
	summarizer.SetConfig(&cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/summarize", httpapi.SummarizeHandler)

	addr := ":" + cfg.Port

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 2 * time.Minute,
	}

	log.Printf("listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}
