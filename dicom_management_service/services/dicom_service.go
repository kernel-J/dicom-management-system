package services

import (
    "fmt"
    "io"
    "strconv"
    "strings"
    "image/png"
    "log"

    "github.com/google/uuid"
    "github.com/suyashkumar/dicom"
    "github.com/suyashkumar/dicom/pkg/tag"
)

//go:generate mockgen -source=dicom_service.go -destination=../mocks/services/dicom_service.go

type UUIDGenerator func() (uuid.UUID, error)

type DICOMService interface {
    UploadFile(file io.Reader) (string, error)
    GetDICOMAttributes(fileId string, dicomTag string) (*dicom.Element, error)
    ConvertDICOMToPNG(fileId string) (string, error)
}

type dicomService struct {
    UploadDir           string
    FileService         FileStorageService
    UuidGenerator       UUIDGenerator
}

func NewDICOMService(u string, fs FileStorageService, uuidGen UUIDGenerator) *dicomService {
    return &dicomService{
        UploadDir: u,
        FileService: fs,
        UuidGenerator: uuidGen,
    }
}

func (s *dicomService) UploadFile(file io.Reader) (string, error) {
    fileId, err := s.UuidGenerator()
    if err != nil {
        log.Printf("Error: failed to upload file\n")
        return "", err
    }

    fileName := fileId.String() + ".dcm"
    out, err := s.FileService.Create(fileName)
    if err != nil {
        log.Printf("Error: failed to upload file, error: %s, fileId: %s\n", err, fileId)
        return "", err
    }
    defer out.Close()

    err = s.FileService.Copy(out, file, fileId.String())
    if err != nil {
        log.Printf("Error: failed to upload file, error: %s, fileId: %s\n", err, fileId)
        return "", err
    }

    log.Printf("file was successfully uploaded, fileId: %s\n", fileId)
    return fileName, nil
}

/*
** copied from https://github.com/suyashkumar/dicom/blob/41fd08ede4de358a6a4514370cfdcb012b4008d6/pkg/tag/tag.go#L236
*/
func (s *dicomService) parse(dicomTag string) (tag.Tag, error) {
	parts := strings.Split(strings.Trim(dicomTag, "()"), ",")
	group, err := strconv.ParseInt(parts[0], 16, 0)
	if err != nil {
		return tag.Tag{}, err
	}
	elem, err := strconv.ParseInt(parts[1], 16, 0)
	if err != nil {
		return tag.Tag{}, err
	}
	return tag.Tag{Group: uint16(group), Element: uint16(elem)}, nil
}

func (s *dicomService) GetDICOMAttributes(fileId string, dicomTag string) (*dicom.Element, error) {
    filePath := s.UploadDir + "/" + fileId + ".dcm"
    dataSet, err := dicom.ParseFile(filePath, nil)
    if err != nil {
        log.Printf("Error: failed to get attributes, error: %s, fileId: %s\n", err, fileId)
        return nil, err
    }

    parsedTag, err := s.parse(dicomTag)
    if err != nil {
        log.Printf("Error: failed to get attributes, error: %s, fileId: %s\n", err, fileId)
        return nil, err
    }

    elem, err := dataSet.FindElementByTagNested(parsedTag)
    if err != nil {
        log.Printf("Error: failed to get attributes, error: %s, fileId: %s\n", err, fileId)
        return nil, err
    }

    log.Printf("Attributes were successfully retreived, fileId: %s\n", fileId)
    return elem, nil
}

/*
** TO DO: Send all the frames
*/
func (s *dicomService) ConvertDICOMToPNG(fileId string) (string, error) {
    filePath := s.UploadDir + "/" + fileId + ".dcm"

    dataSet, err := dicom.ParseFile(filePath, nil)
    if err != nil {
        log.Printf("Error: failed to convert dicom to png, error: %s, fileId: %s\n", err, fileId)
        return "", err
    }
    
    pixelDataElement, err := dataSet.FindElementByTag(tag.PixelData)
    if err != nil {
        log.Printf("Error: failed to convert dicom to png, error: %s, fileId: %s\n", err, fileId)
        return "", err
    }

    pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	for i, fr := range pixelDataInfo.Frames {
		img, _ := fr.GetImage()
		file, _ := s.FileService.Create(fmt.Sprintf("%s-%d.png", fileId, i))
		_ = png.Encode(file, img)
		_ = file.Close()
    }
    
    log.Printf("File was successfully converted to PNG, fileId: %s\n", fileId)
    return fmt.Sprintf("%s/%s-%d.png", s.UploadDir, fileId, 0), nil
}