package handler

import (
	"net/http"
	"strconv"
	"time"
)

// GET /api/notifications
func (h *APIHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	items, err := h.store.GetNotifications(cu.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "内部エラーが発生しました")
		return
	}
	h.store.MarkAllRead(cu.ID)

	result := make([]NotificationJSON, len(items))
	for i, item := range items {
		nj := NotificationJSON{
			ID:         item.Notification.ID,
			ActionType: item.Notification.ActionType,
			Read:       item.Notification.Read,
			Actor:      userToJSON(item.Actor),
			CreatedAt:  item.Notification.CreatedAt.Format(time.RFC3339),
		}
		if item.Target != nil {
			nj.TargetID = &item.Target.ID
			c := item.Target.Content
			nj.TargetContent = &c
		}
		result[i] = nj
	}
	writeJSON(w, http.StatusOK, map[string]any{"notifications": result})
}

// DELETE /api/notifications/{id}
func (h *APIHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "IDが不正です")
		return
	}
	if err := h.store.DeleteNotification(id, cu.ID); err != nil {
		writeError(w, http.StatusInternalServerError, "内部エラーが発生しました")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// GET /api/notifications/unread_count
func (h *APIHandler) UnreadNotificationCount(w http.ResponseWriter, r *http.Request) {
	cu := currentUser(r)
	if cu == nil {
		writeJSON(w, http.StatusOK, map[string]int{"count": 0})
		return
	}
	items, err := h.store.GetNotifications(cu.ID)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]int{"count": 0})
		return
	}
	count := 0
	for _, item := range items {
		if !item.Notification.Read {
			count++
		}
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}
