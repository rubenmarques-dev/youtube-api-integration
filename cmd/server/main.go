package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "github.com/rubenmarques/youtube-api-integration/docs"
	httphandler "github.com/rubenmarques/youtube-api-integration/internal/adapter/http"
	youtubeadapter "github.com/rubenmarques/youtube-api-integration/internal/adapter/youtube"
	httpSwagger "github.com/swaggo/http-swagger"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// @title           YouTube API Integration
// @version         1.0
// @description     POC service for fetching YouTube video metadata.
// @host            localhost:8080
// @BasePath        /
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
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	logger.Info("starting server", "port", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(mux)); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
