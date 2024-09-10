package services

import (
    "fmt"
    "io"
    "os"
    "path/filepath"

    "github.com/google/uuid"
)

type DICOMService struct {
    UploadDir string
}

func NewDICOMService(uploadDir string) *DICOMService {
    return &DICOMService{
        UploadDir: uploadDir,
    }
}

func (s *DICOMService) UploadFile(file io.Reader) (string, error) {
    file_id, err := uuid.NewRandom()
    if err != nil {
        return "", err
    }

    filePath := filepath.Join(s.UploadDir, file_id.String() + ".dcm")
    out, err := os.Create(filePath)
    if err != nil {
        return "", err
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        return "", err
    }

    return filePath, nil
}