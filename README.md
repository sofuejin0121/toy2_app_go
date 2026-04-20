# Chirp - GoとReactで作るマイクロブログSNS

GoバックエンドとReactフロントエンドで構築したマイクロブログSNSアプリケーションです。

## 技術スタック

### バックエンド
- **言語**: Go
- **データベース**: SQLite（[sqlc](https://sqlc.dev/) でクエリ自動生成）
- **画像ストレージ**: Cloudflare R2
- **認証**: セッション管理 + メール認証（Resend）

### フロントエンド
- **フレームワーク**: React + TypeScript
- **ビルドツール**: Vite
- **スタイリング**: Tailwind CSS

### インフラ
- **バックエンド**: Render
- **フロントエンド**: Netlify

### Render のスリープ対策（10 分ごとのヘルスチェック）

Render の無料 Web サービスは、しばらくアクセスが無いとスピンダウンします。本リポジトリでは **GitHub Actions** が **約 10 分ごと**に `GET /api/health` へリクエストし、起動状態を保ちやすくしています（`.github/workflows/render-healthcheck.yml` の `schedule`）。

1. GitHub リポジトリを開く → **Settings** → **Secrets and variables** → **Actions** → **New repository secret**
2. **Name**: `RENDER_HEALTHCHECK_URL`
3. **Secret**: Render の本番 URLにパスを付けたもの（例）  
   `https://<サービス名>.onrender.com/api/health`  
   ※フロントを Netlify に置いていても、**ping 先は API を動かしている Render のオリジン**にしてください。
4. **Actions** タブでワークフロー「Render healthcheck」を選び、**Run workflow** から手動実行して成功することを確認する。

`cron` は **UTC** で動きます。長期間リポジトリに活動が無いと GitHub 側でスケジュール実行が止まることがあります（再び push 等で復帰）。無料枠では外部からの ping でも完全にスリープしない保証はなく、確実に常時起動したい場合は Render の有料プラン（Always On 等）の利用を検討してください。

## 機能

- ユーザー登録・ログイン（メール認証）
- マイクロポスト投稿・削除（画像添付対応）
- フォロー / フォロワー
- いいね・ブックマーク
- 通知
- 管理者機能
- パスワードリセット

## 開発環境のセットアップ

### バックエンド

依存パッケージをダウンロードします。

```
$ go mod tidy
```

サーバーを起動します。

```
$ go run cmd/server/main.go
```

テストを実行します。

```
$ go test ./...
```

### フロントエンド

```
$ cd frontend
$ npm install
$ npm run dev
```
