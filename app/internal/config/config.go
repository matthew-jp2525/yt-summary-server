package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port string

	GeminiAPIKey   string
	YTDLPPath      string
	YTDLCookiePath *string

	Debug bool
}

func Load() Config {
	var ytdlCookiePath *string
	maybeYTDLCookiePath := os.Getenv("YTDLP_COOKIE_PATH")

	if maybeYTDLCookiePath != "" {
		_, err := os.Stat(maybeYTDLCookiePath)
		if err == nil {
			ytdlCookiePath = &maybeYTDLCookiePath
		}
	}

	return Config{
		Port: getOr("PORT", "8080"),

		GeminiAPIKey:   mustGet("GEMINI_API_KEY"),
		YTDLPPath:      getOr("YTDLP_PATH", "yt-dlp"),
		YTDLCookiePath: ytdlCookiePath,

		Debug: getBool("DEBUG"),
	}
}

// ===== helpers =====

func mustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("%s no set", key))
	}
	return v
}

func getOr(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getBool(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		return false
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}
