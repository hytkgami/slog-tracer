package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

const defaultPort = "8080"

func Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}))

	s := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux,
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		log.Println("listening on port", port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %v", err)
	}
	return eg.Wait()
}
