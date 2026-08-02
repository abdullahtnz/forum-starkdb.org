package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/abdullahtnz/forum-starkdb/internal/middleware"
	"github.com/abdullahtnz/forum-starkdb/internal/services"
	"github.com/abdullahtnz/forum-starkdb/internal/utils"
	"github.com/google/uuid"
)

type UploadHandler struct {
	svc *services.UploadService
}

func NewUploadHandler(svc *services.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	if err := r.ParseMultipartForm(h.svc.GetMaxSize()); err != nil {
		writeError(w, http.StatusBadRequest, "file too large")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	if err := utils.ValidateFileUpload(file, header, h.svc.GetMaxSize()); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	ext := filepath.Ext(header.Filename)
	newName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().UnixNano(), ext)
	filePath := filepath.Join(h.svc.GetUploadDir(), newName)

	if err := os.MkdirAll(h.svc.GetUploadDir(), 0755); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create upload directory")
		return
	}

	dst, err := os.Create(filePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := file.Seek(0, 0); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read file")
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to write file")
		return
	}

	response := map[string]string{
		"id":        newName,
		"file_name": header.Filename,
		"file_path": "/uploads/" + newName,
		"mime_type": header.Header.Get("Content-Type"),
		"size":      fmt.Sprintf("%d", header.Size),
	}
	writeJSON(w, http.StatusCreated, response)
}

type TagHandler struct {
	svc *services.TagService
}

func NewTagHandler(svc *services.TagService) *TagHandler {
	return &TagHandler{svc: svc}
}

func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	tags, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load tags")
		return
	}
	writeJSON(w, http.StatusOK, tags)
}
