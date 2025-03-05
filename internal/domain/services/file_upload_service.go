package services

import (
	"context"
	"go-starter/internal/adapters/http/helpers"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"io"
	"strings"
)

// FileUploadService is a service that uploads files to a file upload service.
type FileUploadService struct {
	adapter ports.FileUploadAdapter
}

// NewFileUploadService creates a new instance of FileUploadService.
func NewFileUploadService(adapter ports.FileUploadAdapter) *FileUploadService {
	return &FileUploadService{
		adapter: adapter,
	}
}

// UserAvatarPath is the path to the user avatar directory.
const UserAvatarPath = "avatars"

// UploadAvatar uploads a user avatar to the file upload service.
// Returns the URL of the uploaded file or an error if the upload fails.
func (s *FileUploadService) UploadAvatar(ctx context.Context, userID entities.UserID, filename string, body io.Reader) (string, error) {
	key := UserAvatarPath + "/" + helpers.GenerateFileKey(userID.String(), filename)
	url, err := s.adapter.Upload(ctx, key, body)
	if err != nil {
		return "", domain.ErrFileUpload
	}
	return url, nil
}

// DeleteAvatar deletes a user avatar from the file upload service.
// Returns an error if the deletion fails.
func (s *FileUploadService) DeleteAvatar(ctx context.Context, userID entities.UserID, avatarURL string) error {
	split := strings.Split(avatarURL, ".")
	key := UserAvatarPath + "/" + userID.String() + "." + split[len(split)-1]
	err := s.adapter.Delete(ctx, key)
	if err != nil {
		return domain.ErrFileUpload
	}
	return nil
}
