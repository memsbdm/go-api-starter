package helpers

import (
	"errors"
	"go-starter/internal/domain"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

// ImageExtensions is a list of allowed image file extensions.
var (
	ImageExtensions = []string{".png", ".jpg", ".jpeg"}
)

// MultipartFormParser is a helper struct for parsing multipart form data.
type MultipartFormParser struct {
	maxMemory         int64
	allowedExtensions []string
}

// NewMultipartFormParser creates a new MultipartFormParser instance.
func NewMultipartFormParser(maxMemory int64, allowedExtensions []string) *MultipartFormParser {
	normalizedExtensions := make([]string, len(allowedExtensions))
	for i, ext := range allowedExtensions {
		normalizedExtensions[i] = strings.ToLower(strings.TrimPrefix(ext, "."))
	}
	return &MultipartFormParser{
		maxMemory:         maxMemory,
		allowedExtensions: normalizedExtensions,
	}
}

// Parse parses the multipart form data from the request.
func (p *MultipartFormParser) Parse(r *http.Request) error {
	if err := r.ParseMultipartForm(p.maxMemory); err != nil {
		switch {
		case errors.Is(err, http.ErrNotMultipart):
			return domain.ErrInvalidMultipartForm
		case errors.Is(err, http.ErrMissingBoundary):
			return domain.ErrMissingBoundary
		default:
			return err
		}
	}
	return nil
}

func (p *MultipartFormParser) GetFile(r *http.Request, fieldName string, maxFileSize int64) (*multipart.File, *multipart.FileHeader, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	if header.Size > maxFileSize {
		return nil, nil, domain.ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	ext = strings.TrimPrefix(ext, ".")

	if !p.isExtensionAllowed(ext) {
		return nil, nil, domain.ErrInvalidFileType
	}

	return &file, header, nil
}

// isExtensionAllowed checks if the file extension is allowed.
func (p *MultipartFormParser) isExtensionAllowed(ext string) bool {
	if len(p.allowedExtensions) == 0 {
		return true
	}
	return slices.Contains(p.allowedExtensions, ext)
}

// GenerateFileKey generates a file key for the file.
func GenerateFileKey(filePrefix string, filename string) string {
	return filePrefix + "." + getFileExtension(filename)
}

// getFileExtension returns the file extension from the filename.
// Returns an empty string if no extension is found.
func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) <= 1 {
		return ""
	}
	return parts[len(parts)-1]
}
