package domain

import (
	"context"
	"errors"
)

// VideoRepository is the port that must be implemented by any video data adapter.
type VideoRepository interface {
	GetVideoByID(ctx context.Context, id string) (*Video, error)
}

// ErrNotFound is returned when a video with the requested ID cannot be found.
var ErrNotFound = errors.New("video not found")
