package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/matthew-jp2525/yt-summary-server/internal/config"
	"github.com/matthew-jp2525/yt-summary-server/internal/logger"
	"github.com/matthew-jp2525/yt-summary-server/internal/subtitle"
	"github.com/matthew-jp2525/yt-summary-server/internal/summarizer"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: <command> <youtube_url>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	url := flag.Arg(0)

	_ = godotenv.Load()

	cfg := config.Load()
	logger.Init(cfg.Debug)
	summarizer.SetConfig(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	info, err := subtitle.FetchVideoInfo(ctx, url)
	if err != nil {
		logger.Error.Printf("failed: %v", err)
		os.Exit(1)
	}

	summary, err := summarizer.Summarize(ctx, info)
	if err != nil {
		logger.Error.Printf("failed: %v", err)
		os.Exit(1)
	}

	fmt.Println(info.Title + "\n\n" + summary)
}
