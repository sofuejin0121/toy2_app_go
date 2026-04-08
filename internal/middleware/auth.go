package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/store"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type contextKey string

const (
	CurrentUserKey        contextKey = "currentUser"       // コンテキスト内のユーザーキー
	csrfTokenKey          contextKey = "csrfToken"         // コンテキスト内のCSRFトークンキー
	csrfCookieName                   = "csrf_token"        // ブラウザのcookie名
	signedCookieSeparator            = "--"                // 値と署名の区切り文字
	PermanentCookieExpiry            = 30 * 24 * time.Hour // 30日間を表す定数
)

var cookieSecret = []byte("sample-app-cookie-secret") // クッキーの秘密鍵

// Auth はセッションCookieからユーザーを解決し、CSRFトークンも扱うミドルウェア
// なぜ3層のネストになるのか？
// Auth(store) storeを受け取る(DBアクセスに必要)
// func(next) Handler 次のハンドラーを受け取る (ミドルウェアチェーン)
// HandlerFunc 実際のリクエスト処理
func Auth(s *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// ① CSRFトークンをコンテキストに入れる
			csrfToken := getOrCreateCSRFToken(w, r)
			ctx = context.WithValue(ctx, csrfTokenKey, csrfToken)
			// ② 書き込み系リクエストはCSRFトークンを照合する
			switch r.Method {
			case http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete:
				if r.FormValue("csrf_token") != csrfToken {
					http.Error(w, "Invalid CSRF token", http.StatusForbidden) // 不一致なら即拒否
					return
				}
			}
			// ③ セッションCookieからユーザーIDを取得してDBで検索
			userIDStr := GetSessionValue(r, "user_id")
			if userIDStr != "" {
				id, err := strconv.ParseInt(userIDStr, 10, 64)
				if err == nil {
					user, err := s.GetUser(id)
					if err == nil {
						ctx = context.WithValue(ctx, CurrentUserKey, user)
					}
				}
			}
			// ④ 更新されたコンテキストを持つリクエストを次のハンドラーへ渡す
			// このコンテキストにはCSRFトークンとユーザー情報が含まれている
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

		})
	}
}

// CurrentUser はリクエストコンテキストから現在のユーザーを取得します
func CurrentUser(r *http.Request) *model.User {
	user, _ := r.Context().Value(CurrentUserKey).(*model.User)
	return user
}

// CSRFTokenFromContext はリクエストコンテキストからCSRFトークンを取得する
func CSRFTokenFromContext(r *http.Request) string {
	token, _ := r.Context().Value(csrfTokenKey).(string)
	return token
}

// GetCSRFToken はCSRF用Cookieに保存されたトークンを取得
func GetCSRFToken(r *http.Request) string {
	return GetCookieValue(r, csrfCookieName)
}

// SetSessionValue は署名付きCookieにセッション値を保存する
func SetSessionValue(w http.ResponseWriter, key, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    signValue(value),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// GetSessionValue は署名付きセッションCookie値を検証して返す
func GetSessionValue(r *http.Request, key string) string {
	value := GetCookieValue(r, key)
	if value == "" {
		return ""
	}
	plain, ok := verifySignedValue(value)
	if !ok {
		return ""
	}
	return plain
}

// ClearSessionCookie は一時セッション用Cookieを削除する
func ClearSessionCookie(w http.ResponseWriter) {
	DeleteCookie(w, "user_id")
	DeleteCookie(w, "csrf_token")
}

// SignUserID はユーザーIDを署名付き文字列に変換
func SignUserID(userID int64) string {
	return signValue(strconv.FormatInt(userID, 10))
}

// VerifyUserID は署名付きユーザーIDを検証して数値に戻します
func VerifyUserID(signed string) (int64, bool) {
	plain, ok := verifySignedValue(signed)
	if !ok {
		return 0, false
	}
	id, err := strconv.ParseInt(plain, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

// GetCookieValue は指定された名前のCookie値を取得
func GetCookieValue(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// DeleteCookie は指定したCookieを削除する
func DeleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// 署名を作る
func signValue(value string) string {
	mac := hmac.New(sha256.New, cookieSecret) // 秘密鍵 + SHA256 でHMACを初期化
	mac.Write([]byte(value))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil)) // 値を入力
	return value + "--" + sig
}

// 署名を検証する
// hmac.Equal で定数時間比較をすることで、時間のかかるブルートフォース攻撃を防ぐ
func verifySignedValue(signed string) (string, bool) {
	parts := strings.SplitN(signed, "--", 2) // "5"と"AbCdEf"を分割
	value, sig := parts[0], parts[1]

	// 受け取ったvalueで署名を再計算
	mac := hmac.New(sha256.New, cookieSecret)
	mac.Write([]byte(value))
	expected := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	// hmac.Equal で定数時間比較
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return "", false // 一致しなければfalseを返す
	}
	return value, true
}
func getOrCreateCSRFToken(w http.ResponseWriter, r *http.Request) string {
	// すでにCookieにあればそれを使う(毎回新規生成しない)
	if token := GetCookieValue(r, csrfCookieName); token != "" {
		return token
	}
	// なければ暗号論的乱数で32バイト生成
	b := make([]byte, 32)
	rand.Read(b)
	token := base64.URLEncoding.EncodeToString(b) // URL安全なBase64文字列に変換
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return token
}