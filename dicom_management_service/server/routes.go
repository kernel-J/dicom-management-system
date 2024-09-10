package server

import (
	"net/http"
	"dicom_management_service/handlers"
)

func RegisterRoutes(handler *handlers.DICOMHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", handler.UploadDICOM)

	return mux
}