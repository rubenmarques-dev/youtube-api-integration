package domain

// Chapter represents a single chapter extracted from a video description.
type Chapter struct {
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"` // human-readable, e.g. "1:23" or "1:23:45"
	Seconds   int    `json:"seconds"`   // total seconds from start
}

// Video represents the core video entity returned by the application.
type Video struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ViewCount    string    `json:"viewCount"`
	LikeCount    string    `json:"likeCount"`
	ChannelTitle string    `json:"channelTitle"`
	PublishedAt  string    `json:"publishedAt"`
	Chapters     []Chapter `json:"chapters"`
}
