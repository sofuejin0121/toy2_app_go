package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// GET /api/feed?page=
func (h *APIHandler) Feed(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 30
	items, err := h.store.Feed(cu.ID, page, perPage)
	if err != nil {
		log.Printf("Feed: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	total, _ := h.store.CountMicropostsByUserID(cu.ID)
	writeJSON(w, http.StatusOK, map[string]any{
		"items":      feedItemsToJSON(items),
		"pagination": makePagination(page, perPage, total),
	})
}

// GET /api/microposts/{id}
func (h *APIHandler) GetMicropost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	cu := currentUser(r)
	var viewerID int64
	if cu != nil {
		viewerID = cu.ID
	}
	post, err := h.store.GetMicropostAsFeedItem(id, viewerID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	replies, _ := h.store.GetReplies(id, viewerID)
	replyCount, _ := h.store.CountReplies(id)
	writeJSON(w, http.StatusOK, map[string]any{
		"post":        feedItemToJSON(*post),
		"replies":     feedItemsToJSON(replies),
		"reply_count": replyCount,
	})
}

// POST /api/microposts
func (h *APIHandler) CreateMicropost(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	r.ParseMultipartForm(maxImageSize)
	micropost := &model.Micropost{
		Content: r.FormValue("content"),
		UserID:  cu.ID,
	}
	if replyStr := r.FormValue("in_reply_to_id"); replyStr != "" {
		if replyID, err := strconv.ParseInt(replyStr, 10, 64); err == nil && replyID > 0 {
			micropost.InReplyToID = &replyID
		}
	}
	if errs := micropost.Validate(); len(errs) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": errs})
		return
	}
	imagePath, imageErrs := processImageUpload(r, h.imageDir)
	if len(imageErrs) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": imageErrs})
		return
	}
	micropost.ImagePath = imagePath
	if err := h.store.CreateMicropost(micropost); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	post, err := h.store.GetMicropostAsFeedItem(micropost.ID, cu.ID)
	if err != nil {
		writeJSON(w, http.StatusCreated, map[string]any{"id": micropost.ID})
		return
	}
	writeJSON(w, http.StatusCreated, feedItemToJSON(*post))
}

// DELETE /api/microposts/{id}
func (h *APIHandler) DeleteMicropost(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	mp, err := h.store.GetMicropostByUserIDAndID(cu.ID, id)
	if err != nil || mp == nil {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	if mp.ImagePath != "" {
		fullPath := filepath.Join(h.imageDir, mp.ImagePath)
		if removeErr := os.Remove(fullPath); removeErr != nil && !os.IsNotExist(removeErr) {
			log.Printf("Remove image %s: %v", fullPath, removeErr)
		}
	}
	if err := h.store.DeleteMicropost(id); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
