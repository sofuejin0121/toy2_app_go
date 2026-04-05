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
	_ = components.HomePage().Render(r.Context(), w)
}

// Help はHelpページを表示します。
func (h *StaticHandler) Help(w http.ResponseWriter, r *http.Request) {
	_ = components.HelpPage().Render(r.Context(), w)
}

// About はAboutページを表示します。
func (h *StaticHandler) About(w http.ResponseWriter, r *http.Request) {
	_ = components.AboutPage().Render(r.Context(), w)
}
