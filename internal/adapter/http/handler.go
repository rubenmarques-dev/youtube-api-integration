package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/rubenmarques/youtube-api-integration/internal/domain"
)

// Handler holds the HTTP handler dependencies.
type Handler struct {
	repo   domain.VideoRepository
	logger *slog.Logger
}

// New creates a new HTTP Handler.
func New(repo domain.VideoRepository, logger *slog.Logger) *Handler {
	return &Handler{repo: repo, logger: logger}
}

// RegisterRoutes registers the handler's routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/video", h.getVideo)
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) getVideo(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("incoming request", "method", r.Method, "path", r.URL.Path)

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "missing required query parameter: id"})
		return
	}

	video, err := h.repo.GetVideoByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.logger.Warn("video not found", "id", id)
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "video not found"})
			return
		}
		h.logger.Error("failed to fetch video", "id", id, "error", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to fetch video"})
		return
	}

	writeJSON(w, http.StatusOK, video)
}
