package server

import (
	"net/http"
	"ImageManagementService/handler"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", handler.ImageUploadHandler)

	return mux
}