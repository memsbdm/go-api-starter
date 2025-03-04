package ports

import (
	"context"
	"go-starter/internal/domain/entities"
	"io"
)

// FileUploadService is a service that uploads files to a file upload service.
type FileUploadService interface {
	UploadAvatar(ctx context.Context, userID entities.UserID, filename string, body io.Reader) (string, error)
}

// FileUploadAdapter is an adapter for the FileUploadService interface.
type FileUploadAdapter interface {
	Upload(ctx context.Context, key string, body io.Reader) (string, error)
}
