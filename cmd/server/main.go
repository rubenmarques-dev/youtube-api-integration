package main

import (
	"log/slog"
	"net/http"
	"os"

	httphandler "github.com/rubenmarques/youtube-api-integration/internal/adapter/http"
	youtubeadapter "github.com/rubenmarques/youtube-api-integration/internal/adapter/youtube"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		logger.Error("YOUTUBE_API_KEY environment variable is required")
		os.Exit(1)
	}

	ytClient, err := youtubeadapter.New(apiKey)
	if err != nil {
		logger.Error("failed to create YouTube client", "error", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	h := httphandler.New(ytClient, logger)
	h.RegisterRoutes(mux)

	logger.Info("starting server", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
