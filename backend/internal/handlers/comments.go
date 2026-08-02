package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/abdullahtnz/forum-starkdb/internal/middleware"
	"github.com/abdullahtnz/forum-starkdb/internal/models"
	"github.com/abdullahtnz/forum-starkdb/internal/services"
)

type CommentHandler struct {
	svc *services.CommentService
}

func NewCommentHandler(svc *services.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	viewerID := ""
	if claims := middleware.GetUserClaims(r); claims != nil {
		viewerID = claims.UserID
	}
	comments, err := h.svc.GetByPost(r.Context(), postID, viewerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load comments")
		return
	}
	writeJSON(w, http.StatusOK, comments)
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	postID := r.PathValue("id")
	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	comment, err := h.svc.Create(r.Context(), postID, claims.UserID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

func (h *CommentHandler) Reply(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	parentID := r.PathValue("id")
	var body struct {
		Content string `json:"content"`
		PostID  string `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req := models.CreateCommentRequest{
		Content:  body.Content,
		ParentID: &parentID,
	}
	comment, err := h.svc.Create(r.Context(), body.PostID, claims.UserID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	commentID := r.PathValue("id")
	var req models.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	comment, err := h.svc.Update(r.Context(), commentID, claims.UserID, req)
	if err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, comment)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	commentID := r.PathValue("id")
	if err := h.svc.Delete(r.Context(), commentID, claims.UserID, claims.IsAdmin); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CommentHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	commentID := r.PathValue("id")
	liked, err := h.svc.ToggleLike(r.Context(), commentID, claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to toggle like")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}
