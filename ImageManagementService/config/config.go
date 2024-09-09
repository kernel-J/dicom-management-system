package config

import (
	"log"
	"os"
)

type Config struct {
	Port     string
}

func New() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// database := os.Getenv("DATABASE_URL")
	// if database == "" {
	// 	database = "postgres://user:password@localhost:5432/mydatabase"
	// }

	return &Config{
		Port:     port,
	}
}

func (c *Config) Validate() {
	if c.Port == "" {
		log.Fatal("Port is required")
	}
	// if c.Database == "" {
	// 	log.Fatal("Database URL is required")
	// }
}
