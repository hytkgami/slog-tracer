package main

import (
	"context"
	"log"

	"github.com/hytkgami/slog-tracer/cmd/server"
)

func main() {
	if err := server.Run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}
