package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/matthew-jp2525/yt-summary-server/internal/config"
	"github.com/matthew-jp2525/yt-summary-server/internal/httpapi"
	"github.com/matthew-jp2525/yt-summary-server/internal/logger"
	"github.com/matthew-jp2525/yt-summary-server/internal/subtitle"
	"github.com/matthew-jp2525/yt-summary-server/internal/summarizer"
)

func isAPIKeyAuthEnabled(env string, disable bool) bool {
	if env == "prod" && disable {
		log.Fatal("API key auth cannot be disabled in prod.")
	}

	if disable {
		log.Println("WARNING: API key auth is DISABLED")
		return false
	}

	return true
}

func apiKeyMiddleware(validKeys map[string]struct{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				http.Error(w, "missing api key", http.StatusUnauthorized)
				return
			}

			// 定数時間比較でチェック
			for k := range validKeys {
				if subtle.ConstantTimeCompare([]byte(key), []byte(k)) == 1 {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "invalid api key", http.StatusForbidden)
		})
	}
}

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	logger.Init(cfg.Debug)

	if cfg.Debug {
		logger.Debug.Println("Debug log is enabled")
	}

	httpapi.SetConfig(&cfg)
	subtitle.SetConfig(&cfg)
	summarizer.SetConfig(&cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/summarize", httpapi.SummarizeHandler)
	mux.HandleFunc("/transcript", httpapi.TranscriptHandler)

	authEnabled := isAPIKeyAuthEnabled(cfg.Env, cfg.DisableAPIKeyAuth)

	var handler http.Handler = mux

	if authEnabled {
		if cfg.APIKeys == nil {
			log.Fatal("API key auth enabled but APIKeys is not configured")
		}

		handler = apiKeyMiddleware(cfg.APIKeys)(handler)
	}

	addr := ":" + cfg.Port

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 2 * time.Minute,
	}

	if cfg.YTDLPCookiePath != nil {
		log.Printf("using cookies: %q", *cfg.YTDLPCookiePath)
	}

	if cfg.YTDLPUserAgent != nil {
		log.Printf("using user agent: %q", *cfg.YTDLPUserAgent)
	}

	if cfg.YTDLPJSRuntimes != nil {
		log.Printf("using js runtimes: %q", *cfg.YTDLPJSRuntimes)
	}

	log.Printf("listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}
