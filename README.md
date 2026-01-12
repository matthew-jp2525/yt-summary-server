# yt-summary-server

YouTube の URL を受け取り、  
字幕を取得・整形したうえで **要約文または文字起こし全文** を生成し、  
**プレーンテキストで返す**シンプルな HTTP サーバーおよび CLI ツールです。

curl や iOS ショートカット、  
あるいは ChatGPT / Gemini などの LLM チャット UI と組み合わせて使う  
個人利用向けのツールを想定しています。

---

## 概要

このツールは、YouTube URL を入力として以下の処理を行います。

### summarize（要約）

1. YouTube URL を受け取る
2. yt-dlp で字幕を取得（日本語優先）
3. 字幕をプレーンテキストに整形 
4. Gemini API で要約文を生成
5. 要約結果をそのまま返す（text/plain）

### transcript（文字起こし）

1. YouTube URL を受け取る
2. yt-dlp で字幕を取得（日本語優先）
3. 字幕をプレーンテキストに整形 
4. 動画タイトル + 空行 + 字幕全文 を返す（text/plain）

動画ファイルや字幕は永続化せず、  
処理後はすべて破棄されます。

---

## 想定ユースケース

- 動画を「観る前に」内容を把握したい
- 通勤中に動画リンクを投げて要約だけ読む
- 文字起こし全文を LLM に貼り付けて独自に要約・分析したい
- 要約文や文字起こしをそのままメモアプリに保存したい
- VPS 上で常駐させ、個人用の知的フィルターとして使う

---

## 依存関係

- Go 1.25 以上
- yt-dlp（PATH に存在すること） 
- Gemini API Key（Google AI Studio）

---

## セットアップ

### 1. 環境変数

`.env.example` を参考に `.env` を作成してください。

```
GEMINI_API_KEY=your_api_key_here
```

---

### 2. 起動（HTTP サーバー）

```bash
go run cmd/server/main.go
```

デフォルトで `:8080` にバインドします。

---

### 3. 起動（CLI）

HTTP サーバーを立てずに、  
コマンドラインから直接実行することもできます。

#### 要約（summarize）

```bash
go run cmd/cli/main.go summarize https://www.youtube.com/watch?v=xxxx
```

#### 文字起こし（transcript）

```bash
go run cmd/cli/main.go transcript https://www.youtube.com/watch?v=xxxx
```

いずれも、標準出力に結果がそのまま表示されます。  
パイプやリダイレクトと組み合わせての利用を想定しています。

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
- 本文: 動画タイトル、空行、要約文

```
動画タイトル

この動画では〜について解説しており、主なポイントは次の通りです……
```

---

### POST /transcript

#### リクエスト

- Content-Type: application/x-www-form-urlencoded
- パラメータ:
  - url: YouTube の URL 

```bash
curl -X POST \
  -d "url=https://www.youtube.com/watch?v=xxxx" \
  http://localhost:8080/transcript
```

#### レスポンス

- Content-Type: text/plain
- 本文: 動画タイトル、空行、字幕全文

```
動画タイトル

ここに字幕の全文が続きます……
```

エラー時も、人間が読めるプレーンテキストを返します。

---

## 設計方針

- 人間が直接叩きやすいインターフェースを優先
- JSON による過剰な構造化は行わない
- 一次資料（transcript）と二次資料（summary）を明確に分離
- ステートレス（永続データを持たない）
- 外部依存（字幕取得・要約生成）は差し替え可能な構造にする

---

## 注意事項

- Gemini API の利用料金に注意してください
- 字幕が存在しない動画では失敗します
- 非公開・有料動画は対象外です
- 字幕が非常に長い場合、要約精度が落ちることがあります
