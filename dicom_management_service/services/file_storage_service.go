package services

import (
	"io"
	"os"
	"log"
	"path/filepath"
)

//go:generate mockgen -source=file_storage_service.go -destination=../mocks/services/file_storage_service.go

type FileStorageService interface {
	Create(filename string) (*os.File, error)
	Copy(destination io.Writer, source io.Reader, filename string) error
}

type fileStorageService struct {
	UploadDir string
}

func NewFileStorageService(uploadDir string) *fileStorageService {
	return &fileStorageService{
		UploadDir: uploadDir,
	}
}

func (s *fileStorageService) Create(filename string) (*os.File, error) {
	filePath := filepath.Join(s.UploadDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error: failed to create file, error: %s, filename: %s\n", err, filename)
		return nil, err
	}

	log.Printf("File was successfully created, filename: %s\n", filename)
	return file, nil
}

func (s *fileStorageService) Copy(destination io.Writer, source io.Reader, filename string) error {
	_, err := io.Copy(destination, source)
	if err != nil {
		log.Printf("Error: failed to copy file, error: %s, filename: %s\n", err, filename)
		return err
	}

	log.Printf("File was successfully copied, filename: %s\n", filename)
	return nil
}
