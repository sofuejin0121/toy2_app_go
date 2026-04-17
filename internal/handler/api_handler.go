package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sofuejin0121/toy_app_go/internal/mailer"
	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/storage"
	"github.com/sofuejin0121/toy_app_go/internal/store"
)

// APIHandler はJSON APIエンドポイントをすべて担当するハンドラーです。
type APIHandler struct {
	store   *store.Store
	mailer  mailer.Mailer
	storage storage.ImageStorage
}

// NewAPIHandler は新しいAPIHandlerを作成します。
func NewAPIHandler(s *store.Store, m mailer.Mailer, st storage.ImageStorage) *APIHandler {
	return &APIHandler{store: s, mailer: m, storage: st}
}

// ---- JSON レスポンス型 ----

type UserJSON struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	Admin     bool   `json:"admin"`
	Activated bool   `json:"activated"`
	AvatarURL string `json:"avatar_url"`
	CreatedAt string `json:"created_at"`
}

type MicropostJSON struct {
	ID           int64       `json:"id"`
	Content      string      `json:"content"`
	UserID       int64       `json:"user_id"`
	ImagePath    string      `json:"image_path,omitempty"`
	InReplyToID  *int64      `json:"in_reply_to_id,omitempty"`
	LikeCount    int         `json:"like_count"`
	IsLiked      bool        `json:"is_liked"`
	IsBookmarked bool        `json:"is_bookmarked"`
	User         UserJSON    `json:"user"`
	Parent       *ParentJSON `json:"parent,omitempty"`
	CreatedAt    string      `json:"created_at"`
}

type ParentJSON struct {
	ID      int64    `json:"id"`
	Content string   `json:"content"`
	User    UserJSON `json:"user"`
}

type PaginationJSON struct {
	CurrentPage int  `json:"current_page"`
	TotalPages  int  `json:"total_pages"`
	TotalItems  int  `json:"total_items"`
	PerPage     int  `json:"per_page"`
	HasPrev     bool `json:"has_prev"`
	HasNext     bool `json:"has_next"`
}

type UserProfileJSON struct {
	User           UserJSON        `json:"user"`
	IsCurrentUser  bool            `json:"is_current_user"`
	IsFollowing    bool            `json:"is_following"`
	RelationshipID int64           `json:"relationship_id,omitempty"`
	MicropostCount int             `json:"micropost_count"`
	FollowingCount int             `json:"following_count"`
	FollowersCount int             `json:"followers_count"`
	LikedCount     int             `json:"liked_count"`
	BookmarkCount  int             `json:"bookmark_count"`
	Microposts     []MicropostJSON `json:"microposts"`
	Pagination     PaginationJSON  `json:"pagination"`
}

type NotificationJSON struct {
	ID            int64    `json:"id"`
	ActionType    string   `json:"action_type"`
	Read          bool     `json:"read"`
	Actor         UserJSON `json:"actor"`
	TargetID      *int64   `json:"target_id,omitempty"`
	TargetContent *string  `json:"target_content,omitempty"`
	CreatedAt     string   `json:"created_at"`
}

type SettingsJSON struct {
	EmailOnFollow bool `json:"email_on_follow"`
	EmailOnLike   bool `json:"email_on_like"`
}

// ---- 共通ヘルパー ----

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *APIHandler) requireAuth(w http.ResponseWriter, r *http.Request) *model.User {
	cu := currentUser(r)
	if cu == nil {
		writeError(w, http.StatusUnauthorized, "認証が必要です")
		return nil
	}
	return cu
}

func userToJSON(u model.User) UserJSON {
	return UserJSON{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Bio:       u.Bio,
		Admin:     u.Admin,
		Activated: u.Activated,
		AvatarURL: u.GravatarURL(50),
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

func calcTotalPages(total, perPage int) int {
	if perPage == 0 {
		return 0
	}
	t := total / perPage
	if total%perPage != 0 {
		t++
	}
	return t
}

func makePagination(page, perPage, total int) PaginationJSON {
	totalPages := calcTotalPages(total, perPage)
	return PaginationJSON{
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalItems:  total,
		PerPage:     perPage,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}
}

func (h *APIHandler) feedItemToJSON(item store.FeedItem) MicropostJSON {
	imageURL := ""
	if item.Micropost.ImagePath != "" {
		imageURL = h.storage.PublicURL(item.Micropost.ImagePath)
	}
	mj := MicropostJSON{
		ID:           item.Micropost.ID,
		Content:      item.Micropost.Content,
		UserID:       item.Micropost.UserID,
		ImagePath:    imageURL,
		InReplyToID:  item.Micropost.InReplyToID,
		LikeCount:    item.LikeCount,
		IsLiked:      item.IsLiked,
		IsBookmarked: item.IsBookmarked,
		User:         userToJSON(item.User),
		CreatedAt:    item.Micropost.CreatedAt.Format(time.RFC3339),
	}
	if item.ParentMicropost != nil && item.ParentUser != nil {
		mj.Parent = &ParentJSON{
			ID:      item.ParentMicropost.ID,
			Content: item.ParentMicropost.Content,
			User:    userToJSON(*item.ParentUser),
		}
	}
	return mj
}

func (h *APIHandler) feedItemsToJSON(items []store.FeedItem) []MicropostJSON {
	result := make([]MicropostJSON, len(items))
	for i, item := range items {
		result[i] = h.feedItemToJSON(item)
	}
	return result
}
