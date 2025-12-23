package subtitle

import (
	"regexp"
	"strings"
)

const maxTextSize = 200_000

var htmlTagRegexp = regexp.MustCompile(`<[^>]+>`)

func cleanVTT(vtt string) string {
	var b strings.Builder
	var prevLine string

	lines := strings.SplitSeq(vtt, "\n")

	for line := range lines {
		if strings.HasPrefix(line, "WEBVTT") {
			continue
		}
		if strings.HasPrefix(line, "Kind: ") {
			continue
		}
		if strings.HasPrefix(line, "Language: ") {
			continue
		}
		if strings.Contains(line, "-->") {
			continue
		}

		line = htmlTagRegexp.ReplaceAllString(line, "")
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if line == prevLine {
			continue
		}
		prevLine = line

		b.WriteString(line)
		b.WriteByte(' ')
	}

	result := b.String()
	if len(result) > maxTextSize {
		result = result[:maxTextSize]
	}
	return result
}
