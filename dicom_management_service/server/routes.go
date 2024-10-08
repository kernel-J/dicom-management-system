package server

import (
	"github.com/gorilla/mux"

	"dicom_management_service/handlers"
)

func RegisterRoutes(handler *handlers.DICOMHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/upload", handler.UploadDICOM).Methods("POST")
	router.HandleFunc("/dicom/{id}/{tag}", handler.GetDICOMAttributes).Methods("GET")
	router.HandleFunc("/convert/{id}", handler.ConvertDICOM).Methods("GET")

	return router
}