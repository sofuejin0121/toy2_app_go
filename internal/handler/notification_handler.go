package handler

import (
	"net/http"
	"strconv"

	"github.com/sofuejin0121/toy_app_go/internal/store"
	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/web/components"
	"fmt"
	"os"

) 

type NotificationHandler struct {
	store notificationStore
}

type notificationStore interface {
	GetNotifications(userID int64) ([]store.NotificationItem, error)
	MarkAllRead(userID int64) error
	DeleteNotification(id, userID int64) error
}

func NewNotificationHandler(s notificationStore) *NotificationHandler {
	return &NotificationHandler{store: s}
}

func (h *NotificationHandler) Index(w http.ResponseWriter, r *http.Request) {
	cu := currentUser(r)
	if cu == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	items, err := h.store.GetNotifications(cu.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.store.MarkAllRead(cu.ID)
	data := components.NotificationDataPage{
		Title: "Notifications",
		Flash: getFlash(r),
		LoggedIn: isLoggedIn(r),
		CurrentUser: currentUser(r),
		CSRFToken: middleware.CSRFTokenFromContext(r),
		Debug: os.Getenv("APP_ENV") != "production",
		DebugInfo: fmt.Sprintf("%+v", items),
		Items: items,
	}
	_ = components.NotificationIndex(data).Render(r.Context(), w)
}

func (h *NotificationHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	cu := currentUser(r)
	if cu == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := h.store.DeleteNotification(id, cu.ID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	setFlash(w, "success", "Notification deleted")
	http.Redirect(w, r, "/notifications", http.StatusSeeOther)
}