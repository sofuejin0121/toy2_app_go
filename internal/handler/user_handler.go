package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/sofuejin0121/toy_app_go/internal/mailer"
	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/store"
	"github.com/sofuejin0121/toy_app_go/web/components"
)

// UserHandler はユーザーリソースのHTTPハンドラーです。
type UserHandler struct {
	store  *store.Store
	mailer mailer.Mailer
}

// NewUserHandler は新しいUserHandlerを返します。
func NewUserHandler(store *store.Store, m mailer.Mailer) *UserHandler {
	return &UserHandler{store: store, mailer: m}
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

// Index はすべてのユーザーを一覧表示します（検索・ページネーション付き）。
func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	// ?q= が指定されていれば検索モード、なければ全件取得
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 30

	var users []model.User
	var totalUsers int
	var err error

	if query != "" {
		users, err = h.store.SearchActivatedUsers(query, page, perPage)
		if err != nil {
			log.Printf("SearchActivatedUsers: %v", err)
			http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
			return
		}
		totalUsers, err = h.store.CountSearchActivatedUsers(query)
		if err != nil {
			log.Printf("CountSearchActivatedUsers: %v", err)
			http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
			return
		}
	} else {
		users, err = h.store.GetActivatedUsers(page)
		if err != nil {
			log.Printf("GetActivatedUsers: %v", err)
			http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
			return
		}
		totalUsers, err = h.store.CountActivatedUsers()
		if err != nil {
			log.Printf("CountActivatedUsers: %v", err)
			http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
			return
		}
	}

	// ページネーションのベースURL: 検索クエリがあれば保持する
	paginationBase := "/users"
	if query != "" {
		paginationBase = fmt.Sprintf("/users?q=%s", url.QueryEscape(query))
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
		SearchQuery: query,
	}
	_ = components.UserIndex(data, paginationBase).Render(r.Context(), w)
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
	if !user.Activated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 30

	// 閲覧者のIDを取得（未ログイン時は 0 で IsLiked が常に false）
	var viewerID int64
	if cu := currentUser(r); cu != nil {
		viewerID = cu.ID
	}
	microposts, err := h.store.PaginateMicropostsWithStats(user.ID, viewerID, page, perPage)
	if err != nil {
		log.Printf("PaginateMicropostsWithStats(%d): %v", user.ID, err)
	}
	micropostCount, err := h.store.CountMicropostsByUserID(user.ID)
	if err != nil {
		log.Printf("CountMicropostsByUserID(%d): %v", user.ID, err)
	}
	followingCount, _ := h.store.CountFollowing(user.ID)
	followersCount, _ := h.store.CountFollowers(user.ID)
	likedCount, _ := h.store.CountLikedMicroposts(user.ID)
	bookmarkCount, _ := h.store.CountBookmarkedMicroposts(user.ID)

	data := components.UserPageData{
		Title:          user.Name,
		Flash:          getFlash(r),
		LoggedIn:       isLoggedIn(r),
		CurrentUser:    currentUser(r),
		CSRFToken:      middleware.CSRFTokenFromContext(r),
		User:           *user,
		Microposts:     microposts,
		MicropostCount: micropostCount,
		FollowingCount: followingCount,
		FollowersCount: followersCount,
		LikedCount:     likedCount,
		BookmarkCount:  bookmarkCount,
		Pagination:     components.NewPagination(page, perPage, micropostCount),
	}

	cu := currentUser(r)
	if cu != nil {
		data.IsCurrentUser = cu.ID == user.ID
		if !data.IsCurrentUser {
			isFollowing, _ := h.store.IsFollowing(cu.ID, user.ID)
			data.IsFollowing = isFollowing
			if isFollowing {
				rel, err := h.store.GetRelationshipByUsers(cu.ID, user.ID)
				if err == nil {
					data.RelationshipID = rel.ID
				}
			}
		}
	}

	h.setDebugInfo(&data, r)
	_ = components.UserShow(data).Render(r.Context(), w)
}

// Following はフォローしているユーザー一覧を表示する
func (h *UserHandler) Following(w http.ResponseWriter, r *http.Request) {
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
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 30
	users, err := h.store.PaginateFollowing(user.ID, page, perPage)
	if err != nil {
		log.Printf("PaginateFollowing(%d): %v", user.ID, err)
	}
	followingCount, _ := h.store.CountFollowing(user.ID)
	followersCount, _ := h.store.CountFollowers(user.ID)
	likedCountF, _ := h.store.CountLikedMicroposts(user.ID)
	micropostCount, _ := h.store.CountMicropostsByUserID(user.ID)

	data := components.UserPageData{
		Title:          "Following",
		Flash:          getFlash(r),
		LoggedIn:       isLoggedIn(r),
		CurrentUser:    currentUser(r),
		CSRFToken:      middleware.CSRFTokenFromContext(r),
		User:           *user,
		Users:          users,
		FollowingCount: followingCount,
		FollowersCount: followersCount,
		LikedCount:     likedCountF,
		MicropostCount: micropostCount,
		Pagination:     components.NewPagination(page, perPage, followingCount),
	}
	h.setDebugInfo(&data, r)
	_ = components.ShowFollow(data).Render(r.Context(), w)
}

