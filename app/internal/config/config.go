package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port string

	APIKeys map[string]struct{}

	GeminiAPIKey    string
	YTDLPPath       string
	YTDLPCookiePath *string
	YTDLPUserAgent  *string
	YTDLPJSRuntimes *string

	Debug             bool
	Env               string
	DisableAPIKeyAuth bool
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

	var ytdlpJSRuntimes *string
	maybeYTDLPJSRuntimes := os.Getenv("YTDLP_JS_RUNTIMES")

	if maybeYTDLPJSRuntimes != "" {
		ytdlpJSRuntimes = &maybeYTDLPJSRuntimes
	}

	return Config{
		Port: getOr("PORT", "8080"),

		APIKeys: loadAPIKeys(getOr("YT_SUMMARY_SERVER_API_KEYS", "")),

		GeminiAPIKey:    mustGet("GEMINI_API_KEY"),
		YTDLPPath:       getOr("YTDLP_PATH", "yt-dlp"),
		YTDLPCookiePath: ytdlpCookiePath,
		YTDLPUserAgent:  ytdlpUserAgent,
		YTDLPJSRuntimes: ytdlpJSRuntimes,

		Debug:             getBool("DEBUG"),
		Env:               getOr("YT_SUMMARY_SERVER_ENV", "prod"),
		DisableAPIKeyAuth: getBool("DISABLE_API_KEY_AUTH"),
	}
}

func loadAPIKeys(env string) map[string]struct{} {
	keys := make(map[string]struct{})

	for _, k := range strings.Split(env, ",") {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		keys[k] = struct{}{}
	}

	if len(keys) == 0 {
		return nil
	}

	return keys
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
