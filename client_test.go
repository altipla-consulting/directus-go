package directus

import (
	"log/slog"
	"os"
	"testing"
)

func initClient(t *testing.T) *Client {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})

	if os.Getenv("DIRECTUS_TOKEN") == "" {
		t.Skip("DIRECTUS_TOKEN not set")
	}
	return NewClient("http://localhost:8055", os.Getenv("DIRECTUS_TOKEN"), WithLogger(slog.New(handler)), WithBodyLogger())
}
