package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	httphandler "github.com/rubenmarques/youtube-api-integration/internal/adapter/http"
	"github.com/rubenmarques/youtube-api-integration/internal/domain"
)

// mockRepo is a test double for domain.VideoRepository.
type mockRepo struct {
	video *domain.Video
	err   error
}

func (m *mockRepo) GetVideoByID(_ context.Context, _ string) (*domain.Video, error) {
	return m.video, m.err
}

func newTestHandler(repo domain.VideoRepository) *httphandler.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return httphandler.New(repo, logger)
}

func TestGetVideo_MissingID(t *testing.T) {
	h := newTestHandler(&mockRepo{})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/video", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestGetVideo_NotFound(t *testing.T) {
	h := newTestHandler(&mockRepo{err: domain.ErrNotFound})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/video?id=abc", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetVideo_InternalError(t *testing.T) {
	h := newTestHandler(&mockRepo{err: errors.New("upstream failure")})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/video?id=abc", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestGetVideo_Success(t *testing.T) {
	want := &domain.Video{
		Title:        "Test Video",
		Description:  "A test video",
		ViewCount:    "1000",
		LikeCount:    "50",
		ChannelTitle: "Test Channel",
		PublishedAt:  "2024-01-01T00:00:00Z",
	}
	h := newTestHandler(&mockRepo{video: want})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/video?id=abc123", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var got domain.Video
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if got.Title != want.Title {
		t.Errorf("title: got %q, want %q", got.Title, want.Title)
	}
}
