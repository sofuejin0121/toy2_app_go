package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sofuejin0121/toy_app_go/internal/mailer"
	"github.com/sofuejin0121/toy_app_go/internal/store"
)

// RelationshipHandler はリレーションシップ関連のHTTPハンドラー
type RelationshipHandler struct {
	store  *store.Store
	mailer mailer.Mailer
}

// NewRelationshipHandler は新しいRelationshipHandlerを作成する
func NewRelationshipHandler(s *store.Store, m mailer.Mailer) *RelationshipHandler {
	return &RelationshipHandler{store: s, mailer: m}
}

// Create はフォローを作成する（RequireLoginミドルウェア適用済み）
func (h *RelationshipHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	followedID, err := strconv.ParseInt(r.FormValue("followed_id"), 10, 64)
	if err != nil {
		http.Error(w, "フォロー対象IDが不正です", http.StatusBadRequest)
		return
	}
	if err := h.store.Follow(user.ID, followedID); err != nil {
		log.Printf("Follow %d -> %d: %v", user.ID, followedID, err)
	}

	// フォローされたユーザーの通知設定を確認してメール送信
	if pref, err := h.store.GetOrCreateUserPreference(followedID); err == nil && pref.EmailOnFollow {
		if followed, err := h.store.GetUser(followedID); err == nil {
			if err := h.mailer.SendFollowNotification(followed, user); err != nil {
				log.Printf("SendFollowNotification: %v", err)
			}
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d", followedID), http.StatusSeeOther)
}

// Destroy はフォローを解除する（RequireLoginミドルウェア適用済み）
func (h *RelationshipHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	relationshipID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "リレーションシップIDが不正です", http.StatusBadRequest)
		return
	}
	relationship, err := h.store.GetRelationship(relationshipID)
	if err != nil {
		http.Error(w, "リレーションシップが見つかりません", http.StatusNotFound)
		return
	}
	if err := h.store.Unfollow(user.ID, relationship.FollowedID); err != nil {
		log.Printf("Unfollow %d -> %d: %v", user.ID, relationship.FollowedID, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/users/%d", relationship.FollowedID), http.StatusSeeOther)
}
