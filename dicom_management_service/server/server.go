package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"dicom_management_service/config"
	"log"
)

type Server struct {
	config    *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config:    config,
	}
}

func (s *Server) Run(mux *http.ServeMux) {
	srv := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
