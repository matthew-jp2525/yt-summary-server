package httpapi

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/matthew-jp2525/yt-summary-server/internal/logger"
	"github.com/matthew-jp2525/yt-summary-server/internal/subtitle"
	"github.com/matthew-jp2525/yt-summary-server/internal/summarizer"
)

func SummarizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	url := strings.TrimSpace(r.FormValue("url"))
	if url == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Minute)
	defer cancel()

	info, err := subtitle.FetchVideoInfo(ctx, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug.Printf("video title: %s", info.Title)
	logger.Debug.Printf("cleaned video subtitle text: %s", info.Text)

	summary, err := summarizer.Summarize(ctx, info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug.Printf("summary: %s", summary)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(info.Title + "\n\n" + summary))
}
