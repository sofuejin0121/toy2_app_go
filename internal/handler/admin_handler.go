package handler

import (
	"net/http"
	"github.com/sofuejin0121/toy_app_go/internal/store"
	"github.com/sofuejin0121/toy_app_go/web/components"

	"fmt"
	"os"
	"github.com/sofuejin0121/toy_app_go/internal/middleware"

)

// ハンドラー構造体
type AdminHandler struct {
	store adminStore // storeはインターフェース型
}

// インターフェイスDI
type adminStore interface {
	GetAdminStats() (store.AdminStats, error)
}

func NewAdminHandler(s adminStore) *AdminHandler {
	return &AdminHandler{store: s}
}

// ダッシュボード表示
func (h *AdminHandler) Index(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetAdminStats()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := components.AdminPageData{
		Title: "Admin Dashboard",
		Flash: getFlash(r),
		LoggedIn: isLoggedIn(r),
		CSRFToken: middleware.CSRFTokenFromContext(r),
		CurrentUser: currentUser(r),
		Debug: os.Getenv("APP_ENV") != "production",
		DebugInfo: fmt.Sprintf("%+v", stats),
		Stats: stats,
	}
	_ = components.AdminIndex(data).Render(r.Context(), w)
}