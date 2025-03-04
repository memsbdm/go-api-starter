package mocks

import (
	"context"
	"io"
)

// FileUploadAdapterMock is a mock implementation of the ports.FileUploadAdapter interface.
type FileUploadAdapterMock struct{}

// NewFileUploadAdapterMock creates a new FileUploadAdapterMock instance.
func NewFileUploadAdapterMock() *FileUploadAdapterMock {
	return &FileUploadAdapterMock{}
}

// Upload uploads a file to the file upload service.
func (f *FileUploadAdapterMock) Upload(_ context.Context, key string, _ io.Reader) (string, error) {
	return "https://example.com/" + key, nil
}
