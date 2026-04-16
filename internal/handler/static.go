package handler

import (
	"net/http"
	"strconv"

	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/store"
	"github.com/sofuejin0121/toy_app_go/web/components"
)

// StaticHandler は静的ページを扱うハンドラーです。
type StaticHandler struct {
	store *store.Store
}

// NewStaticHandler は新しいStaticHandlerを作成します。
func NewStaticHandler(stores ...*store.Store) *StaticHandler {
	h := &StaticHandler{}
	if len(stores) > 0 {
		h.store = stores[0]
	}
	return h
}

// Home はHomeページを表示します。
func (h *StaticHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := components.StaticPageData{
		Title:       "",
		Flash:       getFlash(r),
		LoggedIn:    isLoggedIn(r),
		CurrentUser: currentUser(r),
		CSRFToken:   middleware.CSRFTokenFromContext(r),
	}

	if data.LoggedIn && data.CurrentUser != nil && h.store != nil {
		user := data.CurrentUser

		if cnt, err := h.store.CountMicropostsByUserID(user.ID); err == nil {
			data.MicropostCount = cnt
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}
		perPage := 30
		if items, err := h.store.Feed(user.ID, page, perPage); err == nil {
			data.Microposts = items
		}
		data.Pagination = components.NewPagination(page, perPage, data.MicropostCount)

		data.FollowingCount, _ = h.store.CountFollowing(user.ID)
		data.FollowersCount, _ = h.store.CountFollowers(user.ID)
		data.LikedCount, _ = h.store.CountLikedMicroposts(user.ID)
	}

	_ = components.StaticHome(data).Render(r.Context(), w)
}

// Help はHelpページを表示します。
func (h *StaticHandler) Help(w http.ResponseWriter, r *http.Request) {
	_ = components.HelpPage().Render(r.Context(), w)
}

// About はAboutページを表示します。
func (h *StaticHandler) About(w http.ResponseWriter, r *http.Request) {
	_ = components.AboutPage().Render(r.Context(), w)
}

// Contact はContactページを表示します。
func (h *StaticHandler) Contact(w http.ResponseWriter, r *http.Request) {
	data := components.StaticPageData{
		Title:       "Contact",
		Flash:       getFlash(r),
		LoggedIn:    isLoggedIn(r),
		CurrentUser: currentUser(r),
		CSRFToken:   "",
	}
	_ = components.StaticContact(data).Render(r.Context(), w)
}
