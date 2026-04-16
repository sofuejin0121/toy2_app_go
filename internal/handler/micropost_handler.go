package handler

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sofuejin0121/toy_app_go/internal/middleware"
	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/store"
	"github.com/sofuejin0121/toy_app_go/web/components"
)

// maxImageSize は画像アップロードの最大サイズ（5MB）です。
const maxImageSize = 5 * 1024 * 1024

// allowedImageTypes は許可する画像MIMEタイプです。
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

// MicropostHandler はマイクロポストリソースのHTTPハンドラーです。
type MicropostHandler struct {
	store    *store.Store
	imageDir string
}

// NewMicropostHandler は新しいMicropostHandlerを作成します。
func NewMicropostHandler(s *store.Store, imageDir string) *MicropostHandler {
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		log.Printf("MkdirAll %s: %v", imageDir, err)
	}
	return &MicropostHandler{
		store:    s,
		imageDir: imageDir,
	}
}

// RequireCorrectUser は現在のユーザーが対象マイクロポストの所有者かを確認します。
func (h *MicropostHandler) RequireCorrectUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		user := currentUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		mp, err := h.store.GetMicropostByUserIDAndID(user.ID, id)
		if err != nil || mp == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// Create はマイクロポストを作成します。
func (h *MicropostHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	micropost := &model.Micropost{
		Content: r.FormValue("content"),
		UserID:  user.ID,
	}

	// リプライ元IDが指定されていれば設定する
	if replyStr := r.FormValue("in_reply_to_id"); replyStr != "" {
		if replyID, err := strconv.ParseInt(replyStr, 10, 64); err == nil && replyID > 0 {
			micropost.InReplyToID = &replyID
		}
	}

	if errs := micropost.Validate(); len(errs) > 0 {
		setFlash(w, "danger", strings.Join(errs, ", "))
		// リプライの場合はリプライ元ページへ、通常投稿はホームへ戻す
		if micropost.InReplyToID != nil {
			http.Redirect(w, r, fmt.Sprintf("/microposts/%d", *micropost.InReplyToID), http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}

	imagePath, imageErrs := h.processImageUpload(r)
	if len(imageErrs) > 0 {
		setFlash(w, "danger", strings.Join(imageErrs, ", "))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	micropost.ImagePath = imagePath

	if err := h.store.CreateMicropost(micropost); err != nil {
		log.Printf("CreateMicropost: %v", err)
		setFlash(w, "danger", "Failed to create micropost")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	setFlash(w, "success", "Micropost created!")
	// リプライの場合はリプライ元の詳細ページへ遷移 → 返信がぶら下がった状態で見える
	if micropost.InReplyToID != nil {
		http.Redirect(w, r, fmt.Sprintf("/microposts/%d", *micropost.InReplyToID), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Show はマイクロポスト詳細ページ（/microposts/{id}）を表示します。
// 指定した投稿とそのリプライ一覧を表示します。
func (h *MicropostHandler) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	cu := currentUser(r)
	var viewerID int64
	if cu != nil {
		viewerID = cu.ID
	}

	post, err := h.store.GetMicropostAsFeedItem(id, viewerID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	replies, _ := h.store.GetReplies(id, viewerID)
	replyCount, _ := h.store.CountReplies(id)

	data := components.MicropostDetailData{
		Title:       post.User.Name + ": " + post.Micropost.Content,
		Flash:       getFlash(r),
		LoggedIn:    cu != nil,
		CurrentUser: cu,
		CSRFToken:   middleware.CSRFTokenFromContext(r),
		Post:        *post,
		Replies:     replies,
		ReplyCount:  replyCount,
	}
	if err := components.MicropostDetail(data).Render(r.Context(), w); err != nil {
		log.Printf("render MicropostDetail: %v", err)
	}
}

// Destroy はマイクロポストを削除します。
func (h *MicropostHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	mp, err := h.store.GetMicropost(id)
	if err == nil && mp.ImagePath != "" {
		fullPath := filepath.Join(h.imageDir, mp.ImagePath)
		if removeErr := os.Remove(fullPath); removeErr != nil && !os.IsNotExist(removeErr) {
			log.Printf("Remove image %s: %v", fullPath, removeErr)
		}
	}

	if err := h.store.DeleteMicropost(id); err != nil {
		log.Printf("DeleteMicropost: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	setFlash(w, "success", "Micropost deleted")
	ref := r.Header.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}

// Index はすべてのマイクロポストを一覧表示します（管理用）。
func (h *MicropostHandler) Index(w http.ResponseWriter, r *http.Request) {
	microposts, err := h.store.AllMicroposts()
	if err != nil {
		log.Printf("AllMicroposts: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := components.MicropostPageData{
		Title:      "Microposts",
		Notice:     noticeFromRequest(r),
		Microposts: microposts,
	}
	_ = components.MicropostIndex(data).Render(r.Context(), w)
}

// processImageUpload はリクエストから画像ファイルを処理します。
// 画像がアップロードされていない場合は空文字列を返します。
func (h *MicropostHandler) processImageUpload(r *http.Request) (string, []string) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", nil
	}
	defer file.Close()

	if header.Size > maxImageSize {
		return "", []string{"Maximum file size is 5MB"}
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return "", []string{"Image must be JPEG, PNG, or GIF format"}
	}

	_, _, err = image.DecodeConfig(file)
	if err != nil {
		return "", []string{"Invalid image file"}
	}
	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), header.Size, strings.ToLower(ext))

	dst, err := os.Create(filepath.Join(h.imageDir, filename))
	if err != nil {
		return "", []string{"Failed to save image"}
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", []string{"Failed to save image"}
	}

	return filename, nil
}
