# yt-summary-server

YouTube の URL を受け取り、  
字幕を取得・整形したうえで Gemini による要約文を生成し、  
**プレーンテキストで返す**シンプルな HTTP サーバーです。

curl や iOS ショートカットなどから気軽に叩く  
個人利用向けのツールを想定しています。

---

## 概要

このサーバーは、1 リクエストにつき以下の処理を行います。

1. YouTube URL を受け取る  
2. yt-dlp で字幕を取得（日本語優先）  
3. 字幕をプレーンテキストに整形  
4. Gemini API で要約文を生成  
5. 要約結果をそのまま返す（text/plain）

動画ファイルや字幕は永続化せず、  
処理後はすべて破棄されます。

---

## 想定ユースケース

- 動画を「観る前に」内容を把握したい  
- 通勤中に動画リンクを投げて要約だけ読む  
- 要約文をそのままメモアプリに保存したい  
- VPS 上で常駐させ、個人用の知的フィルターとして使う  

---

## 依存関係

- Go 1.22 以上  
- yt-dlp（PATH に存在すること）  
- Gemini API Key（Google AI Studio）  

---

## セットアップ

### 1. 環境変数

`.env.example` を参考に `.env` を作成してください。

    GEMINI_API_KEY=your_api_key_here

---

### 2. 起動（HTTP サーバー）

    go run cmd/server/main.go

デフォルトで `:8080` にバインドします。

---

### 3. 起動（CLI）

HTTP サーバーを立てずに、  
コマンドラインから直接要約することもできます。

    go run cmd/cli/main.go https://www.youtube.com/watch?v=xxxx

標準出力に要約結果がそのまま表示されます。

---

## 使い方（HTTP）

### POST /summarize

#### リクエスト

- Content-Type: application/x-www-form-urlencoded  
- パラメータ:
  - url: YouTube の URL  

```bash
curl -X POST \
  -d "url=https://www.youtube.com/watch?v=xxxx" \
  http://localhost:8080/summarize
```

#### レスポンス

- Content-Type: text/plain  
- 本文: 要約文そのもの  

    この動画では〜について解説しており、主なポイントは次の通りです……

エラー時も人間が読めるテキストを返します。

---

## 設計方針

- 人間が直接叩きやすいインターフェースを優先  
- JSON による過剰な構造化は行わない  
- ステートレス（永続データを持たない）  
- 外部依存（字幕取得・要約生成）は差し替え可能な構造にする  

---

## 注意事項

- Gemini API の利用料金に注意してください  
- 字幕が存在しない動画では失敗します  
- 非公開・有料動画は対象外です  
- 字幕が非常に長い場合、要約精度が落ちることがあります  
