package handler

import (
	"log"
	"net/http"
	"strconv"
)

type BookmarkHandler struct {
	store bookmarkStore
}

type bookmarkStore interface {
	Bookmark(userID, micropostID int64) error
	UnBookmark(userID, micropostID int64) error
}

func NewBookmarkHandler(s bookmarkStore) *BookmarkHandler {
	return &BookmarkHandler{store: s}
}

// Create は POST /bookmarks を処理する。
// フォームの micropost_id に対して現在のユーザーがブックマークを作成する。
func (h *BookmarkHandler) Create(w http.ResponseWriter, r *http.Request) {
	cu := currentUser(r)
	if cu == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	micropostID, err := strconv.ParseInt(r.FormValue("micropost_id"), 10, 64)
	if err != nil || micropostID == 0 {
		http.Error(w, "不正なリクエストです", http.StatusBadRequest)
		return
	}

	if err := h.store.Bookmark(cu.ID, micropostID); err != nil {
		log.Printf("Bookmark(%d, %d): %v", cu.ID, micropostID, err)
	}

	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}

// Destroy は DELETE /bookmarks/{micropost_id} を処理する。
func (h *BookmarkHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	cu := currentUser(r)
	if cu == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	micropostID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || micropostID == 0 {
		http.Error(w, "不正なリクエストです", http.StatusBadRequest)
		return
	}

	if err := h.store.UnBookmark(cu.ID, micropostID); err != nil {
		log.Printf("UnBookmark(%d, %d): %v", cu.ID, micropostID, err)
	}

	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}
