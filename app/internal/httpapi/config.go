package httpapi

import (
	"github.com/matthew-jp2525/yt-summary-server/internal/config"
)

var cfg *config.Config

func SetConfig(c *config.Config) {
	cfg = c
}
