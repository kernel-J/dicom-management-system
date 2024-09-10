package handlers

import (
    // "encoding/json"
    "fmt"
    // "io/ioutil"
    // "mime/multipart"
    "net/http"
    // "os"
    // "path/filepath"
		
    "dicom_management_service/services"
)

type DICOMHandler struct {
    Service *services.DICOMService
}

func NewDICOMHandler(service *services.DICOMService) *DICOMHandler {
    return &DICOMHandler{
        Service: service,
    }
}

func (h *DICOMHandler) UploadDICOM(w http.ResponseWriter, r *http.Request) {
    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to upload file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    filePath, err := h.Service.UploadFile(file)
    if err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }

    w.Write([]byte(fmt.Sprintf("File uploaded successfully: %s", filePath)))
}