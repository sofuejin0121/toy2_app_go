package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sofuejin0121/toy_app_go/internal/mailer"
	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// LikeHandler はいいねリソースのHTTPハンドラー。
type LikeHandler struct {
	store  likeStore
	mailer mailer.Mailer
}

// likeStore は LikeHandler が必要とするストアメソッドを定義するインターフェース。
type likeStore interface {
	Like(userID, micropostID int64) error
	Unlike(userID, micropostID int64) error
	CountLikes(micropostID int64) (int, error)
	GetMicropost(id int64) (*model.Micropost, error)
	GetUser(id int64) (*model.User, error)
	GetOrCreateUserPreference(userID int64) (*model.UserPreference, error)
}

// NewLikeHandler は新しい LikeHandler を返す。
func NewLikeHandler(s likeStore, m mailer.Mailer) *LikeHandler {
	return &LikeHandler{store: s, mailer: m}
}

// isAJAX は X-Requested-With: XMLHttpRequest ヘッダーが付いているかを確認する。
func isAJAX(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// Create は POST /likes を処理する。
// フォームの micropost_id に対して現在のユーザーがいいねを作成する。
// AJAX リクエスト（X-Requested-With: XMLHttpRequest）の場合は JSON を返す。
func (h *LikeHandler) Create(w http.ResponseWriter, r *http.Request) {
	cu := currentUser(r)
	if cu == nil {
		if isAJAX(r) {
			http.Error(w, `{"error":"認証が必要です"}`, http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	micropostID, err := strconv.ParseInt(r.FormValue("micropost_id"), 10, 64)
	if err != nil || micropostID == 0 {
		http.Error(w, "不正なリクエストです", http.StatusBadRequest)
		return
	}

	if err := h.store.Like(cu.ID, micropostID); err != nil {
		log.Printf("Like(%d, %d): %v", cu.ID, micropostID, err)
	}

	// 投稿オーナーの通知設定を確認してメール送信
	if mp, err := h.store.GetMicropost(micropostID); err == nil && mp.UserID != cu.ID {
		if pref, err := h.store.GetOrCreateUserPreference(mp.UserID); err == nil && pref.EmailOnLike {
			if owner, err := h.store.GetUser(mp.UserID); err == nil {
				if err := h.mailer.SendLikeNotification(owner, cu, mp.Content); err != nil {
					log.Printf("SendLikeNotification: %v", err)
				}
			}
		}
	}

	count, _ := h.store.CountLikes(micropostID)

	if isAJAX(r) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"liked":true,"count":%d}`, count)
		return
	}

	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}

// Destroy は DELETE /likes/{micropost_id} を処理する。
// URL パスの {id} はいいね解除したいマイクロポストの ID。
// AJAX リクエストの場合は JSON を返す。
func (h *LikeHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	cu := currentUser(r)
	if cu == nil {
		if isAJAX(r) {
			http.Error(w, `{"error":"認証が必要です"}`, http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	micropostID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || micropostID == 0 {
		http.Error(w, "不正なリクエストです", http.StatusBadRequest)
		return
	}

	if err := h.store.Unlike(cu.ID, micropostID); err != nil {
		log.Printf("Unlike(%d, %d): %v", cu.ID, micropostID, err)
	}

	count, _ := h.store.CountLikes(micropostID)

	if isAJAX(r) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"liked":false,"count":%d}`, count)
		return
	}

	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}
