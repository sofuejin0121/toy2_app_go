package mailer

import "github.com/sofuejin0121/toy_app_go/internal/model"

// Mailer はメール送信のインターフェイス
type Mailer interface {
	SendAccountActivation(user *model.User) error
	SendPasswordReset(user *model.User) error
	SendFollowNotification(to *model.User, follower *model.User) error
	SendLikeNotification(to *model.User, liker *model.User, content string) error
}

// SMTPMailer はSMTPを使ったメーラー実装
type SMTPMailer struct {
	Host     string // SMTPサーバーのホスト名
	Port     int    // SMTPサーバーのポート番号
	Username string // SMTP認証ユーザー名
	Password string // SMTP認証パスワード
	From     string // 送信元アドレス
	AppHost  string // アプリケーションのホスト名（例: "example.com"）。メール内のURLに使用
}

// LogMailer は開発環境用のメーラー (ログに出力)
type LogMailer struct {
	From string // 送信元アドレス
	Host string // ホスト名
}

// BrevoMailer は Brevo HTTP API を使ったメーラー実装（SMTPポート不使用）
type BrevoMailer struct {
	APIKey  string // Brevo API キー
	From    string // 送信元アドレス
	AppHost string // アプリケーションのホスト名
}

// EmailData はメールテンプレートに渡すデータ
type EmailData struct {
	User          *model.User
	ActivationURL string
	ResetURL      string
}

// デフォルトの送信元アドレス
const DefaultFrom = "noreply@example.com"
