package subtitle

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

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

func FetchAndClean(ctx context.Context, url string) (string, error) {
	if err := validateYoutubeURL(url); err != nil {
		return "", err
	}

	tmpDir, err := os.MkdirTemp("", "yt-sub")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	outTemplate := filepath.Join(tmpDir, "sub")

	cmd := exec.CommandContext(
		ctx,
		"yt-dlp",
		"--quiet",
		"--no-warnings",
		"--skip-download",
		"--write-auto-subs",
		"--sub-lang", "ja",
		"--sub-format", "vtt",
		"-o", outTemplate,
		url,
	)

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
