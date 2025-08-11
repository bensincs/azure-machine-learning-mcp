package main

import (
	"log"

	"microsoft.com/aml-mcp/internal/server"
)

func main() {
	config := server.Config{
		Name:    "Azure Machine Learning SDK",
		Version: "1.0.0",
	}

	s := server.New(config)
	if err := s.Serve(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}