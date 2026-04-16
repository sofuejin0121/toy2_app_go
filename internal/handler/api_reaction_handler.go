package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// POST /api/relationships
func (h *APIHandler) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	var body struct {
		FollowedID int64 `json:"followed_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.store.Follow(cu.ID, body.FollowedID); err != nil {
		log.Printf("Follow %d -> %d: %v", cu.ID, body.FollowedID, err)
	}
	if pref, err := h.store.GetOrCreateUserPreference(body.FollowedID); err == nil && pref.EmailOnFollow {
		if followed, err := h.store.GetUser(body.FollowedID); err == nil {
			if err := h.mailer.SendFollowNotification(followed, cu); err != nil {
				log.Printf("SendFollowNotification: %v", err)
			}
		}
	}
	rel, _ := h.store.GetRelationshipByUsers(cu.ID, body.FollowedID)
	var relID int64
	if rel != nil {
		relID = rel.ID
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"message":         "followed",
		"relationship_id": relID,
	})
}

// DELETE /api/relationships/{id}
func (h *APIHandler) DeleteRelationship(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	relID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	rel, err := h.store.GetRelationship(relID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if err := h.store.Unfollow(cu.ID, rel.FollowedID); err != nil {
		log.Printf("Unfollow: %v", err)
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "unfollowed"})
}

// POST /api/likes
func (h *APIHandler) CreateLike(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	var body struct {
		MicropostID int64 `json:"micropost_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.store.Like(cu.ID, body.MicropostID); err != nil {
		log.Printf("Like: %v", err)
	}
	if mp, err := h.store.GetMicropost(body.MicropostID); err == nil && mp.UserID != cu.ID {
		if pref, err := h.store.GetOrCreateUserPreference(mp.UserID); err == nil && pref.EmailOnLike {
			if owner, err := h.store.GetUser(mp.UserID); err == nil {
				if err := h.mailer.SendLikeNotification(owner, cu, mp.Content); err != nil {
					log.Printf("SendLikeNotification: %v", err)
				}
			}
		}
	}
	count, _ := h.store.CountLikes(body.MicropostID)
	writeJSON(w, http.StatusOK, map[string]any{"liked": true, "count": count})
}

// DELETE /api/likes/{id}
func (h *APIHandler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	micropostID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.store.Unlike(cu.ID, micropostID); err != nil {
		log.Printf("Unlike: %v", err)
	}
	count, _ := h.store.CountLikes(micropostID)
	writeJSON(w, http.StatusOK, map[string]any{"liked": false, "count": count})
}

// POST /api/bookmarks
func (h *APIHandler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	var body struct {
		MicropostID int64 `json:"micropost_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.store.Bookmark(cu.ID, body.MicropostID); err != nil {
		log.Printf("Bookmark: %v", err)
	}
	writeJSON(w, http.StatusOK, map[string]any{"bookmarked": true})
}

// DELETE /api/bookmarks/{id}
func (h *APIHandler) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	micropostID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.store.UnBookmark(cu.ID, micropostID); err != nil {
		log.Printf("UnBookmark: %v", err)
	}
	writeJSON(w, http.StatusOK, map[string]any{"bookmarked": false})
}
