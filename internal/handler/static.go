package handler

import (
	"net/http"

	"github.com/sofuejin0121/toy_app_go/web/components"
)

// StaticHandler は静的ページを扱うハンドラーです。
type StaticHandler struct{}

// NewStaticHandler は新しいStaticHandlerを作成します。
func NewStaticHandler() *StaticHandler {
	return &StaticHandler{}
}

// Home はHomeページを表示します。
func (h *StaticHandler) Home(w http.ResponseWriter, r *http.Request) {
	data := components.StaticPageData{
		Title: "Home",
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

// Cotact はContactページを表示します。
func (h *StaticHandler) Contact(w http.ResponseWriter, r *http.Request) {
	data := components.StaticPageData{
		Title: "Contact",
		Flash: getFlash(r),
		LoggedIn: isLoggedIn(r),
		CurrentUser: currentUser(r),
		CSRFToken: "",
	}
	_ = components.StaticContact(data).Render(r.Context(), w)
}