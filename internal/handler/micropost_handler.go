package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sofuejin0121/toy_app_go/internal/model"
	"github.com/sofuejin0121/toy_app_go/internal/store"
	"github.com/sofuejin0121/toy_app_go/web/components"
)

// MicropostHandler はマイクロポストリソースのHTTPハンドラーです。
type MicropostHandler struct {
	store *store.Store
}

// NewMicropostHandler は新しいMicropostHandlerを返します。
func NewMicropostHandler(store *store.Store) *MicropostHandler {
	return &MicropostHandler{store: store}
}

func (h *MicropostHandler) formData(
	title string,
	micropost model.Micropost,
	errors []string,
	action string,
	submitLabel string,
) (components.MicropostPageData, error) {
	users, err := h.store.AllUsers()
	if err != nil {
		return components.MicropostPageData{}, err
	}
	return components.MicropostPageData{
		Title:       title,
		Micropost:   micropost,
		Users:       users,
		Errors:      errors,
		Action:      action,
		SubmitLabel: submitLabel,
	}, nil
}

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

func (h *MicropostHandler) Show(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	micropost, err := h.store.GetMicropost(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data := components.MicropostPageData{
		Title:     "Micropost",
		Notice:    noticeFromRequest(r),
		Micropost: *micropost,
	}
	_ = components.MicropostShow(data).Render(r.Context(), w)
}

func (h *MicropostHandler) New(w http.ResponseWriter, r *http.Request) {
	data, err := h.formData("New micropost", model.Micropost{}, nil, "/microposts", "Create Micropost")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = components.MicropostNew(data).Render(r.Context(), w)
}

func (h *MicropostHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	micropost, err := h.store.GetMicropost(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	data, err := h.formData(
		"Editing micropost",
		*micropost,
		nil,
		fmt.Sprintf("/microposts/%d", micropost.ID),
		"Update Micropost",
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = components.MicropostEdit(data).Render(r.Context(), w)
}

func (h *MicropostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	micropost := model.Micropost{
		Content: r.FormValue("content"),
		UserID:  userID,
	}
	if errors := micropost.Validate(); len(errors) > 0 {
		data, err := h.formData("New micropost", micropost, errors, "/microposts", "Create Micropost")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_ = components.MicropostNew(data).Render(r.Context(), w)
		return
	}
	if err := h.store.CreateMicropost(&micropost); err != nil {
		data, dataErr := h.formData("New micropost", micropost, []string{err.Error()}, "/microposts", "Create Micropost")
		if dataErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_ = components.MicropostNew(data).Render(r.Context(), w)
		return
	}

	redirectWithNotice(w, r, fmt.Sprintf("/microposts/%d", micropost.ID), "Micropost was successfully created.")
}

func (h *MicropostHandler) Update(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	micropost, err := h.store.GetMicropost(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	userID, _ := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	micropost.Content = r.FormValue("content")
	micropost.UserID = userID
	if errors := micropost.Validate(); len(errors) > 0 {
		data, dataErr := h.formData(
			"Editing micropost",
			*micropost,
			errors,
			fmt.Sprintf("/microposts/%d", micropost.ID),
			"Update Micropost",
		)
		if dataErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_ = components.MicropostEdit(data).Render(r.Context(), w)
		return
	}
	if err := h.store.UpdateMicropost(micropost); err != nil {
		data, dataErr := h.formData(
			"Editing micropost",
			*micropost,
			[]string{err.Error()},
			fmt.Sprintf("/microposts/%d", micropost.ID),
			"Update Micropost",
		)
		if dataErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_ = components.MicropostEdit(data).Render(r.Context(), w)
		return
	}

	redirectWithNotice(w, r, fmt.Sprintf("/microposts/%d", micropost.ID), "Micropost was successfully updated.")
}

func (h *MicropostHandler) Destroy(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := h.store.DeleteMicropost(id); err != nil {
		http.NotFound(w, r)
		return
	}
	redirectWithNotice(w, r, "/microposts", "Micropost was successfully destroyed.")
}
