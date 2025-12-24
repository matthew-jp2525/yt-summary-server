package summarizer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/matthew-jp2525/yt-summary-server/internal/config"
	"github.com/matthew-jp2525/yt-summary-server/internal/subtitle"
)

const promptTemplate = `
# 指示
以下の Youtube 動画の字幕を要約してください。

# 条件
- 日本語で書く
- 内容の全体像が分かることを重視する
- 箇条書きではなく、読みやすい文章にする
- 話者の主張や論点が分かるようにする
- 細かすぎる冗長な部分は省く
- 4000文字程度でまとめる

# 出力
要約文のみ出力してください。

------

動画タイトル:
%s

字幕テキスト:
%s
`

var cfg *config.Config

func SetConfig(c *config.Config) {
	cfg = c
}

func Summarize(ctx context.Context, info *subtitle.VideoInfo) (string, error) {
	prompt := fmt.Sprintf(promptTemplate, info.Title, info.Text)

	payload := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", cfg.GeminiAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if len(parsed.Candidates) == 0 ||
		len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no response from gemini")
	}

	return parsed.Candidates[0].Content.Parts[0].Text, nil
}
