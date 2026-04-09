package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/sofuejin0121/toy_app_go/internal/middleware"
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

// Index はすべてのユーザーを一覧表示します（ページネーション付き）。
func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	perPage := 30
	users, err := h.store.PaginateUsers(page, perPage)
	if err != nil {
		log.Printf("PaginateUsers: %v", err)
		http.Error(w, "Internal Server Error",
			http.StatusInternalServerError)
		return
	}
	totalUsers, err := h.store.CountUsers()
	if err != nil {
		log.Printf("CountUsers: %v", err)
		http.Error(w, "Internal Server Error",
			http.StatusInternalServerError)
		return
	}
	pagination := components.NewPagination(page, perPage, totalUsers)
	data := components.UserPageData{
		Title:       "All users",
		Flash:       getFlash(r),
		LoggedIn:    isLoggedIn(r),
		CurrentUser: currentUser(r),
		CSRFToken:   middleware.CSRFTokenFromContext(r),
		Users:       users,
		Pagination:  pagination,
	}
	components.UserIndex(data).Render(r.Context(), w)
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
		CSRFToken:   middleware.CSRFTokenFromContext(r),
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
		CSRFToken:   middleware.CSRFTokenFromContext(r),
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
		Title:       "Edit user",
		Flash:       getFlash(r),
		LoggedIn:    isLoggedIn(r),
		CurrentUser: currentUser(r),
		CSRFToken:   middleware.CSRFTokenFromContext(r),
		User:        *user,
	}
	components.UserEdit(data).Render(r.Context(), w)
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
			CSRFToken:   middleware.CSRFTokenFromContext(r),
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
			CSRFToken:   middleware.CSRFTokenFromContext(r),
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
	logIn(w, r, user.ID, h.store)
	setFlash(w, "success", "Welcome to the Sample App!")
	http.Redirect(w, r, fmt.Sprintf("/users/%d", user.ID), http.StatusSeeOther)
}

// Update はユーザー情報を更新する
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password_confirmation")

	user.Name = name
	user.Email = email
	errs := user.Validate()
	if password != "" {
		if password != passwordConfirmation {
			errs = append(errs, "Password confirmation doesn't match Password")
		}
		if err := model.ValidatePassword(password); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		data := components.UserPageData{
			Title:       "Edit user",
			Flash:       getFlash(r),
			LoggedIn:    isLoggedIn(r),
			CurrentUser: currentUser(r),
			CSRFToken:   middleware.CSRFTokenFromContext(r),
			User:        *user,
			Errors:      errs,
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		components.UserEdit(data).Render(r.Context(), w)
		return
	}

	if password != "" {
		if err := h.store.UpdatePassword(id, password); err != nil {
			log.Printf("UpdatePassword(%d): %v", id, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if err := h.store.UpdateUser(user); err != nil {
		log.Printf("UpdateUser: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	setFlash(w, "success", "Profile updated")
	http.Redirect(w, r, fmt.Sprintf("/users/%d", user.ID), http.StatusSeeOther)
}

// Destroy はユーザーを削除します。
func (h *UserHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// 管理者が自分自身を削除することを防止
	if isCurrentUser(r, id) {
		setFlash(w, "danger", "Cannot delete own account")
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}
	if err := h.store.DeleteUser(id); err != nil {
		log.Printf("DeleteUser: %v", err)
		http.Error(w, "Internal Server Error",
			http.StatusInternalServerError)
		return
	}
	setFlash(w, "success", "User deleted")
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// requireLogin はログイン済みユーザーかどうかで確認するミドルウェア
func (h *UserHandler) RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isLoggedIn(r) {
			storeLocation(w, r)
			setFlash(w, "danger", "Please log in.")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// requireCorrectUser は正しいユーザーかどうか確認するミドルウェア
func (h *UserHandler) RequireCorrectUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if !isCurrentUser(r, id) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// requireAdmin は管理者ユーザーかどうか確認するミドルウェア
func (h *UserHandler) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := currentUser(r)
		if user == nil || !user.Admin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}