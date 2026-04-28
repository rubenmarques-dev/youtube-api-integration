package youtube

import (
	"context"
	"fmt"

	"github.com/rubenmarques/youtube-api-integration/internal/domain"
	"google.golang.org/api/option"
	youtubeapi "google.golang.org/api/youtube/v3"
)

// Client implements domain.VideoRepository using the YouTube Data API v3.
type Client struct {
	svc *youtubeapi.Service
}

// New creates a new YouTube API client authenticated with the given API key.
func New(apiKey string) (*Client, error) {
	svc, err := youtubeapi.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("creating youtube service: %w", err)
	}
	return &Client{svc: svc}, nil
}

// GetVideoByID fetches a single video by ID from the YouTube Data API.
func (c *Client) GetVideoByID(ctx context.Context, id string) (*domain.Video, error) {
	call := c.svc.Videos.List([]string{"snippet", "statistics"}).Id(id).Context(ctx)
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("youtube api call failed: %w", err)
	}
	if len(resp.Items) == 0 {
		return nil, domain.ErrNotFound
	}
	item := resp.Items[0]
	v := &domain.Video{
		Title:        item.Snippet.Title,
		Description:  item.Snippet.Description,
		ViewCount:    fmt.Sprintf("%d", item.Statistics.ViewCount),
		LikeCount:    fmt.Sprintf("%d", item.Statistics.LikeCount),
		ChannelTitle: item.Snippet.ChannelTitle,
		PublishedAt:  item.Snippet.PublishedAt,
	}
	v.Chapters = domain.ParseChapters(item.Snippet.Description)
	return v, nil
}
