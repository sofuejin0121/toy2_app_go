package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sofuejin0121/toy_app_go/internal/model"
)

// POST /api/users (サインアップ)
func (h *APIHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name                 string `json:"name"`
		Email                string `json:"email"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	user := model.User{
		Name:                 body.Name,
		Email:                body.Email,
		Password:             body.Password,
		PasswordConfirmation: body.PasswordConfirmation,
	}
	if errs := user.Validate(); len(errs) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": errs})
		return
	}
	if err := user.SetPassword(user.Password); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if err := h.store.CreateUser(&user); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": []string{err.Error()}})
		return
	}
	if err := user.SendActivationEmail(h.mailer); err != nil {
		log.Printf("SendActivationEmail: %v", err)
	}
	writeJSON(w, http.StatusCreated, userToJSON(user))
}

// GET /api/users?page=&q=
func (h *APIHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 30

	var users []model.User
	var total int
	var err error

	if query != "" {
		users, err = h.store.SearchActivatedUsers(query, page, perPage)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		total, _ = h.store.CountSearchActivatedUsers(query)
	} else {
		users, err = h.store.GetActivatedUsers(page)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		total, _ = h.store.CountActivatedUsers()
	}

	ujsons := make([]UserJSON, len(users))
	for i, u := range users {
		ujsons[i] = userToJSON(u)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"users":      ujsons,
		"pagination": makePagination(page, perPage, total),
	})
}

// GET /api/users/{id}
func (h *APIHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	user, err := h.store.GetUser(id)
	if err != nil || !user.Activated {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 30

	cu := currentUser(r)
	var viewerID int64
	if cu != nil {
		viewerID = cu.ID
	}

	microposts, _ := h.store.PaginateMicropostsWithStats(user.ID, viewerID, page, perPage)
	micropostCount, _ := h.store.CountMicropostsByUserID(user.ID)
	followingCount, _ := h.store.CountFollowing(user.ID)
	followersCount, _ := h.store.CountFollowers(user.ID)
	likedCount, _ := h.store.CountLikedMicroposts(user.ID)
	bookmarkCount, _ := h.store.CountBookmarkedMicroposts(user.ID)

	resp := UserProfileJSON{
		User:           userToJSON(*user),
		MicropostCount: micropostCount,
		FollowingCount: followingCount,
		FollowersCount: followersCount,
		LikedCount:     likedCount,
		BookmarkCount:  bookmarkCount,
		Microposts:     feedItemsToJSON(microposts),
		Pagination:     makePagination(page, perPage, micropostCount),
	}
	if cu != nil {
		resp.IsCurrentUser = cu.ID == user.ID
		if !resp.IsCurrentUser {
			isFollowing, _ := h.store.IsFollowing(cu.ID, user.ID)
			resp.IsFollowing = isFollowing
			if isFollowing {
				rel, err := h.store.GetRelationshipByUsers(cu.ID, user.ID)
				if err == nil {
					resp.RelationshipID = rel.ID
				}
			}
		}
	}
	writeJSON(w, http.StatusOK, resp)
}

// PATCH /api/users/{id}
func (h *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || cu.ID != id {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	user, err := h.store.GetUser(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	var body struct {
		Name                 string `json:"name"`
		Email                string `json:"email"`
		Bio                  string `json:"bio"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	user.Name = strings.TrimSpace(body.Name)
	user.Email = strings.ToLower(strings.TrimSpace(body.Email))
	user.Bio = strings.TrimSpace(body.Bio)

	errs := user.Validate()
	if body.Password != "" {
		if body.Password != body.PasswordConfirmation {
			errs = append(errs, "Password confirmation doesn't match Password")
		}
		if err := model.ValidatePassword(body.Password); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": errs})
		return
	}
	if body.Password != "" {
		if err := h.store.UpdatePassword(id, body.Password); err != nil {
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}
	if err := h.store.UpdateUser(user); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	updated, _ := h.store.GetUser(id)
	writeJSON(w, http.StatusOK, userToJSON(*updated))
}

// DELETE /api/users/{id} (管理者のみ)
func (h *APIHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil || !cu.Admin {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if cu.ID == id {
		writeError(w, http.StatusBadRequest, "cannot delete own account")
		return
	}
	if err := h.store.DeleteUser(id); err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// GET /api/users/{id}/following
func (h *APIHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	user, err := h.store.GetUser(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 30
	users, _ := h.store.PaginateFollowing(user.ID, page, perPage)
	followingCount, _ := h.store.CountFollowing(user.ID)
	followersCount, _ := h.store.CountFollowers(user.ID)
	micropostCount, _ := h.store.CountMicropostsByUserID(user.ID)
	ujsons := make([]UserJSON, len(users))
	for i, u := range users {
		ujsons[i] = userToJSON(u)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user":            userToJSON(*user),
		"users":           ujsons,
		"following_count": followingCount,
		"followers_count": followersCount,
		"micropost_count": micropostCount,
		"pagination":      makePagination(page, perPage, followingCount),
	})
}

// GET /api/users/{id}/followers
func (h *APIHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	user, err := h.store.GetUser(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 30
	users, _ := h.store.PaginateFollowers(user.ID, page, perPage)
	followingCount, _ := h.store.CountFollowing(user.ID)
	followersCount, _ := h.store.CountFollowers(user.ID)
	micropostCount, _ := h.store.CountMicropostsByUserID(user.ID)
	ujsons := make([]UserJSON, len(users))
	for i, u := range users {
		ujsons[i] = userToJSON(u)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user":            userToJSON(*user),
		"users":           ujsons,
		"following_count": followingCount,
		"followers_count": followersCount,
		"micropost_count": micropostCount,
		"pagination":      makePagination(page, perPage, followersCount),
	})
}

// GET /api/users/{id}/likes
func (h *APIHandler) GetUserLikes(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	user, err := h.store.GetUser(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 30
	cu := currentUser(r)
	var viewerID int64
	if cu != nil {
		viewerID = cu.ID
	}
	items, _ := h.store.PaginateLikedMicropostsAsFeedItems(user.ID, viewerID, page, perPage)
	likedCount, _ := h.store.CountLikedMicroposts(user.ID)
	followingCount, _ := h.store.CountFollowing(user.ID)
	followersCount, _ := h.store.CountFollowers(user.ID)
	micropostCount, _ := h.store.CountMicropostsByUserID(user.ID)
	writeJSON(w, http.StatusOK, map[string]any{
		"user":            userToJSON(*user),
		"microposts":      feedItemsToJSON(items),
		"liked_count":     likedCount,
		"following_count": followingCount,
		"followers_count": followersCount,
		"micropost_count": micropostCount,
		"pagination":      makePagination(page, perPage, likedCount),
	})
}

// GET /api/users/{id}/bookmarks
func (h *APIHandler) GetUserBookmarks(w http.ResponseWriter, r *http.Request) {
	cu := h.requireAuth(w, r)
	if cu == nil {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if cu.ID != id {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	user, err := h.store.GetUser(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	const perPage = 30
	items, _ := h.store.GetBookmarkedPosts(user.ID, cu.ID, page, perPage)
	bookmarkCount, _ := h.store.CountBookmarkedMicroposts(user.ID)
	followingCount, _ := h.store.CountFollowing(user.ID)
	followersCount, _ := h.store.CountFollowers(user.ID)
	micropostCount, _ := h.store.CountMicropostsByUserID(user.ID)
	likedCount, _ := h.store.CountLikedMicroposts(user.ID)
	writeJSON(w, http.StatusOK, map[string]any{
		"user":            userToJSON(*user),
		"microposts":      feedItemsToJSON(items),
		"bookmark_count":  bookmarkCount,
		"following_count": followingCount,
		"followers_count": followersCount,
		"micropost_count": micropostCount,
		"liked_count":     likedCount,
		"pagination":      makePagination(page, perPage, bookmarkCount),
	})
}
