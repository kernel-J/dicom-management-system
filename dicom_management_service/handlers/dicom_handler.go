package handlers

import (
    "encoding/json"
    "fmt"
    // "io/ioutil"
    // "mime/multipart"
    "net/http"
    // "os"
    // "path/filepath"
    
    "github.com/gorilla/mux"
    "github.com/suyashkumar/dicom"


    "dicom_management_service/services"
)

type DICOMHandler struct {
    Service *services.DICOMService
}

type DICOMElement struct {
    Id      string
    Tag     string
    Element dicom.Element
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

func (h *DICOMHandler) GetDICOMAttributes(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    tag := vars["tag"]

    fmt.Println("id = %s", id)
    if id == "" || tag == "" {
        http.Error(w, "Missing id or tag query parameter", http.StatusBadRequest)
        return
    }

    value, err := h.Service.GetDICOMAttributes(id, tag)
    if err != nil || value == nil {
        http.Error(w, "Failed to get DICOM header", http.StatusInternalServerError)
        return
    }

    ret := DICOMElement{Id: id, Tag: tag, Element: *value}
    json.NewEncoder(w).Encode(ret)
}

func (h *DICOMHandler) GetDICOM(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    if id == "" {
        http.Error(w, "Missing id query parameter", http.StatusBadRequest)
        return
    }

    pngPath, err := h.Service.GetDICOM(id)
    if err != nil {
        http.Error(w, "Failed to convert to PNG", http.StatusInternalServerError)
        return
    }

    http.ServeFile(w, r, pngPath)
}