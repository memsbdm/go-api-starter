package ports

import (
	"context"
	"go-starter/internal/domain/entities"
	"io"
)

// FileUploadService is a service that uploads files to a file upload service.
type FileUploadService interface {
	// UploadAvatar uploads a user avatar to the S3 bucket.
	// Returns the URL of the uploaded file or an error if the upload fails.
	UploadAvatar(ctx context.Context, userID entities.UserID, filename string, body io.Reader) (string, error)
	// DeleteAvatar deletes a user avatar from the S3 bucket.
	// Returns an error if the deletion fails.
	DeleteAvatar(ctx context.Context, userID entities.UserID, avatarURL string) error
}

// FileUploadAdapter is an adapter for the FileUploadService interface.
type FileUploadAdapter interface {
	// Upload uploads a file to the S3 bucket.
	// Returns the URL of the uploaded file or an error if the upload fails.
	Upload(ctx context.Context, key string, body io.Reader) (string, error)
	// Delete deletes a file from the S3 bucket.
	// Returns an error if the deletion fails.
	Delete(ctx context.Context, key string) error
}
