package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg":                   true,
	"image/png":                    true,
	"image/gif":                    true,
	"image/webp":                   true,
	"image/svg+xml":                true,
	"application/pdf":              true,
	"text/plain":                   true,
	"text/csv":                     true,
	"application/zip":              true,
	"application/x-zip-compressed": true,
}

var blockedExtensions = []string{
	".exe", ".sh", ".bat", ".cmd", ".msi",
	".js", ".ts", ".php", ".py", ".rb", ".pl",
	".dll", ".so", ".dylib",
	".vbs", ".ps1", ".wsf",
}

func ValidateFileUpload(file multipart.File, header *multipart.FileHeader, maxSize int64) error {
	if header.Size > maxSize {
		return fmt.Errorf("file too large: %d bytes (max %d bytes)", header.Size, maxSize)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	for _, blocked := range blockedExtensions {
		if ext == blocked {
			return fmt.Errorf("file type not allowed: %s", ext)
		}
	}

	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	mimeType := detectMimeType(buf)
	if !allowedMimeTypes[mimeType] {
		return fmt.Errorf("file type not allowed: %s", mimeType)
	}

	return nil
}

func detectMimeType(buf []byte) string {
	if len(buf) >= 2 && buf[0] == 0xFF && buf[1] == 0xD8 {
		return "image/jpeg"
	}
	if len(buf) >= 4 && buf[0] == 0x89 && buf[1] == 'P' && buf[2] == 'N' && buf[3] == 'G' {
		return "image/png"
	}
	if len(buf) >= 4 && buf[0] == 'G' && buf[1] == 'I' && buf[2] == 'F' {
		return "image/gif"
	}
	if len(buf) >= 4 && buf[0] == 'R' && buf[1] == 'I' && buf[2] == 'F' && buf[3] == 'F' {
		return "image/webp"
	}
	if len(buf) >= 4 && buf[0] == '%' && buf[1] == 'P' && buf[2] == 'D' && buf[3] == 'F' {
		return "application/pdf"
	}
	if len(buf) >= 2 && buf[0] == 'P' && buf[1] == 'K' {
		return "application/zip"
	}
	if len(buf) >= 4 && buf[0] == '<' && buf[1] == '?' && buf[2] == 'x' && buf[3] == 'm' {
		return "image/svg+xml"
	}
	if len(buf) >= 5 && buf[0] == '<' && buf[1] == 's' && buf[2] == 'v' && buf[3] == 'g' {
		return "image/svg+xml"
	}
	return "application/octet-stream"
}
