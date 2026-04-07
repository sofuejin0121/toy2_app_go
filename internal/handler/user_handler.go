package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/store"
	"github.com/sofuejin0121/toy_app_go/web/components"
)

// UserHandler はユーザーリソースのHTTPハンドラーです。
type UserHandler struct {
	store *store.Store
}

// NewUserHandler は新しいUserHandlerを返します。
func NewUserHandler(store *store.Store) *UserHandler {
	return &UserHandler{store: store}
}

func noticeFromRequest(r *http.Request) string {
	return r.URL.Query().Get("notice")
}

func redirectWithNotice(w http.ResponseWriter, r *http.Request, path string, notice string) {
	target := path
	if notice != "" {
		target = target + "?notice=" + url.QueryEscape(notice)
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.AllUsers()
	if err != nil {
		log.Printf("AllUsers: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := components.UserPageData{
		Title:  "Users",
		Notice: noticeFromRequest(r),
		Users:  users,
	}
	_ = components.UserIndex(data).Render(r.Context(), w)
}

func (h *UserHandler) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	user, err := h.store.GetUser(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := components.UserPageData{
		Title:       user.Name,
		Flash:       getFlash(r),
		LoggedIn:    isLoggedIn(r),
		CurrentUser: currentUser(r),
		CSRFToken:   "", // CSRFトークンは未実装
		// Notice: noticeFromRequest(r),
		User: *user,
	}
	h.setDebugInfo(&data, r)
	_ = components.UserShow(data).Render(r.Context(), w)
}

func (h *UserHandler) New(w http.ResponseWriter, r *http.Request) {
	data := components.UserPageData{
		Title:       "Sign up",
		Action:      "/users",
		SubmitLabel: "Create my account",
		Flash:       getFlash(r),
		LoggedIn:    isLoggedIn(r),
		CurrentUser: currentUser(r),
		CSRFToken:   "",
	}
	_ = components.UserNew(data).Render(r.Context(), w)
}

func (h *UserHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	user, err := h.store.GetUser(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := components.UserPageData{
		Title:       "Editing user",
		User:        *user,
		Action:      fmt.Sprintf("/users/%d", user.ID),
		SubmitLabel: "Update User",
	}
	_ = components.UserEdit(data).Render(r.Context(), w)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user := model.User{
		Name:                 r.FormValue("name"),
		Email:                r.FormValue("email"),
		Password:             r.FormValue("password"),
		PasswordConfirmation: r.FormValue("password_confirmation"),
	}
	if errors := user.Validate(); len(errors) > 0 {
		data := components.UserPageData{
			Title:       "Sign up",
			Flash:       getFlash(r),
			LoggedIn:    isLoggedIn(r),
			CurrentUser: currentUser(r),
			CSRFToken:   "",
			User:        user,
			Errors:      errors,
			Action:      "/users",
			SubmitLabel: "Create my account",
		}
		h.setDebugInfo(&data, r)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = components.UserNew(data).Render(r.Context(), w)
		return
	}
	// パスワードのハッシュ化と保存
	if err := user.SetPassword(user.Password); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := h.store.CreateUser(&user); err != nil {
		data := components.UserPageData{
			Title:       "Sign up",
			Flash:       getFlash(r),
			LoggedIn:    isLoggedIn(r),
			CurrentUser: currentUser(r),
			CSRFToken:   "",
			User:        user,
			Errors:      []string{err.Error()},
			Action:      "/users",
			SubmitLabel: "Create my account",
		}
		h.setDebugInfo(&data, r)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = components.UserNew(data).Render(r.Context(), w)
		return
	}
	setFlash(w, "success", "Welcome to the Sample App!")
	http.Redirect(w, r, fmt.Sprintf("/users/%d", user.ID), http.StatusSeeOther)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	user, err := h.store.GetUser(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	user.Name = r.FormValue("name")
	user.Email = r.FormValue("email")
	if errors := user.Validate(); len(errors) > 0 {
		data := components.UserPageData{
			Title:       "Editing user",
			User:        *user,
			Errors:      errors,
			Action:      fmt.Sprintf("/users/%d", user.ID),
			SubmitLabel: "Update User",
		}
		_ = components.UserEdit(data).Render(r.Context(), w)
		return
	}
	if err := h.store.UpdateUser(user); err != nil {
		data := components.UserPageData{
			Title:       "Editing user",
			User:        *user,
			Errors:      []string{err.Error()},
			Action:      fmt.Sprintf("/users/%d", user.ID),
			SubmitLabel: "Update User",
		}
		_ = components.UserEdit(data).Render(r.Context(), w)
		return
	}

	redirectWithNotice(w, r, fmt.Sprintf("/users/%d", user.ID), "User was successfully updated.")
}

func (h *UserHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := h.store.DeleteUser(id); err != nil {
		http.NotFound(w, r)
		return
	}
	redirectWithNotice(w, r, "/users", "User was successfully destroyed.")
}
