package services_test

import (
	"testing"
	"io"
	"os"
	"errors"

	"github.com/google/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	
	mocks "dicom_management_service/mocks/services"
	"dicom_management_service/services"
)

func mockUUIDGenerator() (uuid.UUID, error) {
	return uuid.Must(uuid.Parse("69359037-9599-48e7-b8f2-48393c019135")), nil
}

func TestDICOMService_NewDICOMService(t *testing.T) {
	ctrl := gomock.NewController(t)
	uploadDir := "/path/to/upload/dir"
	mockFileStorageService := mocks.NewMockFileStorageService(ctrl)

	service := services.NewDICOMService(uploadDir, mockFileStorageService, mockUUIDGenerator)

	assert.NotNil(t, service)
	assert.Equal(t, uploadDir, service.UploadDir)
	assert.Equal(t, mockFileStorageService, service.FileService)
}

func TestDICOMService_UploadFile(t *testing.T) {
	t.Run("Should successfully upload a file", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		uploadDir := "/path/to/upload/dir"
		mockFileStorageService := mocks.NewMockFileStorageService(ctrl)

		service := services.NewDICOMService(uploadDir, mockFileStorageService, mockUUIDGenerator)
		fileContent := io.NopCloser(nil) 

		mockFileStorageService.EXPECT().Create(gomock.Any()).Return(&os.File{}, nil).Times(1)
		mockFileStorageService.EXPECT().Copy(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		filePath, err := service.UploadFile(fileContent)
		assert.NoError(t, err)
		assert.Equal(t, filePath, "69359037-9599-48e7-b8f2-48393c019135.dcm")
	})

	t.Run("Should return error when creating the file fails", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		uploadDir := "/path/to/upload/dir"
		mockFileStorageService := mocks.NewMockFileStorageService(ctrl)
		mockError := errors.New("bogus error")

		service := services.NewDICOMService(uploadDir, mockFileStorageService, mockUUIDGenerator)
		fileContent := io.NopCloser(nil) 

		mockFileStorageService.EXPECT().Create(gomock.Any()).Return(nil, mockError).Times(1)

		filePath, err := service.UploadFile(fileContent)
		assert.Error(t, err, mockError)
		assert.Equal(t, filePath, "")
	})

	t.Run("Should return error when copying the file fails", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		uploadDir := "/path/to/upload/dir"
		mockFileStorageService := mocks.NewMockFileStorageService(ctrl)
		mockError := errors.New("bogus error")

		service := services.NewDICOMService(uploadDir, mockFileStorageService, mockUUIDGenerator)
		fileContent := io.NopCloser(nil) 

		mockFileStorageService.EXPECT().Create(gomock.Any()).Return(&os.File{}, nil).Times(1)
		mockFileStorageService.EXPECT().Copy(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockError).Times(1)

		filePath, err := service.UploadFile(fileContent)
		assert.Error(t, err, mockError)
		assert.Equal(t, filePath, "")
	})

}

func TestDICOMService_GetDICOMAttributes(t *testing.T) {
	t.Run("Should successfully return the dicom attributes", func (t *testing.T) {
	})

	t.Run("Should return error when parsing the file fails", func (t *testing.T) {
	})

	t.Run("Should return error when parsing the tag fails", func (t *testing.T) {
	})
}

func TestDICOMService_ConvertDICOMToPNG(t *testing.T) {
	t.Run("Should successfully convert to png", func (t *testing.T) {
	})

	t.Run("Should return error when parsing the file fails", func (t *testing.T) {
	})

	t.Run("Should return error when finding elements fails", func (t *testing.T) {
	})
}