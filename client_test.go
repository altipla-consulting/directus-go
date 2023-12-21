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
	return NewClient("https://compostela.admin.onetbooking.com", os.Getenv("DIRECTUS_TOKEN"), WithLogger(slog.New(handler)), WithBodyLogger())
}
