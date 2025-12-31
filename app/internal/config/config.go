package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port string

	GeminiAPIKey    string
	YTDLPPath       string
	YTDLPCookiePath *string
	YTDLPUserAgent  *string

	Debug bool
}

func Load() Config {
	var ytdlpCookiePath *string
	maybeYTDLPCookiePath := os.Getenv("YTDLP_COOKIE_PATH")

	if maybeYTDLPCookiePath != "" {
		_, err := os.Stat(maybeYTDLPCookiePath)
		if err == nil {
			ytdlpCookiePath = &maybeYTDLPCookiePath
		}
	}

	var ytdlpUserAgent *string
	maybeYTDLPUserAgent := os.Getenv("YTDLP_USER_AGENT")

	if maybeYTDLPUserAgent != "" {
		ytdlpUserAgent = &maybeYTDLPUserAgent
	}

	return Config{
		Port: getOr("PORT", "8080"),

		GeminiAPIKey:    mustGet("GEMINI_API_KEY"),
		YTDLPPath:       getOr("YTDLP_PATH", "yt-dlp"),
		YTDLPCookiePath: ytdlpCookiePath,
		YTDLPUserAgent:  ytdlpUserAgent,

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
