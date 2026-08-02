package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/abdullahtnz/forum-starkdb/internal/middleware"
	"github.com/abdullahtnz/forum-starkdb/internal/services"
)

type AdminHandler struct {
	svc *services.AdminService
}

func NewAdminHandler(svc *services.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	page, pageSize := parsePagination(r)
	users, total, err := h.svc.GetUsers(r.Context(), page, pageSize)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get users")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"users":       users,
		"total_count": total,
		"page":        page,
		"page_size":   pageSize,
	})
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if err := h.svc.DeleteUser(r.Context(), userID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AdminHandler) ToggleAdmin(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	var req struct {
		MakeAdmin bool `json:"make_admin"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.svc.ToggleAdmin(r.Context(), userID, req.MakeAdmin); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update admin status")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AdminHandler) GetBadWords(w http.ResponseWriter, r *http.Request) {
	words, err := h.svc.GetBadWords(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get bad words")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"words": words})
}

func (h *AdminHandler) AddBadWord(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Word string `json:"word"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.svc.AddBadWord(r.Context(), req.Word); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add word")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (h *AdminHandler) RemoveBadWord(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Word string `json:"word"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.svc.RemoveBadWord(r.Context(), req.Word); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to remove word")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AdminHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if err := h.svc.DeletePost(r.Context(), postID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete post")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AdminHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := r.PathValue("id")
	if err := h.svc.DeleteComment(r.Context(), commentID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete comment")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetUserClaims(r)
		if claims == nil || !claims.IsAdmin {
			writeError(w, http.StatusForbidden, "admin access required")
			return
		}
		next(w, r)
	}
}
