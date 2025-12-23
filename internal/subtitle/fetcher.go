package subtitle

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

func FetchAndClean(ctx context.Context, url string) (string, error) {
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
