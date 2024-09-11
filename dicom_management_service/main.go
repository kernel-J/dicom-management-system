package main

import (
	"log"
	"os"
	"github.com/google/uuid"

	"dicom_management_service/config"
	"dicom_management_service/server"
	"dicom_management_service/handlers"
	"dicom_management_service/services"
)

func main() {
	// Load configuration
	cfg := config.New()
	cfg.Validate()

	if _, err := os.Stat(cfg.UploadDir); os.IsNotExist(err) {
		os.Mkdir(cfg.UploadDir, os.ModePerm)
	}

	UUIDGenerator := uuid.NewRandom
	fileStorageService := services.NewFileStorageService(cfg.UploadDir)
	dicomService := services.NewDICOMService(cfg.UploadDir, fileStorageService, UUIDGenerator)
	dicomHandler := handlers.NewDICOMHandler(dicomService)

	srv := server.NewServer(cfg)
	mux := server.RegisterRoutes(dicomHandler)
	log.Printf("Starting server on port %s\n", cfg.Port)
	srv.Run(mux)
}