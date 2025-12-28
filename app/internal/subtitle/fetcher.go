package subtitle

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/matthew-jp2525/yt-summary-server/internal/config"
)

type VideoInfo struct {
	Title string
	Text  string
}

var cfg *config.Config

func SetConfig(c *config.Config) {
	cfg = c
}

func validateYoutubeURL(value string) error {
	u, err := url.ParseRequestURI(value)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	host := u.Hostname()

	if host != "youtu.bu" &&
		host != "youtube.com" &&
		host != "www.youtube.com" &&
		host != "m.youtube.com" {
		return fmt.Errorf("not a youtube host: %s", host)
	}

	return nil
}

func fetchTitle(ctx context.Context, url string) (string, error) {
	args := []string{
		"--quiet",
		"--no-warnings",
		"--print",
		"%(title)s",
	}

	if cfg.YTDLCookiePath != nil {
		args = append(args, "--cookies", *cfg.YTDLCookiePath)
	}

	args = append(args, url)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func fetchAndClean(ctx context.Context, url string) (string, error) {
	if err := validateYoutubeURL(url); err != nil {
		return "", err
	}

	tmpDir, err := os.MkdirTemp("", "yt-sub")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	outTemplate := filepath.Join(tmpDir, "sub")

	args := []string{
		"--quiet",
		"--no-warnings",
		"--skip-download",
		"--write-auto-subs",
		"--sub-lang", "ja",
		"--sub-format", "vtt",
		"-o", outTemplate,
	}

	if cfg.YTDLCookiePath != nil {
		args = append(args, "--cookies", *cfg.YTDLCookiePath)
	}

	args = append(args, url)

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	files, _ := filepath.Glob(filepath.Join(tmpDir, "*.vtt"))
	if len(files) == 0 {
		return "", errors.New("subtitle not found")
	}

	raw, err := os.ReadFile(files[0])
	if err != nil {
		return "", err
	}

	return cleanVTT(string(raw)), nil
}

func FetchVideoInfo(ctx context.Context, url string) (*VideoInfo, error) {
	title, err := fetchTitle(ctx, url)
	if err != nil {
		return nil, err
	}

	text, err := fetchAndClean(ctx, url)
	if err != nil {
		return nil, err
	}

	return &VideoInfo{
		Title: title,
		Text:  text,
	}, nil
}
