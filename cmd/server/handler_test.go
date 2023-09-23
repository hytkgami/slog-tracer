package server_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/hytkgami/slog-tracer/cmd/server"
	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	port := "18080"
	t.Setenv("PORT", port)

	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return server.Run(ctx)
	})

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/ping", port))
	if err != nil {
		t.Fatalf("failed to ping: %v", err)
	}
	defer resp.Body.Close()
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	want := "pong"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}

	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
