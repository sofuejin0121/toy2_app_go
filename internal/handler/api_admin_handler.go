package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// GET /api/admin
func (h *APIHandler) AdminStats(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil || !cu.Admin {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	stats, err := h.store.GetAdminStats()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	dailys := make([]map[string]any, len(stats.DailySignups))
	for i, d := range stats.DailySignups {
		dailys[i] = map[string]any{"date": d.Date, "count": d.Count}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"total_users":   stats.TotalUsers,
		"total_posts":   stats.TotalPosts,
		"today_signups": stats.TodaySignups,
		"daily_signups": dailys,
	})
}

// GET /api/settings
func (h *APIHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	pref, err := h.store.GetOrCreateUserPreference(cu.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, SettingsJSON{
		EmailOnFollow: pref.EmailOnFollow,
		EmailOnLike:   pref.EmailOnLike,
	})
}

// PATCH /api/settings
func (h *APIHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	var body SettingsJSON
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.store.UpdateUserPreference(cu.ID, body.EmailOnFollow, body.EmailOnLike); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, body)
}

// ServeReact はReact SPAのindex.htmlを返す（全フロントエンドルート用）
func ServeReact(distDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexPath := filepath.Join(distDir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			http.Error(w, fmt.Sprintf("React build not found. Run: cd frontend && npm run build"), http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, indexPath)
	}
}
