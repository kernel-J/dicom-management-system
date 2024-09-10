package services

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "image/jpeg"

    "github.com/google/uuid"
    "github.com/suyashkumar/dicom"
    "github.com/suyashkumar/dicom/pkg/tag"
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
    fileId, err := uuid.NewRandom()
    if err != nil {
        return "", err
    }

    fileName := fileId.String() + ".dcm"
    filePath := filepath.Join(s.UploadDir, fileName)
    out, err := os.Create(filePath)
    if err != nil {
        return "", err
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        return "", err
    }

    return fileName, nil
}

/*
** copied from https://github.com/suyashkumar/dicom/blob/41fd08ede4de358a6a4514370cfdcb012b4008d6/pkg/tag/tag.go#L236
*/
func (s *DICOMService) parse(dicomTag string) (tag.Tag, error) {
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

func (s *DICOMService) GetDICOMAttributes(fileId string, dicomTag string) (*dicom.Element, error) {
    filePath := s.UploadDir + "/" + fileId + ".dcm"
    fmt.Println(filePath)

    dataSet, err := dicom.ParseFile(filePath, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to parse", err)
    }

    parsedTag, err := s.parse(dicomTag)
    if err != nil {
        fmt.Println("Error parsing tag:", err)
        return nil, err
    }

    elem, err := dataSet.FindElementByTagNested(parsedTag)
    if err != nil {
        return nil, fmt.Errorf("failed to parse", err)
    }

    return elem, nil
}


func (s *DICOMService) GetDICOM(fileId string) (string, error) {
    filePath := s.UploadDir + "/" + fileId + ".dcm"
    fmt.Println(filePath)

    dataSet, err := dicom.ParseFile(filePath, nil)
    if err != nil {
        return "", fmt.Errorf("failed to parse", err)
    }
    
    pixelDataElement, err := dataSet.FindElementByTag(tag.PixelData)
    if err != nil {
        return "", err
    }

    pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	for i, fr := range pixelDataInfo.Frames {
		img, _ := fr.GetImage() // The Go image.Image for this frame
		f, _ := os.Create(fmt.Sprintf("%s/%s-%d.jpg", s.UploadDir, fileId, i))
		_ = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
		_ = f.Close()
    }
    
    return "", nil
}