// Followers はフォロワー一覧を表示する
func (h *UserHandler) Followers(w http.ResponseWriter, r *http.Request) {
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
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 30
	users, err := h.store.PaginateFollowers(user.ID, page, perPage)
	if err != nil {
		log.Printf("PaginateFollowers(%d): %v", user.ID, err)
	}
	followingCount, _ := h.store.CountFollowing(user.ID)
	followersCount, _ := h.store.CountFollowers(user.ID)
	likedCountR, _ := h.store.CountLikedMicroposts(user.ID)
	micropostCount, _ := h.store.CountMicropostsByUserID(user.ID)

	data := components.UserPageData{
		Title:          "Followers",
		Flash:          getFlash(r),
		LoggedIn:       isLoggedIn(r),
		CurrentUser:    currentUser(r),
		CSRFToken:      middleware.CSRFTokenFromContext(r),
		User:           *user,
		Users:          users,
		FollowingCount: followingCount,
		FollowersCount: followersCount,
		LikedCount:     likedCountR,
		MicropostCount: micropostCount,
		Pagination:     components.NewPagination(page, perPage, followersCount),
	}
	h.setDebugInfo(&data, r)
	_ = components.ShowFollow(data).Render(r.Context(), w)
}

// LikedPosts はユーザーがいいねした投稿一覧を表示する。
func (h *UserHandler) LikedPosts(w http.ResponseWriter, r *http.Request) {
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

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 30

	var viewerID int64
	if cu := currentUser(r); cu != nil {
		viewerID = cu.ID
	}

	items, err := h.store.PaginateLikedMicropostsAsFeedItems(user.ID, viewerID, page, perPage)
	if err != nil {
		log.Printf("PaginateLikedMicropostsAsFeedItems(%d): %v", user.ID, err)
	}
	likedCount, _ := h.store.CountLikedMicroposts(user.ID)
	followingCount, _ := h.store.CountFollowing(user.ID)
	followersCount, _ := h.store.CountFollowers(user.ID)
	micropostCount, _ := h.store.CountMicropostsByUserID(user.ID)

	data := components.UserPageData{
		Title:          fmt.Sprintf("%s's Likes", user.Name),
		Flash:          getFlash(r),
		LoggedIn:       isLoggedIn(r),
		CurrentUser:    currentUser(r),
		CSRFToken:      middleware.CSRFTokenFromContext(r),
		User:           *user,
		Microposts:     items,
		MicropostCount: micropostCount,
		FollowingCount: followingCount,
		FollowersCount: followersCount,
		LikedCount:     likedCount,
		Pagination:     components.NewPagination(page, perPage, likedCount),
	}
	_ = components.ShowLikes(data).Render(r.Context(), w)
}

// BookmarkedPosts はユーザーがブックマークした投稿一覧を表示する。
// 自分のブックマークのみ閲覧可能。
func (h *UserHandler) BookmarkedPosts(w http.ResponseWriter, r *http.Request) {
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

	// 自分以外のブックマークは閲覧不可
	cu := currentUser(r)
	if cu == nil || cu.ID != user.ID {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 30

	var viewerID int64
	if cu := currentUser(r); cu != nil {
		viewerID = cu.ID
	}

	items, err := h.store.GetBookmarkedPosts(user.ID, viewerID, page, perPage)
	if err != nil {
		log.Printf("GetBookmarkedPosts(%d): %v", user.ID, err)
	}
	bookmarkCount, _ := h.store.CountBookmarkedMicroposts(user.ID)
	followingCount, _ := h.store.CountFollowing(user.ID)
	followersCount, _ := h.store.CountFollowers(user.ID)
	micropostCount, _ := h.store.CountMicropostsByUserID(user.ID)
	likedCount, _ := h.store.CountLikedMicroposts(user.ID)

	data := components.UserPageData{
		Title:          fmt.Sprintf("%s's Bookmarks", user.Name),
		Flash:          getFlash(r),
		LoggedIn:       isLoggedIn(r),
		CurrentUser:    currentUser(r),
		CSRFToken:      middleware.CSRFTokenFromContext(r),
		User:           *user,
		Microposts:     items,
		MicropostCount: micropostCount,
		FollowingCount: followingCount,
		FollowersCount: followersCount,
		LikedCount:     likedCount,
		BookmarkCount:  bookmarkCount,
		Pagination:     components.NewPagination(page, perPage, bookmarkCount),
	}
	_ = components.ShowBookmarks(data).Render(r.Context(), w)
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
		http.Error(w, "不正なリクエストです", http.StatusBadRequest)
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
		http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
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
	// リファクタリング後: Userモデルのメソッドを使用
	if err := user.SendActivationEmail(h.mailer); err != nil {
		log.Printf("SendActivationEmail: %v", err)
	}
	setFlash(w, "success", "Chirpへようこそ！")
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
	bio := strings.TrimSpace(r.FormValue("bio"))
	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password_confirmation")

	user.Name = name
	user.Email = email
	user.Bio = bio
	errs := user.Validate()
	if password != "" {
		if password != passwordConfirmation {
			errs = append(errs, "パスワード確認が一致しません")
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
			http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
			return
		}
	}

	if err := h.store.UpdateUser(user); err != nil {
		log.Printf("UpdateUser: %v", err)
		http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
		return
	}
	setFlash(w, "success", "プロフィールを更新しました")
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
		setFlash(w, "danger", "自分のアカウントは削除できません")
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}
	if err := h.store.DeleteUser(id); err != nil {
		log.Printf("DeleteUser: %v", err)
		http.Error(w, "内部サーバーエラー",
			http.StatusInternalServerError)
		return
	}
	setFlash(w, "success", "ユーザーを削除しました")
	http.Redirect(w, r, "/users", http.StatusSeeOther)
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
