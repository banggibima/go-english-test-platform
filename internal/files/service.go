package files

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/banggibima/go-english-test-platform/pkg/metrics"
	"github.com/banggibima/go-english-test-platform/pkg/storage"
	"github.com/minio/minio-go/v7"
)

type Service struct {
	repository *Repository
	storage    *storage.MinIO
}

func NewService(repository *Repository, storage *storage.MinIO) *Service {
	return &Service{
		repository: repository,
		storage:    storage,
	}
}

func (s *Service) Upload(ctx context.Context, userID string, fileHeader *multipart.FileHeader) (*FileResponse, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf("uploads/%s/%d%s", userID, time.Now().UnixNano(), ext)

	_, err = s.storage.Client.PutObject(
		ctx,
		s.storage.Bucket,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		},
	)

	if err != nil {
		return nil, err
	}

	fileURL := fmt.Sprintf("http://localhost:9000/%s/%s", s.storage.Bucket, objectName)

	metadata := &File{
		UserID:       &userID,
		Bucket:       s.storage.Bucket,
		ObjectName:   objectName,
		OriginalName: fileHeader.Filename,
		ContentType:  fileHeader.Header.Get("Content-Type"),
		SizeBytes:    fileHeader.Size,
		FileURL:      &fileURL,
	}

	if err := s.repository.Create(ctx, metadata); err != nil {
		return nil, err
	}

	metrics.FilesUploadedTotal.Inc()

	response := toFileResponse(metadata)
	return &response, nil
}

func (s *Service) FindByID(ctx context.Context, id string) (*FileResponse, error) {
	file, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	response := toFileResponse(file)
	return &response, nil
}

func (s *Service) FindByUserID(ctx context.Context, userID string) ([]FileResponse, error) {
	items, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]FileResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, toFileResponse(&item))
	}

	return responses, nil
}

func toFileResponse(file *File) FileResponse {
	return FileResponse{
		ID:           file.ID,
		UserID:       file.UserID,
		Bucket:       file.Bucket,
		ObjectName:   file.ObjectName,
		OriginalName: file.OriginalName,
		ContentType:  file.ContentType,
		SizeBytes:    file.SizeBytes,
		FileURL:      file.FileURL,
	}
}
