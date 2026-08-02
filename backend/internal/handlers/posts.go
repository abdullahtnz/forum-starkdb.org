package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abdullahtnz/forum-starkdb/internal/middleware"
	"github.com/abdullahtnz/forum-starkdb/internal/models"
	"github.com/abdullahtnz/forum-starkdb/internal/services"
)

type PostHandler struct {
	svc *services.PostService
}

func NewPostHandler(svc *services.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.svc.Create(r.Context(), claims.UserID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	viewerID := ""
	if claims := middleware.GetUserClaims(r); claims != nil {
		viewerID = claims.UserID
	}
	post, err := h.svc.GetByID(r.Context(), postID, viewerID)
	if err != nil {
		writeError(w, http.StatusNotFound, "post not found")
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	page, pageSize := parsePagination(r)
	sort := r.URL.Query().Get("sort")
	tag := r.URL.Query().Get("tag")
	search := r.URL.Query().Get("q")
	viewerID := ""
	if claims := middleware.GetUserClaims(r); claims != nil {
		viewerID = claims.UserID
	}
	result, err := h.svc.List(r.Context(), page, pageSize, sort, tag, search, viewerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list posts")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	postID := r.PathValue("id")
	var req models.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.svc.Update(r.Context(), postID, claims.UserID, req)
	if err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, post)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	postID := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), postID, claims.UserID, claims.IsAdmin); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	postID := r.PathValue("id")
	liked, err := h.svc.ToggleLike(r.Context(), postID, claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to toggle like")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}

func (h *PostHandler) PinPost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	var req struct {
		Pinned bool `json:"pinned"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.svc.PinPost(r.Context(), postID, req.Pinned); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to pin post")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *PostHandler) ClosePost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	var req struct {
		Closed bool `json:"closed"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.svc.ClosePost(r.Context(), postID, req.Closed); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to close post")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func parsePagination(r *http.Request) (int, int) {
	page := 1
	pageSize := 20
	if p := r.URL.Query().Get("page"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			page = n
		}
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if n, err := strconv.Atoi(ps); err == nil && n > 0 && n <= 50 {
			pageSize = n
		}
	}
	return page, pageSize
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}
